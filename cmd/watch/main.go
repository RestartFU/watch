package watch

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/restartfu/watch/internal/parser"
)

func Execute() {
	wd, _ := os.Getwd()
	rciFilePath := wd + "/Watchfile"
	dep := parser.Parse(rciFilePath)

	deploy(dep)
}

func deploy(d parser.Result) {
	clone := fmt.Sprintf("git clone --depth=1 %s %s", d.RepositoryURL, d.RepositoryTemporaryPath)

	split := strings.Split(d.RepositoryURL, "@")
	if len(split) == 2 {
		clone = fmt.Sprintf("git clone --depth=1 %s %s --branch %s", split[0], d.RepositoryTemporaryPath, split[1])
	}

	cloning := []string{
		"cd /tmp",
		clone,
	}
	sh(strings.Join(cloning, " && "))

	for _, c := range d.Commands {
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
