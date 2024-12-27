package handlers

import (
	"RamadanRangkuti/service-tickitz/internal/repository"
	"RamadanRangkuti/service-tickitz/pkg"
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

	movies, err := repository.FindAllMovies(page, limit, order, sortBy, search)
	if err != nil {
		response.InternalServerError("Failed to fetch movies", err.Error())
		return
	}

	if len(movies) == 0 {
		response.NotFound("No movies found", movies)
		return
	}
	count := repository.CountMovie(search)
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

func GetMovieById(c *gin.Context) {
	response := pkg.NewResponse(c)
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}
	movie, err := repository.FindMovieById(id)
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
