package web

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	m "github.com/keighl/metabolize"
)

// MetaData from URL
type MetaData struct {
	Title       string    `meta:"og:title,title" json:"title"`
	Description string    `meta:"og:description,description" json:"description"`
	Type        string    `meta:"og:type" json:"type"`
	URL         url.URL   `meta:"og:url" json:"url"`
	Image       string    `meta:"og:image" json:"image"`
	Time        time.Time `meta:"article:published_time,parsely-pub-date" json:"time"`
	VideoWidth  int64     `meta:"og:video:width" json:"video_width"`
	VideoHeight int64     `meta:"og:video:height" json:"video_height"`
}

// URLMeta for fetch metadata from URL
func URLMeta(c *gin.Context) {
	var data FormURL
	if err := c.ShouldBindWith(&data, binding.JSON); err != nil {
		errorJSON(c, http.StatusBadRequest, errBadRequest)
		return
	}

	res, _ := http.Get(data.URL)

	meta := new(MetaData)

	if err := m.Metabolize(res.Body, meta); err != nil {
		errorJSON(c, http.StatusInternalServerError, errInternalServer)
		return
	}

	c.JSON(
		http.StatusOK,
		meta,
	)
}
