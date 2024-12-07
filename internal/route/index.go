package route

import (
	"github.com/SeakMengs/yato-cdn/internal/controller"
	"github.com/gin-gonic/gin"
)

func V1_Index(r *gin.RouterGroup, indexController *controller.IndexController) {
	v1 := r.Group("/v1")
	{
		v1.GET("/regions", indexController.GetAllRegions)
	}
}
