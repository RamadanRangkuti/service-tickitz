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

		claims, err := pkg.VerifyToken(token)
		if err != nil {
			response.Unauthorized("Unauthorized", nil)
		}
		c.Set("UserId", claims.UserId)
		c.Set("UserRole", claims.UserRole)
		c.Next()
	}
}

func RoleCheck(requiredRole int) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := pkg.NewResponse(c)

		role, exists := c.Get("UserRole")
		if !exists {
			//forbidden
			response.Forbidden("Access denied", nil)
			return
		}

		userRole, ok := role.(int)
		if !ok {
			response.InternalServerError("Failed to parse user role from token", nil)
			return
		}

		if userRole != requiredRole {
			response.Forbidden("Access denied", nil)
			return
		}

		c.Next()
	}
}
