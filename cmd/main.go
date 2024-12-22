package main

import (
	"RamadanRangkuti/service-tickitz/internal/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routers.Routers(r)
	r.Run("localhost:8083") // listen and serve on 0.0.0.0:8080
}
