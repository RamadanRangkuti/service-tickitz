package handlers

import (
	"RamadanRangkuti/service-tickitz/internal/dto"
	"RamadanRangkuti/service-tickitz/internal/repositories"
	"RamadanRangkuti/service-tickitz/pkg"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func GetCinema(c *gin.Context) {
	response := pkg.NewResponse(c)
	var input dto.GetCinemaDTO

	if err := c.ShouldBind(&input); err != nil {
		fmt.Println("Binding Error:", err)
		response.BadRequest("Invalid input", err.Error())
		return
	}

	fmt.Println("input", input)
	// Query cinema
	cinema, err := repositories.FindCinema(input)
	if err != nil {
		response.InternalServerError("Failed to get cinema", err.Error())
		return
	}

	response.Success("Cinema result", cinema)
}

func GetAvailableSeats(c *gin.Context) {
	response := pkg.NewResponse(c)
	var input dto.GetAvailableSeatsDTO

	if err := c.ShouldBind(&input); err != nil {
		fmt.Println("Binding Error:", err)
		response.BadRequest("Invalid input", err.Error())
		return
	}

	seats, err := repositories.FindAvailableSeats(input)
	if err != nil {
		response.InternalServerError("Failed to get available seats", err.Error())
		return
	}

	response.Success("Available seats", seats)
}

func CreateTransaction(c *gin.Context) {
	response := pkg.NewResponse(c)
	var input dto.CreateTransactionDTO

	type orderResponse struct {
		VirtualAcoout string `json:"virtualAccount"`
		TotalPayment  int    `json:"totalPayment"`
		Due           string `json:"due"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("Binding Error:", err)
		response.BadRequest("Invalid input", err.Error())
		return
	}

	total_price, err := repositories.CreateTransaction(input)
	if err != nil {
		response.InternalServerError("Failed to create transaction", err.Error())
		return
	}
	dueTime := time.Now().Add(2 * time.Hour).Format("2006-01-02 15:04:05")

	response.Success("Order Succes", orderResponse{
		VirtualAcoout: "324234534324",
		TotalPayment:  total_price,
		Due:           dueTime,
	})
}

func Payment(c *gin.Context) {
	response := pkg.NewResponse(c)
	var input dto.CreatePaymentDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	if err := repositories.CreatePayment(input); err != nil {
		response.InternalServerError("Failed to process payment", err.Error())
		return
	}

	response.Success("Payment processed successfully", nil)
}

func GetTicket(c *gin.Context) {
	response := pkg.NewResponse(c)
	var input dto.GetTicketDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	ticket, err := repositories.FetchTicketDetails(input.TransactionId)
	if err != nil {
		response.InternalServerError("Failed to fetch ticket details", err.Error())
		return
	}

	response.Success("Ticket details fetched successfully", ticket)
}
