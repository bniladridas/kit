package api

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLiveIntegrationGetUser(t *testing.T) {
	if os.Getenv("KIT_INTEGRATION_TESTS") != "1" {
		t.Skip("skipping live integration test; set KIT_INTEGRATION_TESTS=1 to run")
	}

	token := os.Getenv("KIT_GITHUB_TOKEN")
	if token == "" {
		t.Skip("skipping live integration test; set KIT_GITHUB_TOKEN to run")
	}

	client := NewClient("https://api.github.com/", token)

	resp, err := client.Get(context.Background(), "/user")
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}
	_ = resp

	var user map[string]interface{}
	if err := resp.Decode(&user); err != nil {
		t.Fatalf("failed to decode user: %v", err)
	}

	login, ok := user["login"].(string)
	assert.True(t, ok, "login should be a string")
	assert.NotEmpty(t, login, "login should not be empty")
	t.Logf("Authenticated as: %s", login)
}

func TestLiveIntegrationListRepos(t *testing.T) {
	if os.Getenv("KIT_INTEGRATION_TESTS") != "1" {
		t.Skip("skipping live integration test; set KIT_INTEGRATION_TESTS=1 to run")
	}

	token := os.Getenv("KIT_GITHUB_TOKEN")
	if token == "" {
		t.Skip("skipping live integration test; set KIT_GITHUB_TOKEN to run")
	}

	client := NewClient("https://api.github.com/", token)

	resp, err := client.Get(context.Background(), "/user/repos", WithQueryParams(map[string]string{
		"per_page": "1",
	}))
	if err != nil {
		t.Fatalf("failed to list repos: %v", err)
	}
	_ = resp

	var repos []map[string]interface{}
	if err := resp.Decode(&repos); err != nil {
		t.Fatalf("failed to decode repos: %v", err)
	}

	assert.NotEmpty(t, repos, "should have at least one repo")
	t.Logf("Found %d repos", len(repos))
}
