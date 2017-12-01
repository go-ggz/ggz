package model

import (
	"fmt"
	"time"

	"github.com/go-ggz/ggz/module/meta"

	"github.com/appleboy/com/random"
)

// Shorten shortener URL
type Shorten struct {
	Slug        string    `xorm:"pk VARCHAR(14)" json:"slug"`
	UserID      int64     `xorm:"INDEX" json:"-"`
	User        *User     `xorm:"-" json:"user"`
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
	shorten := new(Shorten)
	has, err := x.
		Where("url = ?", url).
		Get(shorten)

	if err != nil {
		return nil, err
	}

	// get user data
	if shorten.UserID != 0 {
		u, _ := GetUserByID(shorten.UserID)
		if u != nil {
			shorten.User = u
		}
	}

	if has {
		return shorten, ErrURLExist{shorten.Slug, url}
	}

	return nil, nil
}

// NewShortenURL create url item
func NewShortenURL(url string, size int, user *User) (_ *Shorten, err error) {
	row := &Shorten{
		Date: time.Now(),
		URL:  url,
	}
	exists := true
	slug := ""

	if user != nil {
		row.UserID = user.ID
		row.User = user
	}

	for exists == true {
		slug = random.String(size)
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

func (s *Shorten) getUser(e Engine) (err error) {
	if s.User != nil {
		return nil
	}

	s.User, err = getUserByID(e, s.UserID)
	return err
}

// GetUser returns the shorten owner
func (s *Shorten) GetUser() error {
	return s.getUser(x)
}
