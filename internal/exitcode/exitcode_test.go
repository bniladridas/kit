package exitcode

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodedError(t *testing.T) {
	err := New(ErrAPI, errors.New("api failed"))

	assert.Equal(t, ErrAPI, Code(err))
	assert.True(t, Is(err))
	assert.Equal(t, "api failed", err.Error())
}

func TestCodedErrorUnwrap(t *testing.T) {
	inner := errors.New("inner")
	err := New(ErrNetwork, inner)

	assert.Equal(t, inner, err.Unwrap())
	assert.True(t, errors.Is(err, inner))
}

func TestCodePlainError(t *testing.T) {
	err := errors.New("plain error")

	assert.False(t, Is(err))
	assert.Equal(t, ErrGeneral, Code(err))
}
