package routers

import (
	handlers "RamadanRangkuti/service-tickitz/internal/handlers/users"
	"RamadanRangkuti/service-tickitz/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.RouterGroup) {
	r.GET("/", middlewares.ValidateToken(), handlers.GetAllUser)
	r.GET("/:id", handlers.GetUserById)
	r.POST("/", handlers.CreateUser)
	r.PATCH("/:id", handlers.UpdateUser)
	r.DELETE("/:id", handlers.DeleteUser)
}
