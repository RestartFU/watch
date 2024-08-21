package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

var wd, _ = os.Getwd()

type deployement struct {
	repository repository
	commands   []string
}

func main() {
	rciFilePath := wd + "/Watchfile"
	dep := parse(rciFilePath)

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	deploy(dep)
}

func deploy(d deployement) {
	clone := fmt.Sprintf("git clone --depth=1 %s %s", d.repository, d.repository)

	split := strings.Split(d.repository.url, "@")
	if len(split) == 2 {
		clone = fmt.Sprintf("git clone --depth=1 %s %s --branch %s", split[0], d.repository.path, split[1])
	}

	cloning := []string{
		"cd /tmp",
		clone,
	}
	sh(strings.Join(cloning, " && "))

	for _, c := range d.commands {
		sh(c)
	}
}

func sh(args string) {
	fmt.Println(args)
	cmd := exec.Command("sh", []string{"-c", args}...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
