package routers

import (
	"RamadanRangkuti/service-tickitz/internal/handlers"

	"github.com/gin-gonic/gin"
)

func OrderRouter(r *gin.RouterGroup) {
	r.POST("/transactions", handlers.CreateTransaction)
	r.GET("/cinemas", handlers.GetCinema)
	r.POST("/seats", handlers.GetAvailableSeats)
	r.POST("/payments", handlers.Payment)
	r.GET("/ticket", handlers.GetTicket)
}
