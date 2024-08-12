package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Checker struct {
	tokenizer *tokenizer
	prevToken token
	currToken token

	filename string
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
	dep.url = tok.text
}

func (c *Checker) beginDecl(dep *deployement) {
	c.Next()
	tok := c.Expect(String)
	dep.begin = tok.text
}

func (c *Checker) thenDecl(dep *deployement) {
	c.Next()
	tok := c.Expect(String)
	dep.then = tok.text
}

func (c *Checker) endDecl(dep *deployement) {
	c.Next()
	tok := c.Expect(String)
	dep.end = tok.text
}

func (c *Checker) extractDecl(dep *deployement) {
	c.Next()
	tok := c.Expect(String)
	split := strings.Split(tok.text, " ")
	if len(split) != 2 {
		c.Fatalf(c.tokenizer.pos, "expected two arguments but got %d", len(split))
	}
	dep.extracts = append(dep.extracts, extractment{
		name: split[0],
		out:  split[1],
	})
}

func parse(filename string) deployement {
	dep := &deployement{}
	buf, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	c := &Checker{
		tokenizer: newTokenizer(string(buf)),
		filename:  filename,
	}
	c.Next()
decls:
	for {
		switch curr := c.Current(); curr {
		case Clone:
			c.deployDecl(dep)
		case Begin:
			c.beginDecl(dep)
		case Extract:
			c.extractDecl(dep)
		case Then:
			c.thenDecl(dep)
		case End:
			c.endDecl(dep)
		default:
			break decls
		}
	}
	c.Allow(Semicolon)
	c.Expect(EOF)

	return *dep
}
