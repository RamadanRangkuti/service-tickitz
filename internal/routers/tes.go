package routers

import (
	handlers "RamadanRangkuti/service-tickitz/internal/handlers/users"

	"github.com/gin-gonic/gin"
)

func TesRouter(r *gin.RouterGroup) {
	r.GET("/", handlers.GetTes)
}
