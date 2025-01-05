package routers

import (
	"RamadanRangkuti/service-tickitz/internal/handlers"

	"github.com/gin-gonic/gin"
)

func OrderRouter(r *gin.RouterGroup) {
	r.POST("/", handlers.CreateTransaction)
	r.GET("/cinema", handlers.GetCinema)
	r.GET("/seats", handlers.GetAvailableSeats)
	r.POST("/payment", handlers.Payment)
	r.GET("/ticket", handlers.GetTicket)
}
