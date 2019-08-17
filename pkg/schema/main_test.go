package schema

import (
	"testing"

	"github.com/go-ggz/ggz/pkg/config"
	"github.com/go-ggz/ggz/pkg/model"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

func TestMain(m *testing.M) {
	if err := envconfig.Process("GGZ", config.Server); err != nil {
		log.Fatal().Err(err).Msg("can't load server config")
	}

	model.MainTest(m, "..")
}
