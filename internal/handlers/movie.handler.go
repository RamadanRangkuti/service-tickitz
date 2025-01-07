package handlers

import (
	"RamadanRangkuti/service-tickitz/internal/models"
	"RamadanRangkuti/service-tickitz/internal/repositories"
	"RamadanRangkuti/service-tickitz/pkg"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetAllMOvie godoc
// @Summary Get all movies
// @Schemes
// @Description Get all movies with pagination, sorting, and search
// @Tags Movies
// @Accept x-www-form-urlencoded
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(5)
// @Param sortBy query string false "Field to sort by" default(id)
// @Param order query string false "Order of sorting (asc or desc)" default(asc)
// @Param search query string false "Search query"
// @Success 200 {object} pkg.Response{data=[]models.MovieDetails,meta=pkg.PageInfo}
// @Failure 500 {object} pkg.Response{error=string}
// @Success 200 {object} models.MovieDetails
// @Security ApiKeyAuth
// @Router /movies [get]
func GetMovies(c *gin.Context) {
	response := pkg.NewResponse(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	sortBy := c.DefaultQuery("sortBy", "id")
	order := c.DefaultQuery("order", "asc")
	search := c.Query("search")
	if order != "asc" {
		order = "desc"
	}

	cacheKey := fmt.Sprintf("movies:list:page=%d:limit=%d:sort=%s:order=%s:search=%s", page, limit, sortBy, order, search)
	countKey := "movies:count"

	var movies []models.MovieDetails

	cachedMovies := pkg.Redis().Get(context.Background(), cacheKey)

	if cachedMovies.Val() != "" {
		movies, err := repositories.FindAllMovies(page, limit, order, sortBy, search)
		if err != nil {
			response.InternalServerError("Failed to fetch movies", err.Error())
			return
		}
		encoded, _ := json.Marshal(movies)
		pkg.Redis().Set(context.Background(), cacheKey, string(encoded), time.Minute*10)
	} else {
		json.Unmarshal([]byte(cachedMovies.Val()), &movies)
	}

	fmt.Println(movies)
	// cachedMovies := pkg.Redis().Get(context.Background(), cacheKey)
	// if cachedMovies.Val() != "" {
	// 	json.Unmarshal([]byte(cachedMovies.Val()), &movies)
	// } else {
	// 	movies, err := repositories.FindAllMovies(page, limit, order, sortBy, search)
	// 	if err != nil {
	// 		response.InternalServerError("Failed to fetch movies", err.Error())
	// 		return
	// 	}
	// 	encoded, _ := json.Marshal(movies)
	// 	pkg.Redis().Set(context.Background(), cacheKey, string(encoded), time.Minute*10)
	// }

	var count int
	cachedCount := pkg.Redis().Get(context.Background(), countKey)
	if cachedCount.Val() != "" {
		json.Unmarshal([]byte(cachedCount.Val()), &count)
	} else {
		count = repositories.CountMovie(search)
		encoded, _ := json.Marshal(count)
		pkg.Redis().Set(context.Background(), countKey, string(encoded), time.Minute*10)
	}

	totalPage := int(math.Ceil(float64(count) / float64(limit)))
	pageInfo := &pkg.PageInfo{
		CurrentPage: page,
		NextPage:    page + 1,
		PrevPage:    page - 1,
		TotalPage:   totalPage,
		TotalData:   count,
	}
	if page >= totalPage {
		pageInfo.NextPage = 0
	}
	if page <= 1 {
		pageInfo.PrevPage = 0
	}

	response.GetAllSuccess("Success Get All Movies", movies, pageInfo)
}

// GetMovieById godoc
// @Summary Get movie by ID
// @Schemes
// @Description Get details of a specific movie by its ID
// @Tags Movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} pkg.Response{data=models.MovieDetails}
// @Failure 400 {object} pkg.Response{error=string} "Invalid input"
// @Failure 404 {object} pkg.Response{error=string} "Movie not found"
// @Failure 500 {object} pkg.Response{error=string} "Internal server error"
// @Security ApiKeyAuth
// @Router /movies/{id} [get]
func GetMovieById(c *gin.Context) {
	response := pkg.NewResponse(c)
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}
	movie, err := repositories.FindMovieById(id)
	if err != nil {
		response.InternalServerError("Failed to fetch movie", err.Error())
		return
	}
	if movie == nil {
		response.NotFound(fmt.Sprintf("Movie with ID %d not found", id), nil)
		return
	}

	response.Success("Success get movie", movie)
}

// CreateMovie godoc
// @Summary Create a new movie
// @Schemes
// @Description Create a new movie with the provided details
// @Tags Movies
// @Accept json
// @Produce json
// @Param movie body models.MovieDetails true "Movie details"
// @Success 201 {object} pkg.Response{data=models.MovieDetails}
// @Failure 400 {object} pkg.Response{error=string}
// @Failure 500 {object} pkg.Response{error=string}
// @Security ApiKeyAuth
// @Router /movies [post]
func CreateMovie(c *gin.Context) {
}

// UpdateMovie godoc
// @Summary Update an existing movie
// @Schemes
// @Description Update the details of a movie by its ID
// @Tags Movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Param movie body models.MovieDetails true "Updated movie details"
// @Success 200 {object} pkg.Response{data=models.MovieDetails}
// @Failure 400 {object} pkg.Response{error=string}
// @Failure 404 {object} pkg.Response{error=string}
// @Failure 500 {object} pkg.Response{error=string}
// @Security ApiKeyAuth
// @Router /movies/{id} [patch]
func Update(c *gin.Context) {
}

// DeleteMovie godoc
// @Summary Delete a movie
// @Schemes
// @Description Delete a movie by its ID
// @Tags Movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Success 204 {object} pkg.Response
// @Failure 400 {object} pkg.Response{error=string}
// @Failure 404 {object} pkg.Response{error=string}
// @Failure 500 {object} pkg.Response{error=string}
// @Security ApiKeyAuth
// @Router /movies/{id} [delete]
func Delete(c *gin.Context) {
}
