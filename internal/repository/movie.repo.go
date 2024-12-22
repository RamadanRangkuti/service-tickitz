package repository

import (
	"RamadanRangkuti/service-tickitz/internal/models"
	"RamadanRangkuti/service-tickitz/pkg"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// func GetMovies() ([]models.MovieDetails, error) {
// 	conn, err := pkg.DB()
// 	if err != nil {
// 		return nil, fmt.Errorf("connection failed: %v", err)
// 	}
// 	defer conn.Close(context.Background())

// 	query := `
// 		SELECT
// 			m.id AS movie_id,
// 			m.title AS movie_title,
// 			m.synopsis AS movie_synopsis,
// 			m.duration AS movie_duration,
// 			m.realease_date AS movie_release_date,
// 			m.image AS movie_image,
// 			m.banner AS movie_banner,
// 			md.name AS director_name,
// 			ARRAY_AGG(DISTINCT mc.name) AS casts,
// 			ARRAY_AGG(DISTINCT mg.name) AS genres,
// 			m.created_at AS movie_created_at,
// 			m.updated_at AS movie_updated_at
// 		FROM
// 			movies m
// 		LEFT JOIN
// 			movie_directors md ON m.director_id = md.id
// 		LEFT JOIN
// 			movie_casts mc ON m.id = mc.movie_id
// 		LEFT JOIN
// 			movie_genres mg ON m.id = mg.movie_id
// 		GROUP BY
// 			m.id, md.name
// 	`

// 	rows, err := conn.Query(context.Background(), query)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to execute query: %v", err)
// 	}
// 	defer rows.Close()

// 	var movies []models.MovieDetails
// 	for rows.Next() {
// 		var movie models.MovieDetails
// 		err = rows.Scan(
// 			&movie.Id,
// 			&movie.Title,
// 			&movie.Synopsis,
// 			&movie.Duration,
// 			&movie.ReleaseDate,
// 			&movie.Image,
// 			&movie.Banner,
// 			&movie.DirectorName,
// 			&movie.Casts,
// 			&movie.Genres,
// 			&movie.CreatedAt,
// 			&movie.UpdatedAt,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan row: %v", err)
// 		}
// 		movies = append(movies, movie)
// 	}

// 	return movies, nil
// }

func FindAllMovies(page int, limit int, order string, sortBy string, search string) (models.ListMovies, error) {
	conn, err := pkg.DB()
	if err != nil {
		return nil, fmt.Errorf("connection failed: %v", err)
	}
	defer conn.Close(context.Background())

	offset := (page - 1) * limit
	search = fmt.Sprintf("%%%s%%", search)
	query := fmt.Sprintf(`
    SELECT 
        m.id, 
        m.title, 
        m.synopsis, 
        m.duration, 
        m.realease_date, 
        m.image, 
        m.banner,
        COALESCE(ARRAY_AGG(DISTINCT mc.name), '{}') AS casts,
    	COALESCE(ARRAY_AGG(DISTINCT mg.name), '{}') AS genres,
        md.name AS director,
		m.created_at,
		m.updated_at 
    FROM movies m
    LEFT JOIN movie_casts mc ON mc.movie_id = m.id
    LEFT JOIN movie_genres mg ON mg.movie_id = m.id
    LEFT JOIN movie_directors md ON md.id = m.director_id
    WHERE LOWER(m.title) ILIKE LOWER($1)
    GROUP BY m.id, md.name
    ORDER BY %s %s
    LIMIT $2 OFFSET $3
`, sortBy, order)
	rows, err := conn.Query(context.Background(), query, search, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	movies, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.MovieDetails])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %v", err)
	}
	return movies, nil
}

func FindMovieById(id int) (*models.MovieDetails, error) {
	var movie models.MovieDetails
	conn, err := pkg.DB()
	if err != nil {
		fmt.Println("connection failed", err)
	}
	defer conn.Close(context.Background())

	query := ` SELECT
        m.id,
        m.title,
        m.synopsis,
        m.duration,
        m.realease_date,
        m.image,
        m.banner,
        COALESCE(ARRAY_AGG(DISTINCT mc.name), '{}') AS casts,
        COALESCE(ARRAY_AGG(DISTINCT mg.name), '{}') AS genres,
        md.name AS director,
        m.created_at,
        m.updated_at
    FROM movies m
    LEFT JOIN movie_casts mc ON mc.movie_id = m.id
    LEFT JOIN movie_genres mg ON mg.movie_id = m.id
    LEFT JOIN movie_directors md ON md.id = m.director_id
    WHERE m.id = $1
    GROUP BY m.id, md.name`

	rows, err := conn.Query(context.Background(), query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()
	// Memastikan hanya satu baris yang dikembalikan
	if rows.Next() {
		err :=
			rows.Scan(&movie.Id, &movie.Title, &movie.Synopsis, &movie.Duration, &movie.RealeaseDate, &movie.Image, &movie.Banner, &movie.Casts, &movie.Genres, &movie.Director, &movie.CreatedAt, &movie.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
	} else {
		return nil, fmt.Errorf("movie with ID %d not found", id)
	}

	return &movie, nil
}

func CountMovie(search string) int {
	conn, err := pkg.DB()
	if err != nil {
		fmt.Println("connection failed", err)
	}
	defer conn.Close(context.Background())
	var total int
	search = fmt.Sprintf("%%%s%%", search)

	conn.QueryRow(context.Background(), `
	SELECT COUNT(id) FROM movies WHERE title ILIKE $1
	`, search).Scan(&total)

	return total
}
