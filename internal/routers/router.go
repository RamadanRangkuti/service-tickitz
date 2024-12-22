package routers

import "github.com/gin-gonic/gin"

func Routers(router *gin.Engine) {
	TesRouter(router.Group("/tes"))
	MovieRouter(router.Group("/movies"))
	UserRouter(router.Group("/users"))
	AuthRouter(router.Group("/auth"))
}
