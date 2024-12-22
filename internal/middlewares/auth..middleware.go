package middlewares

import (
	"RamadanRangkuti/service-tickitz/pkg"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := pkg.NewResponse(c)

		head := c.GetHeader("Authorization")
		if head == "" {
			response.Unauthorized("Unauthorized", nil)
			return
		}
		token := strings.Split(head, " ")[1]

		if token == "" {
			response.Unauthorized("Unauthorized", nil)
		}

		_, err := pkg.VerifyToken(token)
		if err != nil {
			response.Unauthorized("Unauthorized", nil)
		}
		c.Next()
	}
}
