package api

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-ggz/ui/dist"

	"github.com/gin-gonic/gin"
)

// Index represents the index page.
func Index(c *gin.Context) {
	file, _ := dist.ReadFile("index.html")
	etag := fmt.Sprintf("%x", md5.Sum(file))
	c.Header("ETag", etag)
	c.Header("Cache-Control", "no-cache")

	if match := c.GetHeader("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			c.Status(http.StatusNotModified)
			return
		}
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", file)
}
