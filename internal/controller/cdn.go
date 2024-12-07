package controller

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/SeakMengs/yato-cdn/internal/config"
	"github.com/SeakMengs/yato-cdn/internal/model"
	"github.com/SeakMengs/yato-cdn/internal/util"
	"github.com/gin-gonic/gin"
)

type CDNController struct {
	*baseController
	CDN config.CDN
}

const (
	REGION_UPLOAD_API  = "/api/v1/files/upload"
	CDN_SERVE_FILE_API = "/api/v1/cdn/"
)

func (cdnc *CDNController) distributeFile(file *multipart.FileHeader, domain string, wg *sync.WaitGroup) {
	defer wg.Done()

	src, err := file.Open()
	if err != nil {
		cdnc.app.Logger.Debugf("Error opening file: %v\n", err)
		return
	}
	defer src.Close()

	// Create a buffer to hold the new multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file to the form data
	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		cdnc.app.Logger.Debugf("Error creating form file: %v\n", err)
		return
	}

	// Copy the file data into the multipart form
	_, err = io.Copy(part, src)
	if err != nil {
		cdnc.app.Logger.Debugf("Error writing file data to form: %v\n", err)
		return
	}

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		cdnc.app.Logger.Debugf("Error closing multipart writer: %v\n", err)
		return
	}

	url := fmt.Sprintf("%s%s", domain, REGION_UPLOAD_API)
	cdnc.app.Logger.Debugf("Distribute file to %s", url)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		cdnc.app.Logger.Debugf("Error creating request: %v\n", err)
		return
	}

	// Set content-type header to multipart/form-data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		cdnc.app.Logger.Debugf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body = &bytes.Buffer{}
	body.ReadFrom(resp.Body)
	cdnc.app.Logger.Debugf("Response body of %s: %s\n", domain, body.String())

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		cdnc.app.Logger.Debugf("Failed to upload file, status code: %d\n", resp.StatusCode)
		return
	}
}

func (cdnc *CDNController) UploadFile(ctx *gin.Context) {
	if !cdnc.CDN.IsCDN {
		util.ResponseFailed(ctx, http.StatusBadRequest, "This server is being serve as an edge server", nil, nil)
	}

	_file, err := ctx.FormFile("file")
	if err != nil {
		cdnc.app.Logger.Debugf("Error getting file: %v\n", err)
		util.ResponseFailed(ctx, http.StatusBadRequest, "No file is attached", err, nil)
		return
	}

	regions, err := cdnc.app.Repository.Region.GetAll(ctx, nil)
	if err != nil {
		cdnc.app.Logger.Debugf("Error getting regions: %v\n", err)
		util.ResponseFailed(ctx, http.StatusInternalServerError, "Failed to retrieve region informations", err, nil)
		return
	}

	wg := sync.WaitGroup{}

	// Intentionally not handle failure to distribute file to a region
	for _, r := range regions {
		wg.Add(1)
		go cdnc.distributeFile(_file, r.Domain, &wg)
	}

	wg.Wait()

	err = cdnc.app.Repository.File.Save(ctx, nil, model.File{
		Name: _file.Filename,
	})
	if err != nil {
		// TODO: if have time, create a delete file route such that we can remove the file in case we failed to save file name to database

		util.ResponseFailed(ctx, http.StatusInternalServerError, "Failed to save file", err, nil)
		return
	}

	fileUrl := fmt.Sprintf("%s%s%s", ctx.Request.Host, CDN_SERVE_FILE_API, _file.Filename)
	util.ResponseSuccess(ctx, gin.H{"message": "File uploaded successfully", "file": fileUrl})
}

func (cdnc *CDNController) findNearestRegion(regions []*model.Region, geo *util.GeoLocation) (*model.Region, error) {
	var nearestRegion *model.Region
	minDistance := math.MaxFloat64

	for _, r := range regions {
		serverGeo, err := util.FindGeoLocation(r.IP)
		if err != nil {
			cdnc.app.Logger.Debugf("Error getting geo location: %v\n", err)
			return nil, err
		}

		distance := util.Distance(geo.Latitude, geo.Longitude, serverGeo.Latitude, serverGeo.Longitude)
		if distance < minDistance {
			minDistance = distance
			nearestRegion = r
			cdnc.app.Logger.Debugf("Server region is %s, set in db is %s", serverGeo.Region, r.Name)
		}
	}

	return nearestRegion, nil
}

func (cdnc *CDNController) ServeFile(ctx *gin.Context) {
	if !cdnc.CDN.IsCDN {
		util.ResponseFailed(ctx, http.StatusBadRequest, "This server is being serve as an edge server", nil, nil)
	}
	filename := ctx.Param("filename")

	_, err := cdnc.app.Repository.File.GetByName(ctx, nil, filename)
	if err != nil {
		util.ResponseFailed(ctx, http.StatusNotFound, "File not found", err, nil)
		return
	}

	clientIp := ctx.ClientIP()
	// clientIp = "34.82.2.21" // US IP test, expect near US server
	// clientIp = "189.202.84.129" // Mexico IP test, expect near US server
	// clientIp = "203.113.152.18" // Vietnam IP test, expect near Korea server
	// clientIp = "110.74.193.124" // Cambodia IP test, expect near Singapore server
	// clientIp = "119.28.48.45" // Japan IP test, expect near Korea server
	// clientIp = "14.207.161.4" // Thailand IP test, expect near Singapore server

	cdnc.app.Logger.Debugf("Find geolocation for client %s\n", clientIp)
	clientGeo, err := util.FindGeoLocation(clientIp)
	if err != nil {
		util.ResponseFailed(ctx, http.StatusInternalServerError, "Failed to get client geo location", err, nil)
		return
	}

	regions, err := cdnc.app.Repository.Region.GetAll(ctx, nil)
	if err != nil {
		util.ResponseFailed(ctx, http.StatusInternalServerError, "Failed to retrieve region informations", err, nil)
		return
	}

	nearestRegion, err := cdnc.findNearestRegion(regions, clientGeo)
	cdnc.app.Logger.Debugf("Nearest region found: %s, %s for client %s, %s\n", nearestRegion.Name, nearestRegion.IP, clientGeo.Region, clientIp)
	if err != nil {
		util.ResponseFailed(ctx, http.StatusInternalServerError, "Failed to find nearest region", err, nil)
		return
	}

	url := fmt.Sprintf("%s%s%s", nearestRegion.Domain, FILE_SERVE_FILE_API, filename)
	cdnc.app.Logger.Debugf("Redirecting to %s\n", url)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}
