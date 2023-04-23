package middleware

import (
	"MyGramAPI/app/entity"
	"MyGramAPI/pkg/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		verifyToken, err := helpers.VerifyToken(c)

		if err != nil {
			c.JSON(http.StatusUnauthorized, entity.Response{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		c.Set("userData", verifyToken)
		c.Next()
	}
}
