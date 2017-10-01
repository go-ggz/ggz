package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index represents the index page.
func Index(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":  http.StatusOK,
			"error": "Welcome gzz shorten URL server",
		},
	)
}
