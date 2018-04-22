package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUserExist(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())
	exists, err := IsUserExist(0, "test@gmail.com")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = IsUserExist(0, "test123456@gmail.com")
	assert.NoError(t, err)
	assert.False(t, exists)

	exists, err = IsUserExist(1, "test1234@gmail.com")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = IsUserExist(1, "test123456@gmail.com")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestGetUserByEmail(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())

	t.Run("missing email", func(t *testing.T) {
		user, err := GetUserByEmail("")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.True(t, IsErrUserNotExist(err))
	})

	t.Run("test exist email", func(t *testing.T) {
		user, err := GetUserByEmail("test@gmail.com")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, int64(1), user.ID)
	})

	t.Run("email not found", func(t *testing.T) {
		user, err := GetUserByEmail("test123456@gmail.com")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.True(t, IsErrUserNotExist(err))
	})
}
