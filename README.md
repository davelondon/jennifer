[![Go Reference](https://pkg.go.dev/badge/github.com/dave/jennifer/jen.svg)](https://pkg.go.dev/github.com/dave/jennifer/jen)
[![Build Status](https://travis-ci.org/dave/jennifer.svg?branch=master)](https://travis-ci.org/dave/jennifer) [![Go Report Card](https://goreportcard.com/badge/github.com/dave/jennifer)](https://goreportcard.com/report/github.com/dave/jennifer)
[![codecov](https://img.shields.io/badge/codecov-100%25-brightgreen.svg)](https://codecov.io/gh/dave/jennifer)
![stability-stable](https://img.shields.io/badge/stability-stable-brightgreen.svg)
[![Sourcegraph](https://sourcegraph.com/github.com/dave/jennifer/-/badge.svg)](https://sourcegraph.com/github.com/dave/jennifer?badge)

# Jennifer
Jennifer is a code generator for Go.

```go
package main

import (
    "fmt"

    . "github.com/dave/jennifer/jen"
)

func main() {
	f := NewFile("main")
	f.Func().Id("main").Params().Block(
		Qual("fmt", "Println").Call(Lit("Hello, world")),
	)
	fmt.Printf("%#v", f)
}
```
Output:
```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, world")
}
```

### Install
```
go get -u github.com/dave/jennifer/jen
```

### Need help? 
If you get stuck, have a question, would like a code review, or just want a 
chat: I'm happy to help! Feel free to open an issue, email me or mention @dave 
in your PR.

### Examples
Jennifer has a comprehensive suite of examples - see [godoc](https://godoc.org/github.com/dave/jennifer/jen#pkg-examples) for an index. Here's some examples of jennifer being used in the real-world:

* [genjen](genjen/render.go) (which generates much of jennifer, using data in [data.go](genjen/data.go))
* [zerogen](https://github.com/mrsinham/zerogen/blob/master/generator.go)
* [go-contentful-generator](https://github.com/nicolai86/go-contentful-generator)

### Rendering
For testing, a File or Statement can be rendered with the fmt package 
using the %#v verb.

```go
c := Id("a").Call(Lit("b"))
fmt.Printf("%#v", c)
// Output:
// a("b")
```

This is not recommended for use in production because any error will cause a 
panic. For production use, [File.Render](#render) or [File.Save](#save) are 
preferred.

# Identifiers
**Identifiers** [Keywords](#keywords) [Operators](#operators) [Braces](#braces) [Parentheses](#parentheses) [Control flow](#control-flow) [Collections](#collections) [Literals](#literals) [Comments](#comments) [Helpers](#helpers) [Misc](#misc) [File](#file)

### Id
Id renders an identifier.

```go
c := If(Id("i").Op("==").Id("j")).Block(
	Return(Id("i")),
)
fmt.Printf("%#v", c)
// Output:
// if i == j {
// 	return i
// }
```

### Dot
Dot renders a period followed by an identifier. Use for fields and selectors. 

```go
c := Qual("a.b/c", "Foo").Call().Dot("Bar").Index(Lit(0)).Dot("Baz")
fmt.Printf("%#v", c)
// Output:
// c.Foo().Bar[0].Baz
```

### Qual
Qual renders a qualified identifier.

```go
c := Qual("encoding/gob", "NewEncoder").Call()
fmt.Printf("%#v", c)
// Output:
// gob.NewEncoder()
```

Imports are automatically added when
used with a File. If the path matches the local path, the package name is
omitted. If package names conflict they are automatically renamed.

```go
f := NewFilePath("a.b/c")
f.Func().Id("init").Params().Block(
	Qual("a.b/c", "Foo").Call().Comment("Local package - name is omitted."),
	Qual("d.e/f", "Bar").Call().Comment("Import is automatically added."),
	Qual("g.h/f", "Baz").Call().Comment("Colliding package name is renamed."),
)
fmt.Printf("%#v", f)
// Output:
// package c
//
// import (
// 	f "d.e/f"
// 	f1 "g.h/f"
// )
//
// func init() {
// 	Foo()    // Local package - name is omitted.
// 	f.Bar()  // Import is automatically added.
// 	f1.Baz() // Colliding package name is renamed.
// }
```

Note that
it is not possible to reliably determine the package name given an arbitrary
package path, so a sensible name is guessed from the path and added as an
alias. The names of all standard library packages are known so these do not
need to be aliased. If more control is needed of the aliases, see
[File.ImportName](#importname) or [File.ImportAlias](#importalias).

### List
List renders a comma separated list. Use for multiple return functions.

```go
c := List(Id("a"), Err()).Op(":=").Id("b").Call()
fmt.Printf("%#v", c)
// Output:
// a, err := b()
```

# Keywords
[Identifiers](#identifiers) **Keywords** [Operators](#operators) [Braces](#braces) [Parentheses](#parentheses) [Control flow](#control-flow) [Collections](#collections) [Literals](#literals) [Comments](#comments) [Helpers](#helpers) [Misc](#misc) [File](#file)

Simple keywords, predeclared identifiers and built-in functions are self 
explanatory:

| Construct        | Name |
| ---------------- | ---- |
| Keywords         | Break, Chan, Const, Continue, Default, Defer, Else, Fallthrough, Func, Go, Goto, Range, Select, Type, Var |
| Functions        | Append, Cap, Close, Complex, Copy, Delete, Imag, Len, Make, New, Panic, Print, Println, Real, Recover |
| Types            | Bool, Byte, Complex64, Complex128, Error, Float32, Float64, Int, Int8, Int16, Int32, Int64, Rune, String, Uint, Uint8, Uint16, Uint32, Uint64, Uintptr |
| Constants        | True, False, Iota, Nil |
| Helpers          | Err |

Built-in functions take a list of parameters and render them appropriately:

```go
c := Id("a").Op("=").Append(Id("a"), Id("b").Op("..."))
fmt.Printf("%#v", c)
// Output:
// a = append(a, b...)
```

Special cases for [If, For](#if-for), [Interface, Struct](#interface-struct), [Switch, Case](#switch-select), [Return](#return) and [Map](#map) are explained below.

# Operators
[Identifiers](#identifiers) [Keywords](#keywords) **Operators** [Braces](#braces) [Parentheses](#parentheses) [Control flow](#control-flow) [Collections](#collections) [Literals](#literals) [Comments](#comments) [Helpers](#helpers) [Misc](#misc) [File](#file)

Op renders the provided operator / token.

```go
c := Id("a").Op(":=").Id("b").Call()
fmt.Printf("%#v", c)
// Output:
// a := b()
```

```go
c := Id("a").Op("=").Op("*").Id("b")
fmt.Printf("%#v", c)
// Output:
// a = *b
```

```go
c := Id("a").Call(Id("b").Op("..."))
fmt.Printf("%#v", c)
// Output:
// a(b...)
```

```go
c := If(Parens(Id("a").Op("||").Id("b")).Op("&&").Id("c")).Block()
fmt.Printf("%#v", c)
// Output:
// if (a || b) && c {
// }
```

# Braces
[Identifiers](#identifiers) [Keywords](#keywords) [Operators](#operators) **Braces** [Parentheses](#parentheses) [Control flow](#control-flow) [Collections](#collections) [Literals](#literals) [Comments](#comments) [Helpers](#helpers) [Misc](#misc) [File](#file)

Several methods render curly braces, summarized below: 

| Name                           | Prefix       | Separator | Example                              |
| ------------------------------ | ------------ | --------- | -------------------------------------|
| [Block](#block)                |              | `\n`      | `func a() { ... }` or `if a { ... }` |
| [Interface](#interface-struct) | `interface`  | `\n`      | `interface { ... }`                  |
| [Struct](#interface-struct)    | `struct`     | `\n`      | `struct { ... }`                     |
| [Values](#values)              |              | `,`       | `[]int{1, 2}` or `A{B: "c"}`         |

### Block
Block renders a statement list enclosed by curly braces. Use for code blocks.

```go
c := Func().Id("foo").Params().String().Block(
	Id("a").Op("=").Id("b"),
	Id("b").Op("++"),
	Return(Id("b")),
)
fmt.Printf("%#v", c)
// Output:
// func foo() string {
// 	a = b
// 	b++
// 	return b
// }
```

```go
c := If(Id("a").Op(">").Lit(10)).Block(
	Id("a").Op("=").Id("a").Op("/").Lit(2),
)
fmt.Printf("%#v", c)
// Output:
// if a > 10 {
// 	a = a / 2
// }
```

A special case applies when used directly after Case or Default, where the braces are omitted. This allows use in switch and select statements. [See example](#switch-select).

### Interface, Struct
Interface and Struct render the keyword followed by a statement list enclosed 
by curly braces.

```go
c := Var().Id("a").Interface()
fmt.Printf("%#v", c)
// Output:
// var a interface{}
```

```go
c := Type().Id("a").Interface(
	Id("b").Params().String(),
)
fmt.Printf("%#v", c)
// Output:
// type a interface {
// 	b() string
// }
```

```go
c := Id("c").Op(":=").Make(Chan().Struct())
fmt.Printf("%#v", c)
// Output:
// c := make(chan struct{})
```

```go
c := Type().Id("foo").Struct(
	List(Id("x"), Id("y")).Int(),
	Id("u").Float32(),
)
fmt.Printf("%#v", c)
// Output:
// type foo struct {
// 	x, y int
// 	u    float32
// }
```

# Parentheses
[Identifiers](#identifiers) [Keywords](#keywords) [Operators](#operators) [Braces](#braces) **Parentheses** [Control flow](#control-flow) [Collections](#collections) [Literals](#literals) [Comments](#comments) [Helpers](#helpers) [Misc](#misc) [File](#file)

Several methods output parenthesis, summarized below:

| Name              | Prefix | Separator | Example                           |
| ----------------- | ------ | --------- | --------------------------------- |
| [Call](#call)     |        | `,`       | `fmt.Println(b, c)`               |
| [Params](#params) |        | `,`       | `func (a *A) Foo(i int) { ... }`  |
| [Defs](#defs)     |        | `\n`      | `const ( ... )`                   |
| [Parens](#parens) |        |           | `[]byte(s)` or `a / (b + c)`      |
| [Assert](#assert) | `.`    |           | `s, ok := i.(string)`             |

### Call
Call renders a comma separated list enclosed by parenthesis. Use for function calls.

```go
c := Qual("fmt", "Printf").Call(
	Lit("%#v: %T\n"),
	Id("a"),
	Id("b"),
)
fmt.Printf("%#v", c)
// Output:
// fmt.Printf("%#v: %T\n", a, b)
```

### Params
Params renders a comma separated list enclosed by parenthesis. Use for function parameters and method receivers.

```go
c := Func().Params(
	Id("a").Id("A"),
).Id("foo").Params(
	Id("b"),
	Id("c").String(),
).String().Block(
	Return(Id("b").Op("+").Id("c")),
)
fmt.Printf("%#v", c)
// Output:
// func (a A) foo(b, c string) string {
// 	return b + c
// }
```

### Defs
Defs renders a statement list enclosed in parenthesis. Use for definition lists.

```go
c := Const().Defs(
	Id("a").Op("=").Lit("a"),
	Id("b").Op("=").Lit("b"),
)
fmt.Printf("%#v", c)
// Output:
// const (
// 	a = "a"
// 	b = "b"
// )
```

### Parens
Parens renders a single item in parenthesis. Use for type conversion or to specify evaluation order.

```go
c := Id("b").Op(":=").Index().Byte().Parens(Id("s"))
fmt.Printf("%#v", c)
// Output:
// b := []byte(s)
```

```go
c := Id("a").Op("/").Parens(Id("b").Op("+").Id("c"))
fmt.Printf("%#v", c)
// Output:
// a / (b + c)
```

### Assert
Assert renders a period followed by a single item enclosed by parenthesis. Use for type assertions.

```go
c := List(Id("b"), Id("ok")).Op(":=").Id("a").Assert(Bool())
fmt.Printf("%#v", c)
// Output:
// b, ok := a.(bool)
```

# Control flow
[Identifiers](#identifiers) [Keywords](#keywords) [Operators](#operators) [Braces](#braces) [Parentheses](#parentheses) **Control flow** [Collections](#collections) [Literals](#literals) [Comments](#comments) [Helpers](#helpers) [Misc](#misc) [File](#file)

### If, For
If and For render the keyword followed by a semicolon separated list.

```go
c := If(
	Err().Op(":=").Id("a").Call(),
	Err().Op("!=").Nil(),
).Block(
	Return(Err()),
)
fmt.Printf("%#v", c)
// Output:
// if err := a(); err != nil {
// 	return err
// }
```

```go
c := For(
	Id("i").Op(":=").Lit(0),
	Id("i").Op("<").Lit(10),
	Id("i").Op("++"),
).Block(
	Qual("fmt", "Println").Call(Id("i")),
)
fmt.Printf("%#v", c)
// Output:
// for i := 0; i < 10; i++ {
// 	fmt.Println(i)
// }
```

### Switch, Select
Switch, Select, Case and Block are used to build switch or select statements:

```go
c := Switch(Id("value").Dot("Kind").Call()).Block(
	Case(Qual("reflect", "Float32"), Qual("reflect", "Float64")).Block(
		Return(Lit("float")),
	),
	Case(Qual("reflect", "Bool")).Block(
		Return(Lit("bool")),
	),
	Case(Qual("reflect", "Uintptr")).Block(
		Fallthrough(),
	),
	Default().Block(
		Return(Lit("none")),
	),
)
fmt.Printf("%#v", c)
// Output:
// switch value.Kind() {
// case reflect.Float32, reflect.Float64:
// 	return "float"
// case reflect.Bool:
// 	return "bool"
// case reflect.Uintptr:
// 	fallthrough
// default:
// 	return "none"
// }
```

### Return
Return renders the keyword followed by a comma separated list.

```go
c := Return(Id("a"), Id("b"))
fmt.Printf("%#v", c)
// Output:
// return a, b
```

# Collections
[Identifiers](#identifiers) [Keywords](#keywords) [Operators](#operators) [Braces](#braces) [Parentheses](#parentheses) [Control flow](#control-flow) **Collections** [Literals](#literals) [Comments](#comments) [Helpers](#helpers) [Misc](#misc) [File](#file)

### Map
Map renders the keyword followed by a single item enclosed by square brackets. Use for map definitions.

```go
c := Id("a").Op(":=").Map(String()).String().Values()
fmt.Printf("%#v", c)
// Output:
// a := map[string]string{}
```

### Index
Index renders a colon separated list enclosed by square brackets. Use for array / slice indexes and definitions.

```go
c := Var().Id("a").Index().String()
fmt.Printf("%#v", c)
// Output:
// var a []string
```

```go
c := Id("a").Op(":=").Id("b").Index(Lit(0), Lit(1))
fmt.Printf("%#v", c)
// Output:
// a := b[0:1]
```

```go
c := Id("a").Op(":=").Id("b").Index(Lit(1), Empty())
fmt.Printf("%#v", c)
// Output:
// a := b[1:]
```

### Values
Values renders a comma separated list enclosed by curly braces. Use for slice or composite literals.

```go
c := Index().String().Values(Lit("a"), Lit("b"))
fmt.Printf("%#v", c)
// Output:
// []string{"a", "b"}
```

Dict renders as key/value pairs. Use with Values for map or composite
literals.

```go
c := Map(String()).String().Values(Dict{
	Lit("a"):	Lit("b"),
	Lit("c"):	Lit("d"),
})
fmt.Printf("%#v", c)
// Output:
// map[string]string{
// 	"a": "b",
// 	"c": "d",
// }
```

```go
c := Op("&").Id("Person").Values(Dict{
	Id("Age"):	Lit(1),
	Id("Name"):	Lit("a"),
})
fmt.Printf("%#v", c)
// Output:
// &Person{
// 	Age:  1,
// 	Name: "a",
// }
```

DictFunc executes a func(Dict) to generate the value.

```go
c := Id("a").Op(":=").Map(String()).String().Values(DictFunc(func(d Dict) {
	d[Lit("a")] = Lit("b")
	d[Lit("c")] = Lit("d")
}))
fmt.Printf("%#v", c)
// Output:
// a := map[string]string{
// 	"a": "b",
// 	"c": "d",
// }
```

Note: the items are ordered by key when rendered to ensure repeatable code.

# Literals
[Identifiers](#identifiers) [Keywords](#keywords) [Operators](#operators) [Braces](#braces) [Parentheses](#parentheses) [Control flow](#control-flow) [Collections](#collections) **Literals** [Comments](#comments) [Helpers](#helpers) [Misc](#misc) [File](#file)

### Lit
Lit renders a literal. Lit supports only built-in types (bool, string, int, complex128, float64,
float32, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
Passing any other type will panic.

```go
c := Id("a").Op(":=").Lit("a")
fmt.Printf("%#v", c)
// Output:
// a := "a"
```

```go
c := Id("a").Op(":=").Lit(1.5)
fmt.Printf("%#v", c)
// Output:
// a := 1.5
```

LitFunc generates the value to render by executing the provided
function.

```go
c := Id("a").Op(":=").LitFunc(func() interface{} { return 1 + 1 })
fmt.Printf("%#v", c)
// Output:
// a := 2
```

For the default constant types (bool, int, float64, string, complex128), Lit 
will render the untyped constant.

| Code          | Output     |
| ------------- | ---------- |
| `Lit(true)`   | `true`     |
| `Lit(1)`      | `1`        |
| `Lit(1.0)`    | `1.0`      |
| `Lit("foo")`  | `"foo"`    |
| `Lit(0 + 1i)` | `(0 + 1i)` |

For all other built-in types (float32, int8, int16, int32, int64, uint, uint8, 
uint16, uint32, uint64, uintptr, complex64), Lit will also render the type.

| Code                     | Output              |
| ------------------------ | ------------------- |
| `Lit(float32(1))`        | `float32(1)`        |
| `Lit(int16(1))`          | `int16(1)`          |
| `Lit(uint8(0x1))`        | `uint8(0x1)`        |
| `Lit(complex64(0 + 1i))` | `complex64(0 + 1i)` |

The built-in alias types byte and rune need a special case. LitRune and LitByte 
render rune and byte literals.

| Code                     | Output      |
| ------------------------ | ----------- |
| `LitRune('x')`           | `'x'`       |
| `LitByte(byte(0x1))`     | `byte(0x1)` |

# Comments
[Identifiers](#identifiers) [Keywords](#keywords) [Operators](#operators) [Braces](#braces) [Parentheses](#parentheses) [Control flow](#control-flow) [Collections](#collections) [Literals](#literals) **Comments** [Helpers](#helpers) [Misc](#misc) [File](#file)

### Comment
Comment adds a comment. If the provided string contains a newline, the
comment is formatted in multiline style.

```go
f := NewFile("a")
f.Comment("Foo returns the string \"foo\"")
f.Func().Id("Foo").Params().String().Block(
	Return(Lit("foo")).Comment("return the string foo"),
)
fmt.Printf("%#v", f)
// Output:
// package a
//
// // Foo returns the string "foo"
// func Foo() string {
// 	return "foo" // return the string foo
// }
```

```go
c := Comment("a\nb")
fmt.Printf("%#v", c)
// Output:
// /*
// a
// b
// */
```

If the comment string starts
with "//" or "/*", the automatic formatting is disabled and the string is
rendered directly.

```go
c := Id("foo").Call(Comment("/* inline */")).Comment("//no-space")
fmt.Printf("%#v", c)
// Output:
// foo( /* inline */ ) //no-space
```

### Commentf
Commentf adds a comment, using a format string and a list of parameters.

```go
name := "foo"
val := "bar"
c := Id(name).Op(":=").Lit(val).Commentf("%s is the string \"%s\"", name, val)
fmt.Printf("%#v", c)
// Output:
// foo := "bar" // foo is the string "bar"
```

# Helpers
[Identifiers](#identifiers) [Keywords](#keywords) [Operators](#operators) [Braces](#braces) [Parentheses](#parentheses) [Control flow](#control-flow) [Collections](#collections) [Literals](#literals) [Comments](#comments) **Helpers** [Misc](#misc) [File](#file)

### Func methods
All constructs that accept a variadic list of items are paired with GroupFunc 
functions that accept a func(*Group). Use for embedding logic.

```go
c := Id("numbers").Op(":=").Index().Int().ValuesFunc(func(g *Group) {
	for i := 0; i <= 5; i++ {
		g.Lit(i)
	}
})
fmt.Printf("%#v", c)
// Output:
// numbers := []int{0, 1, 2, 3, 4, 5}
```

```go
increment := true
name := "a"
c := Func().Id("a").Params().BlockFunc(func(g *Group) {
	g.Id(name).Op("=").Lit(1)
	if increment {
		g.Id(name).Op("++")
	} else {
		g.Id(name).Op("--")
	}
})
fmt.Printf("%#v", c)
// Output:
// func a() {
// 	a = 1
// 	a++
// }
```

### Add
Add appends the provided items to the statement.

```go
ptr := Op("*")
c := Id("a").Op("=").Add(ptr).Id("b")
fmt.Printf("%#v", c)
// Output:
// a = *b
```

```go
a := Id("a")
i := Int()
c := Var().Add(a, i)
fmt.Printf("%#v", c)
// Output:
// var a int
```

### Do
Do calls the provided function with the statement as a parameter. Use for
embedding logic.

```go
f := func(name string, isMap bool) *Statement {
	return Id(name).Op(":=").Do(func(s *Statement) {
		if isMap {
			s.Map(String()).String()
		} else {
			s.Index().String()
		}
	}).Values()
}
fmt.Printf("%#v\n%#v", f("a", true), f("b", false))
// Output:
// a := map[string]string{}
// b := []string{}
```

# Misc
[Identifiers](#identifiers) [Keywords](#keywords) [Operators](#operators) [Braces](#braces) [Parentheses](#parentheses) [Control flow](#control-flow) [Collections](#collections) [Literals](#literals) [Comments](#comments) [Helpers](#helpers) **Misc** [File](#file)

### Tag
Tag renders a struct tag

```go
c := Type().Id("foo").Struct(
	Id("A").String().Tag(map[string]string{"json": "a"}),
	Id("B").Int().Tag(map[string]string{"json": "b", "bar": "baz"}),
)
fmt.Printf("%#v", c)
// Output:
// type foo struct {
// 	A string `json:"a"`
// 	B int    `bar:"baz" json:"b"`
// }
```

Note: the items are ordered by key when rendered to ensure repeatable code.

### Null
Null adds a null item. Null items render nothing and are not followed by a
separator in lists.

In lists, nil will produce the same effect.

```go
c := Func().Id("foo").Params(
	nil,
	Id("s").String(),
	Null(),
	Id("i").Int(),
).Block()
fmt.Printf("%#v", c)
// Output:
// func foo(s string, i int) {}
```

### Empty
Empty adds an empty item. Empty items render nothing but are followed by a
separator in lists.

```go
c := Id("a").Op(":=").Id("b").Index(Lit(1), Empty())
fmt.Printf("%#v", c)
// Output:
// a := b[1:]
```

### Line
Line inserts a blank line.

### Clone
Be careful when passing *Statement. Consider the following... 

```go
a := Id("a")
c := Block(
	a.Call(),
	a.Call(),
)
fmt.Printf("%#v", c)
// Output:
// {
// 	a()()
// 	a()()
// }
```

Id("a") returns a *Statement, which the Call() method appends to twice. To 
avoid this, use Clone. Clone makes a copy of the Statement, so further tokens can be appended
without affecting the original.  

```go
a := Id("a")
c := Block(
	a.Clone().Call(),
	a.Clone().Call(),
)
fmt.Printf("%#v", c)
// Output:
// {
// 	a()
// 	a()
// }
```

### Cgo
The cgo "C" pseudo-package is a special case, and always renders without a package alias. The 
import can be added with `Qual`, `Anon` or by supplying a preamble. The preamble is added with 
`File.CgoPreamble` which has the same semantics as [Comment](#comments). If a preamble is provided, 
the import is separated, and preceded by the preamble. 

```go
f := NewFile("a")
f.CgoPreamble(`#include <stdio.h>
#include <stdlib.h>

void myprint(char* s) {
printf("%s\n", s);
}
`)
f.Func().Id("init").Params().Block(
	Id("cs").Op(":=").Qual("C", "CString").Call(Lit("Hello from stdio\n")),
	Qual("C", "myprint").Call(Id("cs")),
	Qual("C", "free").Call(Qual("unsafe", "Pointer").Parens(Id("cs"))),
)
fmt.Printf("%#v", f)
// Output:
// package a
//
// import "unsafe"
//
// /*
// #include <stdio.h>
// #include <stdlib.h>
//
// void myprint(char* s) {
// 	printf("%s\n", s);
// }
// */
// import "C"
//
// func init() {
// 	cs := C.CString("Hello from stdio\n")
// 	C.myprint(cs)
// 	C.free(unsafe.Pointer(cs))
// }
```  

# File
[Identifiers](#identifiers) [Keywords](#keywords) [Operators](#operators) [Braces](#braces) [Parentheses](#parentheses) [Control flow](#control-flow) [Collections](#collections) [Literals](#literals) [Comments](#comments) [Helpers](#helpers) [Misc](#misc) **File**

File represents a single source file. Package imports are managed
automaticaly by File.

### NewFile
NewFile Creates a new file, with the specified package name.

### NewFilePath
NewFilePath creates a new file while specifying the package path - the
package name is inferred from the path.

### NewFilePathName
NewFilePathName creates a new file with the specified package path and name.

```go
f := NewFilePathName("a.b/c", "main")
f.Func().Id("main").Params().Block(
	Qual("a.b/c", "Foo").Call(),
)
fmt.Printf("%#v", f)
// Output:
// package main
//
// func main() {
// 	Foo()
// }
```

### Save
Save renders the file and saves to the filename provided.

### Render
Render renders the file to the provided writer.

```go
f := NewFile("a")
f.Func().Id("main").Params().Block()
buf := &bytes.Buffer{}
err := f.Render(buf)
if err != nil {
	fmt.Println(err.Error())
} else {
	fmt.Println(buf.String())
}
// Output:
// package a
//
// func main() {}
```

### Anon
Anon adds an anonymous import.

```go
f := NewFile("c")
f.Anon("a")
f.Func().Id("init").Params().Block()
fmt.Printf("%#v", f)
// Output:
// package c
//
// import _ "a"
//
// func init() {}
```

### ImportName
ImportName provides the package name for a path. If specified, the alias will be omitted from the
import block. This is optional. If not specified, a sensible package name is used based on the path
and this is added as an alias in the import block.

```go
f := NewFile("main")

// package a should use name "a"
f.ImportName("github.com/foo/a", "a")

// package b is not used in the code so will not be included
f.ImportName("github.com/foo/b", "b")

f.Func().Id("main").Params().Block(
	Qual("github.com/foo/a", "A").Call(),
)
fmt.Printf("%#v", f)

// Output:
// package main
//
// import "github.com/foo/a"
//
// func main() {
// 	a.A()
// }
```

### ImportNames
ImportNames allows multiple names to be imported as a map. Use the [gennames](gennames) command to
automatically generate a go file containing a map of a selection of package names. 

### ImportAlias
ImportAlias provides the alias for a package path that should be used in the import block. A
period can be used to force a dot-import.

```go
f := NewFile("main")

// package a should be aliased to "b"
f.ImportAlias("github.com/foo/a", "b")

// package c is not used in the code so will not be included
f.ImportAlias("github.com/foo/c", "c")

f.Func().Id("main").Params().Block(
	Qual("github.com/foo/a", "A").Call(),
)
fmt.Printf("%#v", f)

// Output:
// package main
//
// import b "github.com/foo/a"
//
// func main() {
// 	b.A()
// }
```

### Comments
PackageComment adds a comment to the top of the file, above the package
keyword.

HeaderComment adds a comment to the top of the file, above any package
comments. A blank line is rendered below the header comments, ensuring
header comments are not included in the package doc.

CanonicalPath adds a canonical import path annotation to the package clause.

```go
f := NewFile("c")
f.CanonicalPath = "d.e/f"
f.HeaderComment("Code generated by...")
f.PackageComment("Package c implements...")
f.Func().Id("init").Params().Block()
fmt.Printf("%#v", f)
// Output:
// // Code generated by...
//
// // Package c implements...
// package c // import "d.e/f"
//
// func init() {}
```

CgoPreamble adds a cgo preamble comment that is rendered directly before the "C" pseudo-package
import.

### PackagePrefix
If you're worried about generated package aliases conflicting with local variable names, you
can set a prefix here. Package foo becomes {prefix}_foo.

```go
f := NewFile("a")
f.PackagePrefix = "pkg"
f.Func().Id("main").Params().Block(
	Qual("b.c/d", "E").Call(),
)
fmt.Printf("%#v", f)
// Output:
// package a
//
// import pkg_d "b.c/d"
//
// func main() {
// 	pkg_d.E()
// }
```
