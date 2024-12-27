package handlers

import (
	"RamadanRangkuti/service-tickitz/internal/models"
	"RamadanRangkuti/service-tickitz/internal/repository"
	"RamadanRangkuti/service-tickitz/pkg"
	"context"
	"encoding/json"
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
	var movies models.ListMovies
	var count int

	get := pkg.Redis().Get(context.Background(), c.Request.RequestURI)
	if get.Val() != "" {
		rawData := []byte(get.Val())
		json.Unmarshal(rawData, &movies)
	} else {
		movies, err := repository.FindAllMovies(page, limit, order, sortBy, search)
		if err != nil {
			response.InternalServerError("Failed to fetch movies", err.Error())
			return
		}
		encoded, _ := json.Marshal(movies)
		pkg.Redis().Set(context.Background(), c.Request.RequestURI, string(encoded), 0)
	}

	getCount := pkg.Redis().Get(context.Background(), fmt.Sprintf("count+%s", c.Request.RequestURI))

	if getCount.Val() != "" {
		rawData := []byte(getCount.Val())
		json.Unmarshal(rawData, &count)
	} else {
		count = repository.CountMovie(search)
		encoded, _ := json.Marshal(count)
		pkg.Redis().Set(context.Background(), fmt.Sprintf("count+%s", c.Request.RequestURI), string(encoded), 0)
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
