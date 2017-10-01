package web

import (
	"net/http"

	"github.com/go-ggz/ggz/assets"

	"github.com/gin-gonic/gin"
)

// Favicon represents the favicon.
func Favicon(c *gin.Context) {
	file, _ := assets.ReadFile("favicon.ico")

	c.Data(
		http.StatusOK,
		"image/x-icon",
		file,
	)
}
