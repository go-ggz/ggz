package assets

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/go-ggz/ggz/config"

	"github.com/appleboy/com/file"
	"github.com/sirupsen/logrus"
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
			logrus.Warnf("Custom assets directory doesn't exist")
		}
	}

	f, err := FS.OpenFile(CTX, origPath, os.O_RDONLY, 0644)

	if err != nil {
		return nil, err
	}

	return f, nil
}

// ReadSource is adapTed from ioutil
func ReadSource(origPath string) (content []byte, err error) {
	content, err = ReadFile(origPath)

	if err != nil {
		logrus.Warnf("Failed to read builtin %s file. %s", origPath, err)
	}

	if config.Server.Assets != "" && file.IsDir(config.Server.Assets) {
		origPath = path.Join(
			config.Server.Assets,
			origPath,
		)

		if file.IsFile(origPath) {
			content, err = ioutil.ReadFile(origPath)

			if err != nil {
				logrus.Warnf("Failed to read custom %s file. %s", origPath, err)
			}
		}
	}

	return content, err
}
