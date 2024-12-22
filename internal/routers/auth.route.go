package routers

import (
	"RamadanRangkuti/service-tickitz/internal/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRouter(r *gin.RouterGroup) {
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
}
