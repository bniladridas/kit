package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/zalando/go-keyring"
)

var (
	ErrKeyringUnavailable = errors.New("keyring unavailable")
)

type SecureStore struct {
	ServiceName         string
	fallbackDirOverride string
	forceFallback       bool
}

func NewSecureStore(service string) *SecureStore {
	return &SecureStore{ServiceName: service}
}

func (s *SecureStore) SetFallbackDir(dir string) {
	s.fallbackDirOverride = dir
}

func (s *SecureStore) SetForceFallback(force bool) {
	s.forceFallback = force
}

func (s *SecureStore) fallbackDir() (string, error) {
	if s.fallbackDirOverride != "" {
		return s.fallbackDirOverride, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".kit", "credentials"), nil
}

func (s *SecureStore) Save(provider Provider, token string) error {
	if err := keyring.Set(s.ServiceName, string(provider), token); err != nil {
		if errors.Is(err, keyring.ErrUnsupportedPlatform) {
			return s.saveFallback(provider, token)
		}
		return fmt.Errorf("keyring save failed: %w", err)
	}
	return nil
}

func (s *SecureStore) Load(provider Provider) (string, error) {
	if s.forceFallback {
		return s.loadFallback(provider)
	}

	token, err := keyring.Get(s.ServiceName, string(provider))
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return "", ErrNotAuthenticated
		}
		if errors.Is(err, keyring.ErrUnsupportedPlatform) {
			return s.loadFallback(provider)
		}
		return "", fmt.Errorf("keyring load failed: %w", err)
	}
	return token, nil
}

func (s *SecureStore) Delete(provider Provider) error {
	if s.forceFallback {
		return s.deleteFallback(provider)
	}

	if err := keyring.Delete(s.ServiceName, string(provider)); err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return nil
		}
		if errors.Is(err, keyring.ErrUnsupportedPlatform) {
			return s.deleteFallback(provider)
		}
		return fmt.Errorf("keyring delete failed: %w", err)
	}
	return nil
}

func (s *SecureStore) saveFallback(provider Provider, token string) error {
	dir, err := s.fallbackDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	data, err := json.Marshal(map[string]string{
		"token": token,
	})
	if err != nil {
		return err
	}

	path := filepath.Join(dir, string(provider)+".json")
	return os.WriteFile(path, data, 0o600)
}

func (s *SecureStore) loadFallback(provider Provider) (string, error) {
	dir, err := s.fallbackDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(dir, string(provider)+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", ErrNotAuthenticated
	}

	var creds map[string]string
	if err := json.Unmarshal(data, &creds); err != nil {
		return "", err
	}

	token, ok := creds["token"]
	if !ok || token == "" {
		return "", ErrNotAuthenticated
	}

	return token, nil
}

func (s *SecureStore) deleteFallback(provider Provider) error {
	dir, err := s.fallbackDir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, string(provider)+".json")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func IsKeyringAvailable() bool {
	switch runtime.GOOS {
	case "darwin", "windows", "linux":
		return true
	default:
		return false
	}
}
