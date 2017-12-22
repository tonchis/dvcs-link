package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Shell interface {
	Run(command string) string
}

type Bash struct {
}

func (b Bash) Run(command string) string {
	remoteCommand := exec.Command("bash", "-c", command)

	if path, err := exec.LookPath("bash"); err != nil {
		fmt.Println("[bashExec] binary not in path")
		log.Fatal(err)
	} else {
		remoteCommand.Path = string(path)
	}

	var remote string
	if output, err := remoteCommand.Output(); err != nil {
		fmt.Println("[bashExec] command execution failed")
		log.Fatal(err)
	} else {
		remote = strings.TrimSpace(string(output))
	}

	return remote
}

var bashExec = Bash{}.Run

func main() {
	var remote, file, start, end string

	flag.StringVar(&remote, "remote", "origin", "Name of the remote repository. Default 'origin'.")
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `%[1]s is a tool to print the GitHub or GitLab link of a file and range.

Usage:

	$ %[1]s foo.go
	https://github.com/user/repo/blob/commit/foo.go
	
	$ %[1]s foo.go 5
	https://github.com/user/repo/blob/commit/foo.go#L5
	
	$ %[1]s foo.go 5 10
	https://github.com/user/repo/blob/commit/foo.go#L5-L10
	
	$ %[1]s -remote mirror foo.go 5 10
	https://gitlab.com/user/repo/blob/commit/foo.go#L5-L10
`, os.Args[0])
	}

	args := flag.Args()
	switch len(args) {
	case 1:
		file = args[0]
	case 2:
		file = args[0]
		start = args[1]
	case 3:
		file = args[0]
		start = args[1]
		end = args[2]
	default:
		flag.Usage()
		os.Exit(1)
	}

	err := verifyHost(remote)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(resolveLink(remote, file, start, end))
}

func verifyHost(remote string) error {
	host := bashExec(fmt.Sprintf("git config --local --get remote.%v.url", remote))

	if !(strings.Contains(host, "github") || strings.Contains(host, "gitlab")) {
		return errors.New("Unsupported host. Only `github` or `gitlab` are supported.")
	}

	return nil
}

func resolveLink(remote string, file string, start string, end string) string {
	hostUrl := resolveHost(remote)

	commit := bashExec("git rev-parse HEAD")

	link := strings.Join([]string{hostUrl, "blob", commit, file}, "/")

	if start != "" {
		link = link + fmt.Sprintf("#L%v", start)

		link = resolveRange(link, end)
	}

	return link
}

func resolveHost(remote string) string {
	remoteUrl := bashExec(fmt.Sprintf("git config --local --get remote.%v.url", remote))
	remoteUrl = normalizeUrl(remoteUrl)

	return "https://" + remoteUrl
}

func normalizeUrl(remoteUrl string) string {
	if strings.HasPrefix(remoteUrl, "https") {
		return strings.Replace(remoteUrl, "https://", "", 1)
	}

	res := strings.Split(remoteUrl, "@")[1]
	res = strings.Replace(res, ":", "/", 1)
	res = strings.Replace(res, ".git", "", 1)

	return res
}

func resolveRange(link string, end string) string {
	if end == "" {
		return link
	}

	var suffix string
	if strings.Contains(link, "github") {
		suffix = "L" + end
	} else if strings.Contains(link, "gitlab") {
		suffix = end
	}

	return link + fmt.Sprintf("-%v", suffix)
}
