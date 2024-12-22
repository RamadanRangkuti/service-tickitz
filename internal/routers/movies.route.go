package routers

import (
	handlers "RamadanRangkuti/service-tickitz/internal/handlers/movies"

	"github.com/gin-gonic/gin"
)

func MovieRouter(r *gin.RouterGroup) {
	r.GET("/", handlers.GetMovie)
}
