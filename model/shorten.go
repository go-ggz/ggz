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
	URL         string    `xorm:"NOT NULL VARCHAR(620)" json:"url"`
	Date        time.Time `json:"date"`
	Hits        int64     `xorm:"NOT NULL DEFAULT 0" json:"hits"`
	Title       string    `xorm:"VARCHAR(512)" json:"title"`
	Description string    `xorm:"TEXT" json:"description"`
	Type        string    `json:"type"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `xorm:"created" json:"created_at,omitempty"`
	UpdatedAt   time.Time `xorm:"updated" json:"updated_at,omitempty"`

	// reference
	User *User `xorm:"-" json:"user"`
}

func getShortenBySlug(e Engine, slug string) (*Shorten, error) {
	s := new(Shorten)
	has, err := e.Where("slug = ?", slug).Get(s)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrShortenNotExist{slug}
	}
	return s, nil
}

// GetShortenBySlug returns the shorten object by given slug if exists.
func GetShortenBySlug(slug string) (*Shorten, error) {
	return getShortenBySlug(x, slug)
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
		_, err = getShortenBySlug(x, slug)
		if err != nil {
			if IsErrShortenNotExist(err) {
				exists = false
				continue
			}
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

// GetShortenURLs returns a list of urls of given user.
func GetShortenURLs(userID int64, page, pageSize int, orderBy string) ([]*Shorten, error) {
	sess := x.NewSession()

	if len(orderBy) == 0 {
		orderBy = "date DESC"
	}

	if userID != 0 {
		sess = sess.
			Where("user_id = ?", userID)
	}

	sess = sess.OrderBy(orderBy)

	if page <= 0 {
		page = 1
	}
	sess.Limit(pageSize, (page-1)*pageSize)

	urls := make([]*Shorten, 0, pageSize)
	return urls, sess.Find(&urls)
}
