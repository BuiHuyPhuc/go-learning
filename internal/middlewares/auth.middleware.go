package middlewares

import (
	"context"
	"go-learning/internal/utils/auth"
	"log"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	// return func(c *gin.Context) {
	// 	// TODO: Implement authentication logic
	// 	token := c.GetHeader("Authorization")
	// 	if token != "valid-token" {
	// 		response.ErrorResponse(c, response.ErrInvalidToken, "")
	// 		c.Abort()
	// 		return
	// 	}

	// 	c.Next()
	// }
	return func(c *gin.Context) {
		//get the request url path
		uri := c.Request.URL.Path
		log.Println("uri request: ", uri)

		// check headers authorization
		jwtToken, valid := auth.ExtractBearerToken(c)
		if !valid {
			c.AbortWithStatusJSON(401, gin.H{"code": 40001, "err": "Unauthorized", "description": ""})
			return
		}

		// validate token
		claims, err := auth.VerifyTokenSubject(jwtToken)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"code": 40001, "err": "Invalid token", "description": ""})
			return
		}
		log.Println("claims::: UUID::", claims.Subject)

		// update claims to context
		ctx := context.WithValue(c.Request.Context(), "subjectUUID", claims.Subject)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
