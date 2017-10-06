package meta

import (
	"net/http"
	"net/url"
	"time"

	m "github.com/keighl/metabolize"
)

// MetaData struct
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

// FetchData for fetch metadata from header of body
func FetchData(url string) (MetaData, error) {
	res, _ := http.Get(url)

	meta := new(MetaData)

	if err := m.Metabolize(res.Body, meta); err != nil {
		return *meta, err
	}

	return *meta, nil
}
