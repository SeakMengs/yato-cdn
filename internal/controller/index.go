package controller

import (
	"net/http"

	"github.com/SeakMengs/yato-cdn/internal/util"
	"github.com/gin-gonic/gin"
)

type IndexController struct {
	*baseController
}

func (ic IndexController) Index(ctx *gin.Context) {
	util.ResponseSuccess(ctx, gin.H{
		"message": "Welcome to the api",
	})
}

func (ic *IndexController) GetAllRegions(ctx *gin.Context) {
	regions, err := ic.app.Repository.Region.GetAll(ctx, nil)
	if err != nil {
		util.ResponseFailed(ctx, http.StatusInternalServerError, "Failed to retrieve region informations", err, nil)
		return
	}

	util.ResponseSuccess(ctx, regions)
}
