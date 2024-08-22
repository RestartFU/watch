package parser

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/restartfu/watch/internal/tokenizer"
)

type repository struct {
	path string
	url  string
}

type Checker struct {
	tokenizer *tokenizer.Tokenizer
	prevToken tokenizer.Token
	currToken tokenizer.Token

	variables map[string]string
	filename  string
}

func (c *Checker) Fatalf(pos tokenizer.Position, format string, args ...any) {
	fmt.Printf("%s(%d:%d)", c.filename, pos.Line(), pos.Column())
	fmt.Printf(format, args...)
	fmt.Println()
	os.Exit(1)
}

func (c *Checker) Next() (res tokenizer.Token) {
	token, err := c.tokenizer.Token()
	if err != nil && err != io.EOF {
		c.Fatalf(c.tokenizer.Position, " found invalid token: %v", err)
	}
	c.prevToken, c.currToken = c.currToken, token
	return c.prevToken
}

func (c *Checker) Expect(kind tokenizer.TokenKind) tokenizer.Token {
	token := c.Next()
	if token.Kind() != kind {
		c.Fatalf(token.Position, " expected token %v, got %v", kind, token.Kind())
	}
	return token
}

func (c *Checker) Allow(kind tokenizer.TokenKind) bool {
	if c.currToken.Kind() == kind {
		c.Next()
		return true
	}
	return false
}

func (c *Checker) Current() tokenizer.TokenKind {
	if c.currToken.Kind() == tokenizer.Comment {
		c.Next()
		return c.Current()
	}
	return c.currToken.Kind()
}

func (c *Checker) cloneDecl(dep *Result) {
	c.Next()
	tok := c.Expect(tokenizer.String)

	dep.RepositoryURL = "https://" + tok.Text()
	dir := strconv.Itoa(rand.Intn(10000000))
	dep.RepositoryTemporaryPath = "/tmp/" + dir

	if c.Allow(tokenizer.As) {
		str := c.Expect(tokenizer.String)
		c.variables[str.Text()] = dep.RepositoryTemporaryPath
		fmt.Println(c.variables)
	}
}

func (c *Checker) runDecl(dep *Result) {
	c.Next()
	tok := c.Expect(tokenizer.String)
	str := tok.Text()
	for k, v := range c.variables {
		str = strings.ReplaceAll(str, fmt.Sprintf("$[%s]", k), v)
	}

	dep.Commands = append(dep.Commands, str)
	fmt.Println(dep.Commands)
}

func (c *Checker) extractDecl(dep *Result) {
	c.Next()
	tok := c.Expect(tokenizer.String)
	str := tok.Text()
	for k, v := range c.variables {
		str = strings.ReplaceAll(str, fmt.Sprintf("$[%s]", k), v)
	}

	split := strings.Split(str, " ")
	if len(split) != 2 {
		c.Fatalf(c.tokenizer.Position, " expected two arguments but got %d", len(split))
	}

	wd, _ := os.Getwd()
	out := split[1]
	if strings.HasPrefix(out, ".") {
		out = wd + out[1:]
	}

	cmd := fmt.Sprintf("mv %s %s", split[0], out)
	dep.Commands = append(dep.Commands, cmd)
}

func (c *Checker) setDecl() {
	c.Next()
	tok := c.Expect(tokenizer.String)
	str := tok.Text()
	for k, v := range c.variables {
		str = strings.ReplaceAll(str, fmt.Sprintf("$[%s]", k), v)
	}

	split := strings.Split(str, "=")
	if len(split) != 2 {
		c.Fatalf(c.tokenizer.Position, " expected two arguments separated by '=' but got %d", len(split))
	}

	c.variables[split[0]] = split[1]
}

func Parse(filename string) Result {
	dep := &Result{}
	buf, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	c := &Checker{
		tokenizer: tokenizer.NewTokenizer(string(buf)),
		variables: map[string]string{},
		filename:  filename,
	}
	c.Next()
decls:
	for {
		switch curr := c.Current(); curr {
		case tokenizer.Clone:
			c.cloneDecl(dep)
		case tokenizer.Run:
			c.runDecl(dep)
		case tokenizer.Extract:
			c.extractDecl(dep)
		case tokenizer.Set:
			c.setDecl()
		default:
			break decls
		}
	}
	c.Allow(tokenizer.Semicolon)
	c.Expect(tokenizer.EOF)

	return *dep
}

type Result struct {
	RepositoryURL,
	RepositoryTemporaryPath string

	Commands []string
}
