package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type deployement struct {
	wd       string
	url      string
	begin    string
	then     string
	end      string
	extracts []extractment
}

type extractment struct {
	name string
	out  string
}

func (e extractment) extract(dir string, wd string) {
	if strings.HasPrefix(e.out, ".") {
		e.out = wd + e.out[1:]
	}
	sh(fmt.Sprintf("mv %s/%s %s", dir, e.name, e.out))
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	rciFilePath := wd + "/Coolfile"
	dep := parse(rciFilePath)
	dep.wd = wd
	deploy(dep)
}

func deploy(d deployement) {
	d.url = "https://" + d.url
	dir := strconv.Itoa(rand.Intn(10000000))
	output := "/tmp/" + dir
	clone := fmt.Sprintf("git clone --depth=1 %s %s", d.url, output)

	split := strings.Split(d.url, "@")
	if len(split) == 2 {
		clone = fmt.Sprintf("git clone --depth=1 %s %s --branch %s", split[0], output, split[1])
	}

	cloning := []string{
		"cd /tmp",
		clone,
	}
	begining := []string{
		"cd " + output,
		d.begin,
	}

	sh(strings.Join(cloning, " && "))
	sh(strings.Join(begining, " && "))
	for _, e := range d.extracts {
		e.extract(output, d.wd)
	}
	sh(d.then)
	sh(d.end)
}

func sh(args string) {
	fmt.Println(args)
	cmd := exec.Command("sh", []string{"-c", args}...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
