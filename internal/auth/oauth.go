package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cli/browser"
	"github.com/cli/oauth"
	"github.com/cli/oauth/device"
)

const (
	gitHubDeviceFlowURL = "https://github.com/login/device/code"
	gitHubTokenURL      = "https://github.com/login/oauth/access_token"
	gitHubAPIURL        = "https://api.github.com"
)

func DeviceFlowLogin(ctx context.Context, clientID string, scopes []string) (string, error) {
	if clientID == "" {
		clientID = "Iv1.5ad3b2eda3f1e2b7"
	}

	host := oauth.GitHubHost("https://github.com")

	code, err := device.RequestCode(http.DefaultClient, host.DeviceCodeURL, clientID, scopes)
	if err != nil {
		return "", fmt.Errorf("failed to request device code: %w", err)
	}

	fmt.Printf("First, copy your one-time code: %s\n", code.UserCode)
	fmt.Printf("Then visit %s in your browser and paste the code.\n", code.VerificationURI)
	fmt.Println("(You can press Ctrl+C to cancel)")

	if err := browser.OpenURL(code.VerificationURI); err != nil {
		return "", fmt.Errorf("failed to open browser: %w", err)
	}

	token, err := device.Wait(ctx, http.DefaultClient, host.TokenURL, device.WaitOptions{
		ClientID:   clientID,
		DeviceCode: code,
	})
	if err != nil {
		return "", fmt.Errorf("device flow wait failed: %w", err)
	}

	return token.Token, nil
}
