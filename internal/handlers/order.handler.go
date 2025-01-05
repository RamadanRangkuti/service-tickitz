package handlers

import (
	"RamadanRangkuti/service-tickitz/internal/dto"
	"RamadanRangkuti/service-tickitz/internal/repositories"
	"RamadanRangkuti/service-tickitz/pkg"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Find cinemas by movie and location
// @Schemes
// @Description Retrieve a list of cinemas based on the movie ID, location, date, and time.
// @Tags Cinemas
// @Accept mpfd
// @Produce json
// @Param movie_id query int true "Movie ID"
// @Param location query string true "Location"
// @Param date query string true "Date (format: YYYY-MM-DD)"
// @Param time query string true "Time (format: HH:mm:ss)"
// @Success 200 {object} pkg.Response{data=[]models.GetCinema}
// @Failure 400 {object} pkg.Response{error=string} "Invalid input"
// @Failure 500 {object} pkg.Response{error=string} "Internal server error"
// @Router /order/cinemas [get]
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

// @Summary Retrieve available seats for a cinema
// @Schemes
// @Description Get a list of available seats for a specific movie and cinema schedule.
// @Tags Seats
// @Accept mpfd
// @Produce json
// @Param movie_id formData int true "Movie ID"
// @Param cinema_id formData int true "Cinema ID"
// @Param date formData string true "Date (format: YYYY-MM-DD)"
// @Param time formData string true "Time (format: HH:mm:ss)"
// @Success 200 {object} pkg.Response{data=[]models.CinemaSeat}
// @Failure 400 {object} pkg.Response{error=string} "Invalid input"
// @Failure 500 {object} pkg.Response{error=string} "Internal server error"
// @Router /order/seats [post]
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

// CreateTransaction godoc
// @Summary Process a payment for a transaction
// @Schemes
// @Description Process a payment for a given transaction ID and update its status.
// @Tags Transactions
// @Accept mpfd
// @Produce json
// @Param user_id formData int true "User ID"
// @Param movie_id formData int true "Movie ID"
// @Param cinema_id formData int true "Cinema ID"
// @Param date formData string true "Date (format: YYYY-MM-DD)"
// @Param time formData string true "Time (format: HH:mm:ss)"
// @Param location formData string true "Location"
// @Param seat_numbers[] formData []string true "Seat Numbers" collectionFormat(multi)
// @Success 200 {object} pkg.Response{data=string} "Transaction created successfully"
// @Failure 400 {object} pkg.Response{error=string} "Invalid input"
// @Failure 500 {object} pkg.Response{error=string} "Internal server error"
// @Router /order/transactions [post]
func CreateTransaction(c *gin.Context) {
	response := pkg.NewResponse(c)
	var input dto.CreateTransactionDTO
	type orderResponse struct {
		VirtualAcoout string `json:"virtualAccount"`
		TotalPayment  int    `json:"totalPayment"`
		Due           string `json:"due"`
	}

	if err := c.ShouldBind(&input); err != nil {
		log.Printf("Binding Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	// Log array setelah parsing
	log.Printf("Received input: %+v\n", input.SeatNumbers)

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

// @Summary Process a payment for a transaction
// @Schemes
// @Description Process a payment for a given transaction ID and update its status.
// @Tags Transactions
// @Accept mpfd
// @Produce json
// @Param transaction_id formData int true "Transaction ID"
// @Param payment_method formData string true "Payment Method"
// @Param virtual_account_number formData string true "VA Number"
// @Param amount formData string true "Ammount"
// @Success 200 {object} pkg.Response{data=string} "Payment processed successfully"
// @Failure 400 {object} pkg.Response{error=string} "Invalid input"
// @Failure 500 {object} pkg.Response{error=string} "Internal server error"
// @Router /order/payments [post]
func Payment(c *gin.Context) {
	response := pkg.NewResponse(c)
	var input dto.CreatePaymentDTO

	if err := c.ShouldBind(&input); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	if err := repositories.CreatePayment(input); err != nil {
		response.InternalServerError("Failed to process payment", err.Error())
		return
	}

	response.Success("Payment processed successfully", nil)
}

// GetTicket godoc
// @Summary Retrieve ticket details for a transaction
// @Schemes
// @Description Get ticket information based on a transaction ID.
// @Tags Tickets
// @Accept mpfd
// @Produce json
// @Param transaction_id query int true "Transaction ID"
// @Success 200 {object} pkg.Response{data=models.Ticket} "Ticket details retrieved successfully"
// @Failure 400 {object} pkg.Response{error=string} "Invalid input"
// @Failure 500 {object} pkg.Response{error=string} "Internal server error"
// @Router /order/ticket [get]
func GetTicket(c *gin.Context) {
	response := pkg.NewResponse(c)
	var input dto.GetTicketDTO

	if err := c.ShouldBind(&input); err != nil {
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
