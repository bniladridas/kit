package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	ErrNotInRepo = fmt.Errorf("not in a git repository")
)

type Context struct {
	Owner  string
	Repo   string
	Branch string
	Remote string
	Root   string
}

var remotePattern = regexp.MustCompile(`(?:https://github\.com/|git@github\.com:)([^/]+)/([^/.]+)(?:\.git)?`)

func GetContext() (*Context, error) {
	root, err := findGitRoot()
	if err != nil {
		return nil, err
	}

	branch, err := getCurrentBranch(root)
	if err != nil {
		return nil, err
	}

	remote, owner, repo, err := getRemoteInfo(root)
	if err != nil {
		return nil, err
	}

	return &Context{
		Owner:  owner,
		Repo:   repo,
		Branch: branch,
		Remote: remote,
		Root:   root,
	}, nil
}

func findGitRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		gitDir := filepath.Join(dir, ".git")
		if info, err := os.Stat(gitDir); err == nil {
			if info.IsDir() {
				return dir, nil
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", ErrNotInRepo
		}
		dir = parent
	}
}

func getCurrentBranch(gitRoot string) (string, error) {
	headPath := filepath.Join(gitRoot, ".git", "HEAD")
	data, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}

	content := string(data)
	if strings.HasPrefix(content, "ref: ") {
		ref := strings.TrimSpace(content[5:])
		parts := strings.Split(ref, "/")
		return parts[len(parts)-1], nil
	}

	return "(detached)", nil
}

func getRemoteInfo(gitRoot string) (string, string, string, error) {
	configPath := filepath.Join(gitRoot, ".git", "config")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", "", "", err
	}

	content := string(data)

	remote, err := getCurrentRemote(content)
	if err != nil {
		return "", "", "", err
	}

	owner, repo, err := parseRemoteURL(remote)
	if err != nil {
		return "", "", "", err
	}

	return remote, owner, repo, nil
}

func getCurrentRemote(configContent string) (string, error) {
	inRemote := false
	currentRemote := ""

	for _, line := range strings.Split(configContent, "\n") {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "[remote ") {
			inRemote = true
			currentRemote = ""
			continue
		}

		if strings.HasPrefix(line, "[") {
			inRemote = false
			continue
		}

		if inRemote && strings.HasPrefix(line, "url =") {
			currentRemote = strings.TrimSpace(strings.TrimPrefix(line, "url ="))
			return currentRemote, nil
		}
	}

	return "", fmt.Errorf("no remote URL found in git config")
}

func parseRemoteURL(remote string) (string, string, error) {
	matches := remotePattern.FindStringSubmatch(remote)
	if matches == nil || len(matches) < 3 {
		return "", "", fmt.Errorf("not a GitHub remote: %s", remote)
	}

	return matches[1], matches[2], nil
}

func RunGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err, string(out))
	}
	return string(out), nil
}
