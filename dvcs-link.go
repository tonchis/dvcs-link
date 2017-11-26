package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func usage() string {
	return "No file given."
}

func main() {
	var file, start, end string

	args := os.Args[1:]
	switch len(args) {
	case 0:
		fmt.Println(usage())

		return
	case 1:
		file = args[0]
	case 2:
		file = args[0]
		start = args[1]
	default:
		file = args[0]
		start = args[1]
		end = args[2]
	}

	remoteOriginUrl := bashExec("git remote -v | grep 'origin.*push' | awk '{ print $2 }'")

	remoteOriginUrl = convertToHttps(remoteOriginUrl)

	commit := bashExec("git rev-parse HEAD")

	dvcsLink := strings.Join([]string{remoteOriginUrl, "blob", commit, file}, "/")

	if start != "" {
		dvcsLink = dvcsLink + fmt.Sprintf("#L%v", start)

		if end != "" {
			dvcsLink = dvcsLink + fmt.Sprintf("-L%v", end)
		}
	}

	fmt.Println(dvcsLink)
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

func convertToHttps(remoteOriginUrl string) string {
	if strings.HasPrefix(remoteOriginUrl, "https") {
		return remoteOriginUrl
	}

	res := strings.Split(remoteOriginUrl, "@")[1]
	res = strings.Replace(res, ":", "/", 1)
	res = strings.Replace(res, ".git", "", 1)

	return "https://" + res
}
