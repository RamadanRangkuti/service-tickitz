package routers

import (
	handlers "RamadanRangkuti/service-tickitz/internal/handlers/users"
	"RamadanRangkuti/service-tickitz/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.RouterGroup) {
	r.GET("/", middlewares.ValidateToken(), handlers.GetAllUser)
	r.GET("/:id", middlewares.ValidateToken(), handlers.GetUserById)
	r.POST("/", middlewares.ValidateToken(), handlers.CreateUser)
	r.PATCH("/:id", middlewares.ValidateToken(), handlers.UpdateUser)
	r.DELETE("/:id", middlewares.ValidateToken(), handlers.DeleteUser)
}
