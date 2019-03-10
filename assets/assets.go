package assets

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/go-ggz/ggz/config"
	"github.com/go-ggz/ui/dist"

	"github.com/appleboy/com/file"
	"github.com/rs/zerolog/log"
)

//go:generate fileb0x ab0x.yaml

// Load initializes the static files.
func Load() http.FileSystem {
	return ChainedFS{}
}

// ChainedFS is a simple HTTP filesystem including custom path.
type ChainedFS struct {
}

// Open just implements the HTTP filesystem interface.
func (c ChainedFS) Open(origPath string) (http.File, error) {
	if config.Server.Assets != "" {
		if file.IsDir(config.Server.Assets) {
			customPath := path.Join(config.Server.Assets, origPath)

			if file.IsFile(customPath) {
				f, err := os.Open(customPath)

				if err != nil {
					return nil, err
				}

				return f, nil
			}
		} else {
			log.Warn().Msg("Custom assets directory doesn't exist")
		}
	}

	log.Debug().Msgf("origPath is %s", origPath)

	filePath := fmt.Sprintf("/assets/%s", origPath)

	f, err := dist.FS.OpenFile(dist.CTX, filePath, os.O_RDONLY, 0644)

	if err != nil {
		return nil, err
	}

	return f, nil
}

// ReadSource is adapTed from ioutil
func ReadSource(origPath string) (content []byte, err error) {
	content, err = dist.ReadFile(origPath)

	if err != nil {
		log.Warn().Err(err).Msgf("Failed to read builtin %s file.", origPath)
	}

	if config.Server.Assets != "" && file.IsDir(config.Server.Assets) {
		origPath = path.Join(
			config.Server.Assets,
			origPath,
		)

		if file.IsFile(origPath) {
			content, err = ioutil.ReadFile(origPath)

			if err != nil {
				log.Warn().Err(err).Msgf("Failed to read custom %s file", origPath)
			}
		}
	}

	return content, err
}
