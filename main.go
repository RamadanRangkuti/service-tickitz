package main

import (
	"RamadanRangkuti/service-tickitz/internal/routers"

	docs "RamadanRangkuti/service-tickitz/docs"

	"github.com/gin-gonic/gin"
	swaggerfile "github.com/swaggo/files"      // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @host: localhost:8082
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfile.Handler))
	routers.Routers(r)
	r.Run("localhost:8082")
}
