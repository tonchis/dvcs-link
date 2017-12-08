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

func main() {
	var host, file, start, end string

	flag.StringVar(&host, "host", "", "DVCS host: `github` or `gitlab`")
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
	
	$ %[1]s -host gitlab foo.go 5 10
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

	_, err := verifyHost(host)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(resolveLink(host, file, start, end))
}

func bashExec(command string) string {
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

func normalizeUrl(remoteOriginUrl string) string {
	if strings.HasPrefix(remoteOriginUrl, "https") {
		return strings.Replace(remoteOriginUrl, "https://", "", 1)
	}

	res := strings.Split(remoteOriginUrl, "@")[1]
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

func resolveLink(host string, file string, start string, end string) string {
	remoteOriginUrl := bashExec("git config --local --get remote.origin.url")
	remoteOriginUrl = normalizeUrl(remoteOriginUrl)

	if host != "" {
		hostAndPath := strings.SplitN(remoteOriginUrl, "/", 2)
		remoteOriginUrl = fmt.Sprintf("%v.com/%v", host, hostAndPath[1])
	}

	remoteOriginUrl = "https://" + remoteOriginUrl

	commit := bashExec("git rev-parse HEAD")

	link := strings.Join([]string{remoteOriginUrl, "blob", commit, file}, "/")

	if start != "" {
		link = link + fmt.Sprintf("#L%v", start)

		link = resolveRange(link, end)
	}

	return link
}

func verifyHost(host string) (string, error) {
	if host == "" {
		host = bashExec("git config --local --get remote.origin.url")
	}

	if !(strings.Contains(host, "github") || strings.Contains(host, "gitlab")) {
		return "", errors.New("Unsupported host. Only `github` or `gitlab` are supported.")
	}

	return "", nil
}
