package repositories

import (
	"RamadanRangkuti/service-tickitz/internal/dto"
	"RamadanRangkuti/service-tickitz/internal/models"
	"RamadanRangkuti/service-tickitz/pkg"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func FindCinema(input dto.GetCinemaDTO) (models.GetCinemas, error) {
	conn, err := pkg.DB()
	if err != nil {
		fmt.Println("connection failed", err)
	}
	defer conn.Close(context.Background())

	query := `
	SELECT ms.id, 
		c.name, 
		l.city, 
		ms.date, 
		ms.time
	FROM movie_schedules ms
	JOIN cinemas c ON ms.cinema_id = c.id
	JOIN cinema_locations cl ON c.id = cl.cinema_id
	JOIN locations l ON cl.location_id = l.id
	WHERE ms.movie_id = $1 AND l.city = $2 AND ms.date = $3 AND ms.time = $4
	`

	rows, err := conn.Query(context.Background(), query, input.MovieId, input.Location, input.Date, input.Time)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	cinema, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.GetCinema])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %v", err)
	}
	return cinema, nil
}

func FindAvailableSeats(input dto.GetAvailableSeatsDTO) ([]models.CinemaSeat, error) {
	conn, err := pkg.DB()
	if err != nil {
		fmt.Println("connection failed", err)
		return nil, err
	}
	defer conn.Close(context.Background())

	query := `
	SELECT cs.seat_number, cs.is_available, sp.category, sp.price
	FROM cinema_seats cs
	JOIN seat_prices sp ON cs.seat_prices_id = sp.id
	JOIN movie_schedules ms ON cs.cinema_id = ms.cinema_id
	WHERE ms.movie_id = $1
	  AND ms.date = $2
	  AND ms.time = $3
	  AND ms.cinema_id = $4
	  AND cs.is_available = true;
	`

	rows, err := conn.Query(context.Background(), query, input.MovieId, input.Date, input.Time, input.CinemaId)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var seats []models.CinemaSeat
	for rows.Next() {
		var seat models.CinemaSeat
		if err := rows.Scan(&seat.SeatNumber, &seat.IsAvailable, &seat.Category, &seat.Price); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		seats = append(seats, seat)
	}
	return seats, nil
}

func CreateTransaction(input dto.CreateTransactionDTO) (int, error) {
	conn, err := pkg.DB()
	if err != nil {
		fmt.Println("connection failed", err)
		return 0, err
	}
	defer conn.Close(context.Background())

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	var totalPrice float64
	for _, seatNumber := range input.SeatNumbers {
		var price float64
		query := `
			SELECT sp.price
			FROM cinema_seats cs
			JOIN seat_prices sp ON cs.seat_prices_id = sp.id
			WHERE cs.cinema_id = $1 AND cs.seat_number = $2
		`
		err := tx.QueryRow(context.Background(), query, input.CinemaId, seatNumber).Scan(&price)
		if err != nil {
			return 0, fmt.Errorf("failed to get seat price for seat %s: %v", seatNumber, err)
		}
		totalPrice += price
	}

	query := `
		INSERT INTO transactions (user_id, movie_id, date, location, quantity, total_price, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var transactionId int
	err = tx.QueryRow(context.Background(), query, input.UserId, input.MovieId, input.Date+"T"+input.Time, input.Location, len(input.SeatNumbers), totalPrice, "pending").Scan(&transactionId)
	if err != nil {
		return 0, fmt.Errorf("failed to insert transaction: %v", err)
	}

	for _, seatNumber := range input.SeatNumbers {
		query = `
			INSERT INTO transaction_seats (transaction_id, seat_number, cinema_id)
			VALUES ($1, $2, $3)
		`
		_, err = tx.Exec(context.Background(), query, transactionId, seatNumber, input.CinemaId)
		if err != nil {
			return 0, fmt.Errorf("failed to insert transaction seat: %v", err)
		}

		// Update seat availability
		query = `
			UPDATE cinema_seats
			SET is_available = false
			WHERE cinema_id = $1 AND seat_number = $2
		`
		_, err = tx.Exec(context.Background(), query, input.CinemaId, seatNumber)
		if err != nil {
			return 0, fmt.Errorf("failed to update seat availability for seat %s: %v", seatNumber, err)
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return int(totalPrice), nil
}

func CreatePayment(input dto.CreatePaymentDTO) error {
	conn, err := pkg.DB()
	if err != nil {
		log.Printf("connection failed: %v", err)
		return err
	}
	defer conn.Close(context.Background())

	query := `
		INSERT INTO payments (transaction_id, payment_method, virtual_account_number, amount, status)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = conn.Exec(context.Background(), query, input.TransactionId, input.PaymentMethod, input.VirtualAccountNumber, input.Amount, "success")
	if err != nil {
		return fmt.Errorf("failed to insert payment: %v", err)
	}

	// Update transaction status
	query = `
		UPDATE transactions
		SET status = 'completed'
		WHERE id = $1
	`
	_, err = conn.Exec(context.Background(), query, input.TransactionId)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %v", err)
	}

	return nil
}

func FetchTicketDetails(transactionId int) (models.Ticket, error) {
	conn, err := pkg.DB()
	if err != nil {
		return models.Ticket{}, fmt.Errorf("database connection failed: %v", err)
	}
	defer conn.Close(context.Background())

	query := `
		SELECT 
			m.title AS movie_name,
			TO_CHAR(t.date, 'DD Mon YYYY') AS date,
			TO_CHAR(t.date, 'HH12:MI PM') AS time,
			t.quantity,
			ARRAY_AGG(ts.seat_number) AS seats,
			t.total_price
		FROM transactions t
		JOIN transaction_seats ts ON t.id = ts.transaction_id
		JOIN movies m ON t.movie_id = m.id
		WHERE t.id = $1
		GROUP BY m.title, t.date, t.quantity, t.total_price
	`
	var ticket models.Ticket
	err = conn.QueryRow(context.Background(), query, transactionId).Scan(
		&ticket.MovieTitle,
		&ticket.Date,
		&ticket.Time,
		&ticket.Quantity,
		&ticket.Seats,
		&ticket.TotalPrice,
	)
	if err != nil {
		return models.Ticket{}, fmt.Errorf("failed to fetch ticket details: %v", err)
	}

	return ticket, nil
}
