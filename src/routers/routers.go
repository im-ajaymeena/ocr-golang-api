package routers

import (
	"ajaymeena/ocr/api"
	docs "ajaymeena/ocr/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func GetRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())
	r.Use(CORS)

	docs.SwaggerInfo.BasePath = "/"

	r.POST("/image-sync", api.OCRSynchronous)
	r.POST("/upload", api.GetBase64)

	r.POST("/image", api.CreateOCRTask)

	r.GET("/image", api.GetOCRTaskResult)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}
