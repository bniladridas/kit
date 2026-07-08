package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationGetUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/user", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"login": "octocat",
			"name":  "The Octocat",
		})
	}))
	defer server.Close()

	client := NewClient(server.URL+"/", "test-token")

	resp, err := client.Get(context.Background(), "/user")
	assert.NoError(t, err)

	var user map[string]interface{}
	err = resp.Decode(&user)
	assert.NoError(t, err)

	assert.Equal(t, "octocat", user["login"])
	assert.Equal(t, "The Octocat", user["name"])
}

func TestIntegrationListRepos(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/user/repos", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"full_name": "octocat/hello-world", "visibility": "public"},
			{"full_name": "octocat/hello-world-2", "visibility": "private"},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL+"/", "test-token")

	resp, err := client.Get(context.Background(), "/user/repos")
	assert.NoError(t, err)

	var repos []map[string]interface{}
	err = resp.Decode(&repos)
	assert.NoError(t, err)

	assert.Len(t, repos, 2)
	assert.Equal(t, "octocat/hello-world", repos[0]["full_name"])
	assert.Equal(t, "private", repos[1]["visibility"])
}

func TestIntegrationCreateIssue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repos/owner/repo/issues", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "Test Issue", body["title"])

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"number":  1,
			"title":   "Test Issue",
			"html_url": "https://github.com/owner/repo/issues/1",
		})
	}))
	defer server.Close()

	client := NewClient(server.URL+"/", "test-token")

	resp, err := client.Post(context.Background(), "/repos/owner/repo/issues", map[string]interface{}{
		"title": "Test Issue",
	})
	assert.NoError(t, err)

	var issue map[string]interface{}
	err = resp.Decode(&issue)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), issue["number"])
	assert.Equal(t, "Test Issue", issue["title"])
}

func TestIntegrationUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Bad credentials",
		})
	}))
	defer server.Close()

	client := NewClient(server.URL+"/", "bad-token")

	_, err := client.Get(context.Background(), "/user")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}
