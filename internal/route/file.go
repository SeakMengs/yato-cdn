package route

import (
	"github.com/SeakMengs/yato-cdn/internal/controller"
	"github.com/gin-gonic/gin"
)

func V1_File(r *gin.RouterGroup, fileController *controller.FileController) {
	v1 := r.Group("/v1/files")
	{
		v1.GET("", fileController.GetAllFileNames)
		// Test endpoint with curl: curl --output - -X GET http://localhost:8080/api/v1/files/your file name
		v1.GET("/:filename", fileController.ReadFile)
		// Test endpoint with curl: curl -X POST -F "file=@path/to/your/file" http://localhost:8080/api/v1/files/upload
		v1.DELETE("/:filename", fileController.DeleteFile)
		v1.POST("/upload", fileController.UploadFile)
	}
}
