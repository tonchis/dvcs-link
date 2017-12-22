package main

import (
	"strings"
	"testing"
)

type MockBash struct {
}

func (m MockBash) Run(command string) string {
	commandParts := strings.SplitN(command, " ", 3)
	gitCommand := commandParts[1]
	gitArgs := commandParts[2]

	switch gitCommand {
	case "rev-parse":
		return "commit-sha"
	case "config":
		if strings.Contains(gitArgs, "origin") || strings.Contains(gitArgs, "gitlab") {
			return "git@gitlab.com:tonchis/dvcs-link"
		} else if strings.Contains(gitArgs, "github") {
			return "git@github.com:tonchis/dvcs-link.git"
		} else {
			return "https://usersname@bitbucket.org/username/reponame.git"
		}
	}

	return ""
}

func TestMain(m *testing.M) {
	bashExec = MockBash{}.Run
	m.Run()
}

func TestFilenameOnly(t *testing.T) {
	link := resolveLink("origin", "foo.go", "", "")
	expectedLink := "https://gitlab.com/tonchis/dvcs-link/blob/commit-sha/foo.go"

	if link != expectedLink {
		t.Errorf("Expected %v to equal %v", link, expectedLink)
	}
}

func TestWithLine(t *testing.T) {
	link := resolveLink("origin", "foo.go", "6", "")
	expectedLink := "https://gitlab.com/tonchis/dvcs-link/blob/commit-sha/foo.go#L6"

	if link != expectedLink {
		t.Errorf("Expected %v to equal %v", link, expectedLink)
	}
}

func TestWithGitLabRange(t *testing.T) {
	link := resolveLink("origin", "foo.go", "5", "10")
	expectedLink := "https://gitlab.com/tonchis/dvcs-link/blob/commit-sha/foo.go#L5-10"

	if link != expectedLink {
		t.Errorf("Expected %v to equal %v", link, expectedLink)
	}
}

func TestWithGitHubRange(t *testing.T) {
	link := resolveLink("github", "foo.go", "5", "10")
	expectedLink := "https://github.com/tonchis/dvcs-link/blob/commit-sha/foo.go#L5-L10"

	if link != expectedLink {
		t.Errorf("Expected %v to equal %v", link, expectedLink)
	}
}

func TestUnsupportedHost(t *testing.T) {
	err := verifyHost("github")
	if err != nil {
		t.Error("Expected to support github.")
	}

	err = verifyHost("gitlab")
	if err != nil {
		t.Error("Expected to support gitlab.")
	}

	err = verifyHost("bitbucket")
	if err == nil {
		t.Error("Unexpected bitbucket support.")
	}
}
