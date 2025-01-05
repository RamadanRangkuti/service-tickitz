package routers

import (
	"github.com/gin-gonic/gin"
)

func Routers(router *gin.Engine) {
	MovieRouter(router.Group("/movies"))
	UserRouter(router.Group("/users"))
	AuthRouter(router.Group("/auth"))
	OrderRouter(router.Group("/order"))
}
