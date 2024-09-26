package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (cre *Credential) AuthMiddlewareGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", -1)

		err := cre.VerifyToken(token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}
