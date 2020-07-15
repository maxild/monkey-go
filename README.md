## Code organization

In a GIT repo GO code can live in the root dir, or in 3 different subdirs.

There are no offcial conventions for organizing code/packages in GO. This is in
contrast with languages such as [Rust](https://doc.rust-lang.org/cargo/guide/project-layout.html).

Go has some conventions using dirs named `cmd`, `pkg` and `internal`. Of course
if the project is simple without many file (packages) then we could use a flat
structure.

* `cmd`: Dir contains command packages.
* `pkg`: Dir contains public packages.
* `internal`: Dir contains private packages.

The `cmd` dir can be a good idea, if project contains many commands
(executables), or to separate the `main.go` files from the root dir, where many
other files can live.

When you want to build or run something, it will look like `go run cmd/binaryname/main.go`.

An import of a path containing the element “internal” is disallowed
if the importing code is outside the tree rooted at the parent of the “internal” directory.

If you put a package inside an internal directory, then other packages can’t
import it unless they share a common ancestor. Internal packages enable you
to export code for reuse in your project while reducing your public API.

```
cmd/
  binaryname/
    main.go # a small file glueing things together
internal/
  data/
    types.go
    types_test.go # unit tests are right here
    (...)
pkg/
  api/
    types.go  # REST API input and output types
test/
  smoketest.py
ui/
  index.html
README.md
Makefile
(...)
```

The test directory does not contain Go tests! Unit tests live right besides
the code they are supposed to test. Instead, this is the place to put scripts
for external blackbox and smoke tests. I like to use Python for my scripting
needs, as you see.

### Wait, what about GOPATH, src and all that?
One of my favorite changes in Go 1.11, was the addition of Go modules. They
take godep’s job, and finally free your projects from having to be placed
inside of the hierarchical complex folder structure of GOPATH. It’s super neat!

### Want to learn more
There is an [example](https://github.com/golang-standards/project-layout)
project online.

## Pratt Parser Resources

* [Introduction to Pratt Parsing and its terminology](https://abarker.github.io/typped/pratt_parsing_intro.html) Python
* [How Desmos uses Pratt Parsers](https://engineering.desmos.com/articles/pratt-parser/)
* [Pratt Parsing](https://dev.to/jrop/pratt-parsing)

* [Top Down Operator Precedence]() (2007) by Douglas Crockford
* [Simple Top-Down Parsing in Python](https://effbot.org/zone/simple-top-down-parsing.htm) (2008) by Fredrik Lund
* [Top-Down operator precedence parsing](https://eli.thegreenplace.net/2010/01/02/top-down-operator-precedence-parsing) (2010) by Eli Bendersky) by Eli Bendersky
* [Pratt Parsers: Expression Parsing Made Easy](https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/) (2011) by Bob Nystrom

### Great overview of [Andy Chu](http://andychu.net/)

This [overview](http://www.oilshell.org/blog/2017/03/31.html) is great.

* [Pratt Parsing and Precedence Climbing Are the Same Algorithm](http://www.oilshell.org/blog/2016/11/01.html)
* [Review of Pratt/TDOP Parsing Tutorials](http://www.oilshell.org/blog/2016/11/02.html)
* [Pratt Parsing Without Prototypal Inheritance, Global Variables, Virtual Dispatch, or Java](http://www.oilshell.org/blog/2016/11/03.html)
* [Pratt Parsing Demo](https://github.com/andychu/pratt-parsing-demo) in Python

* [Parsing Expressions by Recursive Descent](http://www.engr.mun.ca/~theo/Misc/exp_parsing.htm) by Thedore Norvell
* [From Precedence Climbing to Pratt Parsing](https://www.engr.mun.ca/~theo/Misc/pratt_parsing.htm) by Theodore Norvell

* [Simple but Powerful Pratt Parsing](https://matklad.github.io/2020/04/13/simple-but-powerful-pratt-parsing.html)
  in Rust

* [Jean-Marc Bourguet](https://github.com/bourguet/operator_precedence_parsing)

## Other resources

* [Crafting Interpreters](http://craftinginterpreters.com/) by Bob Nystrom.
* [Let's build a compiler](https://generalproblem.net/lets_build_a_compiler/01-starting-out/) by Noah Zentzis.
* [Write Yourself a Scheme in 48 Hours](https://en.wikibooks.org/wiki/Write_Yourself_a_Scheme_in_48_Hours)
* [Let’s Build A Simple Interpreter](https://ruslanspivak.com/lsbasi-part1/) in
  Python.

  Also when learning Rust (later on) this [repo](https://github.com/pauldix/monkey-rust)
  (and the guy behind it) can be used to help. He also wrote a blog [post](https://www.influxdata.com/blog/rust-can-be-difficult-to-learn-and-frustrating-but-its-also-the-most-exciting-thing-in-software-development-in-a-long-time/)
  about the subject.

## Rust developer environment links

* [How I Start Nix](https://christine.website/blog/how-i-start-nix-2020-03-08)
* [Rust Overlay](https://github.com/mozilla/nixpkgs-mozilla#rust-overlay)
* [Nix Rust Development](https://duan.ca/2020/05/07/nix-rust-development/)
