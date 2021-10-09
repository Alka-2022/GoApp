package middlewares

import (
	"net/http"

	"go-webapp/db"

	"github.com/gin-gonic/gin"
)


func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}


func Connect(c *gin.Context) {
	s := db.Conn.Session.Clone()

	defer s.Close()

	c.Set("db", s)
	c.Next()
}


func ErrorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		c.HTML(http.StatusBadRequest, "400", gin.H{
			"errors": c.Errors,
		})
	}
}
