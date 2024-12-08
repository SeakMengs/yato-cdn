package route

import (
	"github.com/SeakMengs/yato-cdn/internal/controller"
	"github.com/gin-gonic/gin"
)

func V1_CDN(r *gin.RouterGroup, cdnController *controller.CDNController) {
	v1 := r.Group("/v1/cdn")
	{
		v1.GET("/:filename/get-region", cdnController.ServeFile)
		// Test endpoint with curl: curl --output - -X GET http://localhost:8080/api/v1/cdn/your file name
		v1.GET("/:filename", cdnController.ServeFile)
		// Test endpoint with curl: curl -X POST -F "file=@path/to/your/file" http://localhost:8080/api/v1/cdn/upload
		v1.POST("/upload", cdnController.UploadFile)
	}
}
