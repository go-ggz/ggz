package models

import (
	"fmt"
	"time"

	"github.com/go-ggz/ggz/config"
	"github.com/go-ggz/ggz/modules/meta"

	"github.com/appleboy/com/random"
)

// Shorten shortener URL
type Shorten struct {
	Slug        string    `xorm:"pk VARCHAR(14)" json:"slug"`
	URL         string    `xorm:"NOT NULL VARCHAR(620)" json:"url"`
	Date        time.Time `json:"date"`
	Hits        int64     `xorm:"NOT NULL DEFAULT 0" json:"hits"`
	Title       string    `xorm:"VARCHAR(512)"`
	Description string    `xorm:"TEXT"`
	Type        string
	Image       string
}

// GetFromSlug get shorten URL data
func (s *Shorten) GetFromSlug(slug string) (bool, error) {
	return x.
		Where("slug = ?", slug).
		Get(s)
}

// GetShortenFromURL check url exist
func GetShortenFromURL(url string) (*Shorten, error) {
	var data Shorten
	has, err := x.
		Where("url = ?", url).
		Get(&data)

	if err != nil {
		return nil, err
	}

	if has {
		return &data, ErrURLExist{data.Slug, url}
	}

	return &data, nil
}

// NewShortenURL create url item
func NewShortenURL(url string) (_ *Shorten, err error) {
	row := &Shorten{
		Date: time.Now(),
		URL:  url,
	}
	exists := true
	slug := ""

	for exists == true {
		slug = random.String(config.Server.ShortenSize)
		exists, err = row.GetFromSlug(slug)
		if err != nil {
			return nil, err
		}
	}

	row.Slug = slug

	if _, err := x.Insert(row); err != nil {
		return nil, err
	}

	go row.UpdateMetaData()

	return row, nil
}

// UpdateHits udpate hit count
func (s *Shorten) UpdateHits(slug string) error {
	if _, err := x.Exec("UPDATE `shorten` SET hits = hits + 1 WHERE slug = ?", slug); err != nil {
		return err
	}

	return nil
}

// UpdateMetaData form raw body
func (s *Shorten) UpdateMetaData() error {
	data, err := meta.FetchData(s.URL)

	if err != nil {
		return err
	}

	s.Title = data.Title
	s.Description = data.Description
	s.Type = data.Type
	s.Image = data.Image

	if _, err := x.ID(s.Slug).Update(s); err != nil {
		return fmt.Errorf("update shorten [%s]: %v", s.Slug, err)
	}

	return nil
}
