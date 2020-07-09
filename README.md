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
