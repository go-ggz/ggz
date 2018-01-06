package model

import (
	"time"

	"github.com/go-ggz/ggz/module/base"

	"github.com/satori/go.uuid"
)

type AccessToken struct {
	ID     int64 `xorm:"pk autoincr"`
	UserID int64 `xorm:"INDEX"`
	Name   string
	Sha1   string `xorm:"UNIQUE VARCHAR(40)"`

	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
	HasRecentActivity bool      `xorm:"-"`
	HasUsed           bool      `xorm:"-"`
}

// AfterLoad is invoked from XORM after setting the values of all fields of this object.
func (t *AccessToken) AfterLoad() {
	t.HasUsed = t.UpdatedAt.Unix() > t.CreatedAt.Unix()
}

// NewAccessToken creates new access token.
func NewAccessToken(t *AccessToken) error {
	t.Sha1 = base.EncodeSha1(uuid.NewV4().String())
	_, err := x.Insert(t)
	return err
}

// GetAccessTokenBySHA returns access token by given sha1.
func GetAccessTokenBySHA(sha string) (*AccessToken, error) {
	if sha == "" {
		return nil, ErrAccessTokenEmpty{}
	}
	t := &AccessToken{Sha1: sha}
	has, err := x.Get(t)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrAccessTokenNotExist{sha}
	}
	return t, nil
}

// UpdateAccessToken updates information of access token.
func UpdateAccessToken(t *AccessToken) error {
	_, err := x.ID(t.ID).AllCols().Update(t)
	return err
}

// DeleteAccessTokenByID deletes access token by given ID.
func DeleteAccessTokenByID(id, userID int64) error {
	cnt, err := x.ID(id).Delete(&AccessToken{
		UserID: userID,
	})
	if err != nil {
		return err
	} else if cnt != 1 {
		return ErrAccessTokenNotExist{}
	}
	return nil
}
