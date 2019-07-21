package api

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-ggz/ui/dist"
)

// Favicon represents the favicon.
func Favicon(c *gin.Context) {
	file, _ := dist.ReadFile("favicon.ico")
	etag := fmt.Sprintf("%x", md5.Sum(file))
	c.Header("ETag", etag)
	c.Header("Cache-Control", "max-age=0")

	if match := c.GetHeader("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			c.Status(http.StatusNotModified)
			return
		}
	}

	c.Data(
		http.StatusOK,
		"image/x-icon",
		file,
	)
}
