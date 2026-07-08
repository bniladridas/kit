package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewClient("https://api.example.com", "test-token")

	assert.NotNil(t, client)
	assert.Equal(t, "https://api.example.com/", client.BaseURL)
	assert.Equal(t, "test-token", client.Token)
}

func TestNewClientEmptyValues(t *testing.T) {
	client := NewClient("", "")

	assert.NotNil(t, client)
	assert.Equal(t, "/", client.BaseURL)
	assert.Empty(t, client.Token)
}
