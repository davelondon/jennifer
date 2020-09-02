package jen

import (
	"bytes"
	"fmt"
	"testing"
)

func TestGuessAlias(t *testing.T) {

	data := map[string]string{
		"A":             "a",
		"a":             "a",
		"a$":            "a",
		"a/b":           "b",
		"a/b/c":         "c",
		"a/b/c-d":       "cd",
		"a/b/c-d/":      "cd",
		"a.b":           "ab",
		"a/b.c":         "bc",
		"a/b-c.d":       "bcd",
		"a/bb-ccc.dddd": "bbcccdddd",
		"a/foo-go":      "foogo",
		"123a":          "a",
		"a/321a.b":      "ab",
		"a/123":         "pkg",
	}
	for path, expected := range data {
		if guessAlias(path) != expected {
			fmt.Printf("guessAlias test failed %s should return %s but got %s\n", path, expected, guessAlias(path))
			t.Fail()
		}
	}
}

func TestValidAlias(t *testing.T) {
	data := map[string]bool{
		"a":   true,  // ok
		"b":   false, // already registered
		"go":  false, // keyword
		"int": false, // predeclared
		"err": false, // common name
	}
	f := NewFile("test")
	f.register("b")
	for alias, expected := range data {
		if f.isValidAlias(alias) != expected {
			fmt.Printf("isValidAlias test failed %s should return %t but got %t\n", alias, expected, f.isValidAlias(alias))
			t.Fail()
		}
	}
}

func TestFile_ImportComment(t *testing.T) {
	file := NewFile("main")

	file.Anon("fmt")
	file.ImportComment("fmt", "anonymous fmt")

	file.ImportName("io", "ioname")
	file.ImportComment("io", "io comment")

	file.ImportAlias("ioutil", "ioutilalias")
	file.ImportComment("ioutil", "ioutil comment")

	file.Func().Id("main").Params().Block(
		Qual("io", "Pipe").Call(),
		Qual("ioutil", "NopCloser").Call(Nil()),
	)

	var got bytes.Buffer
	err := file.Render(&got)
	if err != nil {
		t.Fatal(err)
	}

	expect := `package main

import (
	_ "fmt"              // anonymous fmt
	"io"                 // io comment
	ioutilalias "ioutil" // ioutil comment
)

func main() {
	ioname.Pipe()
	ioutilalias.NopCloser(nil)
}
`

	if got.String() != expect {
		t.Fatalf("Got: %v, expect: %v", got.String(), expect)
	}
}
