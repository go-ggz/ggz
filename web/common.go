package web

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func errorJSON(c *gin.Context, code int, err InnError) {
	logrus.Errorln(err.Error())

	c.AbortWithStatusJSON(
		code,
		gin.H{
			"code":  err.Code,
			"error": err.Error(),
		},
	)
}
