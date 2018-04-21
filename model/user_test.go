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
