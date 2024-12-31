package routers

import (
	"RamadanRangkuti/service-tickitz/internal/handlers"

	"github.com/gin-gonic/gin"
)

func MovieRouter(r *gin.RouterGroup) {
	r.GET("/", handlers.GetMovies)
	r.GET("/:id", handlers.GetMovieById)
}
