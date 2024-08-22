package watch

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/restartfu/watch/internal/parser"
)

func Execute() {
	wd, _ := os.Getwd()
	rciFilePath := wd + "/Watchfile"

	cmds := parser.Parse(rciFilePath)
	deploy(cmds)
}

func deploy(cmds []string) {
	fmt.Println()
	for _, c := range cmds {
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
