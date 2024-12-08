package controller

import (
	"fmt"
	"io"
	"net/http"

	"github.com/SeakMengs/yato-cdn/internal/config"
	"github.com/SeakMengs/yato-cdn/internal/file"
	"github.com/SeakMengs/yato-cdn/internal/util"
	"github.com/gin-gonic/gin"
)

type FileController struct {
	*baseController
	CDN config.CDN
}

// Region, name
const (
	STORAGE_PATH        = "./storage/%s/%s"
	FILE_SERVE_FILE_API = "/api/v1/files/"
)

func (fc *FileController) GetAllFileNames(ctx *gin.Context) {
	fn, err := fc.app.Repository.File.GetAll(ctx, nil)
	if err != nil {
		util.ResponseFailed(ctx, http.StatusInternalServerError, "Failed to retrieve file informations", err, nil)
		return
	}

	util.ResponseSuccess(ctx, fn)
}

func (fc *FileController) UploadFile(ctx *gin.Context) {
	// if fc.CDN.IsCDN {
	// 	util.ResponseFailed(ctx, http.StatusBadRequest, "This server is being serve as a distribute server", nil, nil)
	// }

	// Get the uploaded file
	_file, err := ctx.FormFile("file")
	if err != nil {
		fc.app.Logger.Debugf("No file is attached. err: %v", err)
		util.ResponseFailed(ctx, http.StatusBadRequest, "No file is attached", err, nil)
		return
	}

	// Set the destination file path
	dst := fmt.Sprintf(STORAGE_PATH, fc.CDN.Region, _file.Filename)
	fc.app.Logger.Debugf("File path: %s\n", dst)

	// Save the uploaded file
	err = file.Save(_file, dst)
	if err != nil {
		fc.app.Logger.Debugf("Failed to save file. err: %v", err)
		util.ResponseFailed(ctx, http.StatusInternalServerError, "Failed to save file", err, nil)
		return
	}

	fileUrl := fmt.Sprintf("%s%s%s", ctx.Request.Host, FILE_SERVE_FILE_API, _file.Filename)
	util.ResponseSuccess(ctx, gin.H{"message": "File uploaded successfully", "file": fileUrl})
}

func (fc *FileController) ReadFile(ctx *gin.Context) {
	filename := ctx.Param("filename")
	dst := fmt.Sprintf(STORAGE_PATH, fc.CDN.Region, filename)

	// Open the file
	file, err := file.Read(dst)
	if err != nil {
		fc.app.Logger.Debugf("Failed to read file. err: %v", err)
		util.ResponseFailed(ctx, http.StatusInternalServerError, "Failed to read file", err, nil)
		return
	}
	defer file.Content.Close()

	// Detect the MIME type
	ctx.Header("Content-Type", file.MimeType)

	_, err = io.Copy(ctx.Writer, file.Content)
	if err != nil {
		fc.app.Logger.Debugf("Failed to serve file. err: %v", err)
		util.ResponseFailed(ctx, http.StatusInternalServerError, "Failed to serve file", err, nil)
	}
}

func (fc *FileController) DeleteFile(ctx *gin.Context) {
	filename := ctx.Param("filename")

	dst := fmt.Sprintf(STORAGE_PATH, fc.CDN.Region, filename)

	err := file.Delete(dst)
	if err != nil {
		fc.app.Logger.Debugf("Failed to delete file. err: %v", err)
		util.ResponseFailed(ctx, http.StatusInternalServerError, "Failed to delete file", err, nil)
		return
	}

	util.ResponseSuccess(ctx, gin.H{"message": "File deleted successfully"})
}
