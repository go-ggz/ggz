package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeMD5(t *testing.T) {
	assert.Equal(t,
		"3858f62230ac3c915f300c664312c63f",
		EncodeMD5("foobar"),
	)
}

func TestEncodeSha1(t *testing.T) {
	assert.Equal(t,
		"8843d7f92416211de9ebb963ff4ce28125932878",
		EncodeSha1("foobar"),
	)
}
