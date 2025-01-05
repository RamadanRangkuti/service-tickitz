package dto

type GetCinemaDTO struct {
	MovieId  int    `form:"movie_id" binding:"required"`
	Location string `form:"location" binding:"required"`
	Date     string `form:"date" binding:"required"`
	Time     string `form:"time" binding:"required"`
}

type GetAvailableSeatsDTO struct {
	MovieId  int    `form:"movie_id" binding:"required"`
	CinemaId int    `form:"cinema_id" binding:"required"`
	Date     string `form:"date" binding:"required"`
	Time     string `form:"time" binding:"required"`
}

type CreateTransactionDTO struct {
	UserId      int      `form:"user_id" binding:"required"`        // ID pengguna yang melakukan transaksi
	MovieId     int      `form:"movie_id" binding:"required"`       // ID film
	CinemaId    int      `form:"cinema_id" binding:"required"`      // ID bioskop
	Date        string   `form:"date" binding:"required"`           // Tanggal jadwal (format: YYYY-MM-DD)
	Time        string   `form:"time" binding:"required"`           // Waktu jadwal (format: HH:mm:ss)
	Location    string   `form:"location" binding:"required"`       // Lokasi bioskop
	SeatNumbers []string `form:"seat_numbers[]" binding:"required"` // Daftar nomor kursi
}

type CreatePaymentDTO struct {
	TransactionId        int     `form:"transaction_id" binding:"required"`
	PaymentMethod        string  `form:"payment_method" binding:"required"`
	VirtualAccountNumber string  `form:"virtual_account_number" binding:"required"`
	Amount               float64 `form:"amount" binding:"required"`
}

type GetTicketDTO struct {
	TransactionId int `form:"transaction_id" binding:"required"`
}

// type CreateTransactionDTO struct {
// 	UserId      int      `json:"user_id" binding:"required"`
// 	MovieId     int      `json:"movie_id" binding:"required"`
// 	CinemaId    int      `json:"cinema_id" binding:"required"`
// 	Date        string   `json:"date" binding:"required"`
// 	Time        string   `json:"time" binding:"required"`
// 	Location    string   `json:"location" binding:"required"`
// 	SeatNumbers []string `json:"seat_numbers" binding:"required"`
// 	TotalPrice  float64  `json:"total_price" binding:"required"`
// }
