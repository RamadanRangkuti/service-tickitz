package handlers

import (
	"RamadanRangkuti/service-tickitz/internal/models"
	"RamadanRangkuti/service-tickitz/pkg"

	"github.com/gin-gonic/gin"
)

func GetTes(c *gin.Context) {
	response := pkg.NewResponse(c)
	tes, err := models.FindALlTes()
	if err != nil {
		response.InternalServerError("Failed get all", err.Error())
		return
	}
	response.Success("Success Get All User", tes)
}
