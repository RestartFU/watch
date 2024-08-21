package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type repository struct {
	path string
	url  string
}

type Checker struct {
	tokenizer *tokenizer
	prevToken token
	currToken token

	variables map[string]string
	filename  string
}

func (c *Checker) Fatalf(pos pos, format string, args ...any) {
	fmt.Printf("%s(%d:%d)", c.filename, pos.line, pos.column)
	fmt.Printf(format, args...)
	fmt.Println()
	os.Exit(1)
}

func (c *Checker) Next() (res token) {
	token, err := c.tokenizer.Token()
	if err != nil && err != io.EOF {
		c.Fatalf(c.tokenizer.pos, " found invalid token: %v", err)
	}
	c.prevToken, c.currToken = c.currToken, token
	return c.prevToken
}

func (c *Checker) Expect(kind tokenKind) token {
	token := c.Next()
	if token.kind != kind {
		c.Fatalf(token.pos, " expected token %v, got %v", kind, token.kind)
	}
	return token
}

func (c *Checker) Allow(kind tokenKind) bool {
	if c.currToken.kind == kind {
		c.Next()
		return true
	}
	return false
}

func (c *Checker) Current() tokenKind {
	if c.currToken.kind == Comment {
		c.Next()
		return c.Current()
	}
	return c.currToken.kind
}

func (c *Checker) deployDecl(dep *deployement) {
	c.Next()
	tok := c.Expect(String)

	dep.repository.url = "https://" + tok.text
	dir := strconv.Itoa(rand.Intn(10000000))
	dep.repository.path = "/tmp/" + dir

	fmt.Println(dep.repository)
	if c.Allow(As) {
		str := c.Expect(String)
		c.variables[str.text] = dep.repository.path
		fmt.Println(c.variables)
	}
}

func (c *Checker) runDecl(dep *deployement) {
	c.Next()
	tok := c.Expect(String)
	str := tok.text
	for k, v := range c.variables {
		str = strings.ReplaceAll(str, fmt.Sprintf("$[%s]", k), v)
	}

	dep.commands = append(dep.commands, str)
	fmt.Println(dep.commands)
}

func (c *Checker) extractDecl(dep *deployement) {
	c.Next()
	tok := c.Expect(String)
	str := tok.text
	for k, v := range c.variables {
		str = strings.ReplaceAll(str, fmt.Sprintf("$[%s]", k), v)
	}

	split := strings.Split(str, " ")
	if len(split) != 2 {
		c.Fatalf(c.tokenizer.pos, " expected two arguments but got %d", len(split))
	}

	out := split[1]
	if strings.HasPrefix(out, ".") {
		out = wd + out[1:]
	}

	cmd := fmt.Sprintf("mv %s %s", split[0], out)
	dep.commands = append(dep.commands, cmd)
}

func parse(filename string) deployement {
	dep := &deployement{}
	buf, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	c := &Checker{
		tokenizer: newTokenizer(string(buf)),
		variables: map[string]string{},
		filename:  filename,
	}
	c.Next()
decls:
	for {
		switch curr := c.Current(); curr {
		case Clone:
			c.deployDecl(dep)
		case Run:
			c.runDecl(dep)
		case Extract:
			c.extractDecl(dep)
		default:
			fmt.Println()
			break decls
		}
	}
	c.Allow(Semicolon)
	c.Expect(EOF)

	return *dep
}
