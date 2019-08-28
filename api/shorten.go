package api

import (
	"net/http"
	"regexp"

	"github.com/go-ggz/ggz/pkg/config"
	"github.com/go-ggz/ggz/pkg/helper"
	"github.com/go-ggz/ggz/pkg/model"
	"github.com/go-ggz/ggz/pkg/router"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var shortenPattern = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

// ShortenedIndex index page.
func ShortenedIndex(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, config.Server.Host)
}

// FormURL URL Struct
type FormURL struct {
	URL string `json:"url,omitempty" binding:"required,url"`
}

// CreateShortenURL create shorten url
func CreateShortenURL(c *gin.Context) {
	var data FormURL
	if err := c.ShouldBindJSON(&data); err != nil {
		errorJSON(c, http.StatusBadRequest, errBadRequest)
		return
	}

	row, err := model.GetShortenFromURL(data.URL)

	if model.IsErrURLExist(err) {
		c.JSON(
			http.StatusOK,
			row,
		)
		return
	}

	if err != nil {
		errorJSON(c, http.StatusInternalServerError, errInternalServer)
		return
	}

	user := helper.GetUserDataFromModel(c.Request.Context())
	row, err = model.CreateShorten(data.URL, config.Server.ShortenSize, user)

	if err != nil {
		errorJSON(c, http.StatusInternalServerError, errInternalServer)
		return
	}

	// upload QRCode image.
	go func(slug string) {
		if err := helper.QRCodeGenerator(slug); err != nil {
			log.Error().Err(err).Msg("QRCode Generator fail")
		}
	}(row.Slug)

	c.JSON(
		http.StatusOK,
		row,
	)
}

// FetchShortenedURL show URL content
func FetchShortenedURL(c *gin.Context) {
	slug := c.Param("slug")

	if !shortenPattern.MatchString(slug) {
		errorJSON(c, http.StatusBadRequest, errSlugNotMatch)
		return
	}

	row, err := model.GetShortenBySlug(slug)
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, errInternalServer)
		return
	}

	if model.IsErrShortenNotExist(err) {
		errorJSON(c, http.StatusNotFound, errSlugNotFound)
		return
	}

	c.JSON(
		http.StatusOK,
		row,
	)
}

// RedirectURL redirect shorten text to origin URL.
func RedirectURL(c *gin.Context) {
	slug := c.Param("slug")

	if !shortenPattern.MatchString(slug) {
		errorJSON(c, http.StatusBadRequest, errSlugNotMatch)
		return
	}

	if slug == "healthz" {
		Heartbeat(c)
		return
	} else if slug == "metrics" {
		if config.Metrics.Enabled {
			router.Metrics(config.Metrics.Token)(c)
		}
		return
	}

	row, err := model.GetShortenBySlug(slug)
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, errInternalServer)
		return
	}

	if model.IsErrShortenNotExist(err) {
		errorJSON(c, http.StatusNotFound, errSlugNotFound)
		return
	}

	err = row.UpdateHits(slug)
	if err != nil {
		errorJSON(c, http.StatusNotFound, errInternalServer)
		return
	}

	c.Redirect(http.StatusMovedPermanently, row.URL)
}
