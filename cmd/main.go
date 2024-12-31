package main

import (
	"RamadanRangkuti/service-tickitz/internal/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routers.Routers(r)
	r.Run("localhost:8080")
}
