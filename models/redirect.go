package models

import (
	"time"

	"github.com/appleboy/com/random"
)

// Redirect shortener URL
type Redirect struct {
	Slug string    `xorm:"pk VARCHAR(14)" json:"slug"`
	URL  string    `xorm:"NOT NULL VARCHAR(620)" json:"url"`
	Date time.Time `json:"date"`
	Hits int64     `xorm:"NOT NULL DEFAULT 0" json:"hits"`
}

// GetFromSlug get shorten URL data
func (shorten *Redirect) GetFromSlug(slug string) (bool, error) {
	return x.
		Where("slug = ?", slug).
		Get(shorten)
}

// GetShortenFromURL check url exist
func GetShortenFromURL(url string) (*Redirect, error) {
	var data Redirect
	has, err := x.
		Where("url = ?", url).
		Get(&data)

	if err != nil {
		return nil, err
	}

	if has {
		return nil, ErrURLExist{data.Slug, url}
	}

	return &data, nil
}

// NewShortenURL create url item
func NewShortenURL(url string) (_ *Redirect, err error) {
	row := &Redirect{
		Date: time.Now(),
		URL:  url,
	}
	exists := true
	slug := ""

	for exists == true {
		slug = random.String(5)
		exists, err = row.GetFromSlug(slug)
		if err != nil {
			return nil, err
		}
	}

	row.Slug = slug

	if _, err := x.Insert(row); err != nil {
		return nil, err
	}

	return row, nil
}

// UpdateHits udpate hit count
func (shorten *Redirect) UpdateHits(slug string) error {
	if _, err := x.Exec("UPDATE `redirect` SET hits = hits + 1 WHERE slug = ?", slug); err != nil {
		return err
	}

	return nil
}
