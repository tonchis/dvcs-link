package main

import (
	"fmt"
	"testing"
)

func TestFilenameOnly(t *testing.T) {
	link := resolveLink("", "foo.go", "", "")
	commit := bashExec("git rev-parse HEAD")
	expectedLink := fmt.Sprintf("https://gitlab.com/tonchis/dvcs-link/blob/%v/foo.go", commit)

	if link != expectedLink {
		t.Errorf("Expected %v to equal %v", link, expectedLink)
	}
}

func TestWithLine(t *testing.T) {
	link := resolveLink("", "foo.go", "5", "")
	commit := bashExec("git rev-parse HEAD")
	expectedLink := fmt.Sprintf("https://gitlab.com/tonchis/dvcs-link/blob/%v/foo.go#L5", commit)

	if link != expectedLink {
		t.Errorf("Expected %v to equal %v", link, expectedLink)
	}
}

func TestWithGitLabRange(t *testing.T) {
	link := resolveLink("", "foo.go", "5", "10")
	commit := bashExec("git rev-parse HEAD")
	expectedLink := fmt.Sprintf("https://gitlab.com/tonchis/dvcs-link/blob/%v/foo.go#L5-10", commit)

	if link != expectedLink {
		t.Errorf("Expected %v to equal %v", link, expectedLink)
	}
}

func TestWithGitHubRange(t *testing.T) {
	link := resolveLink("github", "foo.go", "5", "10")
	commit := bashExec("git rev-parse HEAD")
	expectedLink := fmt.Sprintf("https://github.com/tonchis/dvcs-link/blob/%v/foo.go#L5-L10", commit)

	if link != expectedLink {
		t.Errorf("Expected %v to equal %v", link, expectedLink)
	}
}

func TestWithHost(t *testing.T) {
	link := resolveLink("github", "foo.go", "5", "10")
	commit := bashExec("git rev-parse HEAD")
	expectedLink := fmt.Sprintf("https://github.com/tonchis/dvcs-link/blob/%v/foo.go#L5-L10", commit)

	if link != expectedLink {
		t.Errorf("Expected %v to equal %v", link, expectedLink)
	}
}
