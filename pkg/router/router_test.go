package router

import (
	"net/http"
	"testing"

	"github.com/go-ggz/ggz/pkg/model"

	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"
)

func TestHealthzOnRedirectService(t *testing.T) {
	assert.NoError(t, model.PrepareTestDatabase())

	r := gofight.New()

	t.Run("return 200", func(t *testing.T) {
		r.GET("/healthz").
			Run(LoadRedirct(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusOK, r.Code)
				assert.Equal(t, "ok", r.Body.String())
			})
	})
}
