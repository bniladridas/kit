package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContext(t *testing.T) {
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	os.MkdirAll(gitDir, 0o755)

	os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/main\n"), 0o644)
	os.WriteFile(filepath.Join(gitDir, "config"), []byte(`
[remote "origin"]
	url = https://github.com/owner/repo.git
`), 0o644)

	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	ctx, err := GetContext()
	assert.NoError(t, err)
	assert.Equal(t, "owner", ctx.Owner)
	assert.Equal(t, "repo", ctx.Repo)
	assert.Equal(t, "main", ctx.Branch)
	assert.Equal(t, "https://github.com/owner/repo.git", ctx.Remote)
	expectedRoot, _ := filepath.EvalSymlinks(tmpDir)
	actualRoot, _ := filepath.EvalSymlinks(ctx.Root)
	assert.Equal(t, expectedRoot, actualRoot)
}

func TestGetContextSSHRemote(t *testing.T) {
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	os.MkdirAll(gitDir, 0o755)

	os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/develop\n"), 0o644)
	os.WriteFile(filepath.Join(gitDir, "config"), []byte(`
[remote "upstream"]
	url = git@github.com:octocat/hello-world.git
`), 0o644)

	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	ctx, err := GetContext()
	assert.NoError(t, err)
	assert.Equal(t, "octocat", ctx.Owner)
	assert.Equal(t, "hello-world", ctx.Repo)
	assert.Equal(t, "develop", ctx.Branch)
	assert.Equal(t, "git@github.com:octocat/hello-world.git", ctx.Remote)
	expectedRoot, _ := filepath.EvalSymlinks(tmpDir)
	actualRoot, _ := filepath.EvalSymlinks(ctx.Root)
	assert.Equal(t, expectedRoot, actualRoot)
}

func TestGetContextNotInRepo(t *testing.T) {
	tmpDir := t.TempDir()

	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	_, err := GetContext()
	assert.ErrorIs(t, err, ErrNotInRepo)
}

func TestParseRemoteURL(t *testing.T) {
	tests := []struct {
		remote     string
		owner      string
		repo       string
		expectErr  bool
	}{
		{"https://github.com/owner/repo.git", "owner", "repo", false},
		{"https://github.com/owner/repo", "owner", "repo", false},
		{"git@github.com:owner/repo.git", "owner", "repo", false},
		{"git@github.com:owner/repo", "owner", "repo", false},
		{"https://gitlab.com/owner/repo.git", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.remote, func(t *testing.T) {
			owner, repo, err := parseRemoteURL(tt.remote)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.owner, owner)
				assert.Equal(t, tt.repo, repo)
			}
		})
	}
}
