package main

import (
	"fmt"
	"os"
	"os/exec"
	"log"
	"strings"
)

func main() {
	var file, start, end string

	args := os.Args[1:]
	switch len(args) {
	case 0:
		fmt.Println("No file given.")
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

	remoteUrl := bashExec("git remote -v | grep 'origin' | grep 'push' | awk '{ print $2 }'")

	if strings.HasPrefix(remoteUrl, "git@") {
		remoteUrl = "https://" + strings.Replace(strings.Split(remoteUrl, "@")[1], ":", "/", 1)
	}

	commit := bashExec("git rev-parse HEAD")

	githubUrl := strings.Join([]string {remoteUrl, "blob", commit, file}, "/")

	if start != "" {
		githubUrl = githubUrl + fmt.Sprintf("#L%v", start)
	}
	
	if end != "" {
		githubUrl = githubUrl + fmt.Sprintf("-L%v", end)
	}

	fmt.Println(githubUrl)
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
