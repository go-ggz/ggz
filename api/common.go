package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func errorJSON(c *gin.Context, code int, err InnError) {
	log.Error().Err(err).Msg("json error")

	c.AbortWithStatusJSON(
		code,
		gin.H{
			"code":  err.Code,
			"error": err.Error(),
		},
	)
}
