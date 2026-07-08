package auth

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginAndIsLoggedIn(t *testing.T) {
	Reset()
	tmpDir := t.TempDir()
	testStore := NewSecureStore("kit-test")
	testStore.SetFallbackDir(filepath.Join(tmpDir, ".kit"))
	testStore.SetForceFallback(true)

	SetSecureStoreForTest(testStore)

	err := Login(ProviderGitHub, "test-token-123")
	assert.NoError(t, err)

	assert.True(t, IsLoggedIn(ProviderGitHub))

	Reset()
	SetSecureStoreForTest(NewSecureStore("kit"))
}

func TestIsLoggedInFalse(t *testing.T) {
	Reset()
	tmpDir := t.TempDir()
	testStore := NewSecureStore("kit-test")
	testStore.SetFallbackDir(filepath.Join(tmpDir, ".kit"))
	testStore.SetForceFallback(true)

	SetSecureStoreForTest(testStore)
	defer func() {
		Reset()
		SetSecureStoreForTest(NewSecureStore("kit"))
	}()

	assert.False(t, IsLoggedIn(ProviderGitHub))
}

func TestLogout(t *testing.T) {
	Reset()
	tmpDir := t.TempDir()
	testStore := NewSecureStore("kit-test")
	testStore.SetFallbackDir(filepath.Join(tmpDir, ".kit"))
	testStore.SetForceFallback(true)

	SetSecureStoreForTest(testStore)

	err := Login(ProviderGitHub, "test-token")
	assert.NoError(t, err)
	assert.True(t, IsLoggedIn(ProviderGitHub))

	err = Logout(ProviderGitHub)
	assert.NoError(t, err)
	assert.False(t, IsLoggedIn(ProviderGitHub))

	Reset()
	SetSecureStoreForTest(NewSecureStore("kit"))
}
