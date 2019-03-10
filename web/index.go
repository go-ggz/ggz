package web

import (
	"net/http"

	"github.com/go-ggz/ui/dist"

	"github.com/gin-gonic/gin"
)

// Index represents the index page.
func Index(c *gin.Context) {
	file, _ := dist.ReadFile("index.html")

	c.Data(http.StatusOK, "text/html; charset=utf-8", file)
}
