package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-ggz/ui/dist"
)

// Favicon represents the favicon.
func Favicon(c *gin.Context) {
	file, _ := dist.ReadFile("favicon.ico")

	c.Data(
		http.StatusOK,
		"image/x-icon",
		file,
	)
}
