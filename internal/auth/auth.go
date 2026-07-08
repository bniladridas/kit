package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	ErrNotAuthenticated = errors.New("not authenticated")
	ErrInvalidToken     = errors.New("invalid token")
)

type Provider string

const (
	ProviderGitHub Provider = "github"
)

var (
	secureStore *SecureStore
)

func init() {
	secureStore = NewSecureStore("kit")
}

func SetSecureStoreForTest(s *SecureStore) {
	secureStore = s
}

type Credentials struct {
	Provider Provider
	Token    string
}

type StoredCredentials struct {
	Credentials map[Provider]Credentials
}

var currentCredentials *StoredCredentials

func Login(provider Provider, token string) error {
	if token == "" {
		return ErrInvalidToken
	}

	if currentCredentials == nil {
		currentCredentials = &StoredCredentials{
			Credentials: make(map[Provider]Credentials),
		}
	}

	currentCredentials.Credentials[provider] = Credentials{
		Provider: provider,
		Token:    token,
	}

	if err := secureStore.Save(provider, token); err != nil {
		return err
	}

	return nil
}

func LoginWithOAuth(ctx context.Context, provider Provider, clientID string, scopes []string) (string, error) {
	switch provider {
	case ProviderGitHub:
		token, err := DeviceFlowLogin(ctx, clientID, scopes)
		if err != nil {
			return "", err
		}

		if err := Login(provider, token); err != nil {
			return "", err
		}

		return token, nil
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}
}

func Logout(provider Provider) error {
	if currentCredentials != nil {
		delete(currentCredentials.Credentials, provider)
	}

	return secureStore.Delete(provider)
}

func IsLoggedIn(provider Provider) bool {
	if currentCredentials != nil {
		if creds, ok := currentCredentials.Credentials[provider]; ok && creds.Token != "" {
			return true
		}
	}

	token, err := secureStore.Load(provider)
	if err != nil {
		return false
	}

	if currentCredentials == nil {
		currentCredentials = &StoredCredentials{
			Credentials: make(map[Provider]Credentials),
		}
	}

	currentCredentials.Credentials[provider] = Credentials{
		Provider: provider,
		Token:    token,
	}

	return true
}

func GetToken(provider Provider) (string, error) {
	if currentCredentials != nil {
		if creds, ok := currentCredentials.Credentials[provider]; ok && creds.Token != "" {
			return creds.Token, nil
		}
	}

	token, err := secureStore.Load(provider)
	if err != nil {
		return "", err
	}

	if currentCredentials == nil {
		currentCredentials = &StoredCredentials{
			Credentials: make(map[Provider]Credentials),
		}
	}

	currentCredentials.Credentials[provider] = Credentials{
		Provider: provider,
		Token:    token,
	}

	return token, nil
}

func ValidateToken(ctx context.Context, provider Provider) (string, error) {
	token, err := GetToken(provider)
	if err != nil {
		return "", err
	}

	return validateTokenWithAPI(ctx, token)
}

func validateTokenWithAPI(ctx context.Context, token string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, gitHubAPIURL+"/user", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "kit/0.1")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("invalid token")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	login, _ := result["login"].(string)
	return login, nil
}

func Reset() {
	currentCredentials = nil
}

func ParseProvider(s string) (Provider, error) {
	p := Provider(strings.ToLower(strings.TrimSpace(s)))
	switch p {
	case ProviderGitHub:
		return p, nil
	default:
		return "", fmt.Errorf("unsupported provider: %s", s)
	}
}
