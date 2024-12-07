package controller

import (
	"fmt"
	"io"
	"net/http"

	"github.com/SeakMengs/yato-cdn/internal/file"
	"github.com/SeakMengs/yato-cdn/internal/model"
	"github.com/SeakMengs/yato-cdn/internal/util"
	"github.com/gin-gonic/gin"
)

type FileController struct {
	*baseController
}

func (fc *FileController) GetAllFileNames(ctx *gin.Context) {
	fn, err := fc.app.Repository.File.GetAll(ctx, nil)
	if err != nil {
		util.ResponseFailed(ctx, http.StatusInternalServerError, "Failed to retrieve file informations", err, nil)
		return
	}

	util.ResponseSuccess(ctx, fn)
}

// Region, name
const STORAGE_PATH = "./storage/%s/%s"

func (fc *FileController) UploadFile(ctx *gin.Context) {
	// Get the uploaded file
	_file, err := ctx.FormFile("file")
	if err != nil {
		util.ResponseFailed(ctx, http.StatusBadRequest, "No file is attached", err, nil)
		return
	}

	// Set the destination file path
	dst := fmt.Sprintf(STORAGE_PATH, "singapore", _file.Filename)

	// Save the uploaded file
	err = file.Save(_file, dst)
	if err != nil {
		util.ResponseFailed(ctx, http.StatusInternalServerError, "Failed to save file", err, nil)
		return
	}

	err = fc.app.Repository.File.Save(ctx, nil, model.File{
		Name: _file.Filename,
	})
	if err != nil {
		// If save in db fail, remove the file we save earlier
		file.Delete(dst)

		util.ResponseFailed(ctx, http.StatusInternalServerError, "Failed to save file", err, nil)
		return
	}

	util.ResponseSuccess(ctx, gin.H{"message": "File uploaded successfully", "file": dst})
}

func (fc FileController) ReadFile(ctx *gin.Context) {
	filename := ctx.Param("filename")
	dst := fmt.Sprintf(STORAGE_PATH, "singapore", filename)

	// Open the file
	file, err := file.Read(dst)
	if err != nil {
		util.ResponseFailed(ctx, http.StatusInternalServerError, "Failed to serve file", err, nil)
		return
	}
	defer file.Content.Close()

	// // Detect the MIME type
	ctx.Header("Content-Type", file.MimeType)

	_, err = io.Copy(ctx.Writer, file.Content)
	if err != nil {
		util.ResponseFailed(ctx, http.StatusInternalServerError, "Failed to serve file", err, nil)
	}
}
