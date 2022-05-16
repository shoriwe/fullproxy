# plasma

## Documentation

You can find documentation in:

[Documentation](https://shoriwe.github.io)

## Description

**`plasma`** is a dynamic programming language highly inspired in **`ruby`** syntax and semantics with interfaces and
design focused in application embedding.

## Try it

You can have a working interpreter by compiling `cmd/plasma`.

- You can compile a binary with (using **`Go-1.16`**)

```shell
go install github.com/shoriwe/gplasma/cmd/plasma@latest
```

```
...>plasma [MODE] [FLAG [FLAG [FLAG]]] [PROGRAM [PROGRAM [PROGRAM]]]

[+] Notes
        - No PROGRAM arguments will spawn a REPL

[+] Flags
        -h, --help              Show this help message

[+] Modes
        module          tool to install, uninstall and initialize modules

[+] Environment Variables
        NoColor -> TRUE or FALSE                Disable color printing for this CLI
        SitePackages -> PATH            This is the path to the Site-Packages of the running VM; Default is PATH/TO/PLASMA/EXECUTABLE/site-packages
```

## Features

### Embedding

**`plasma`** was designed to be embedded in other go applications, you should do it like:

```go
package main

import (
	"fmt"
	"github.com/shoriwe/gplasma"
	"github.com/shoriwe/gplasma/pkg/std/features/importlib"
	"os"
)

var (
	files                   []string
	virtualMachine          *gplasma.VirtualMachine
	sitePackagesPath        = "site-packages"
)

// Setup the vm based on the options
func setupVM() {
	virtualMachine = gplasma.NewVirtualMachine()
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = "."
	}
	importSystem := importlib.NewImporter()
	// Load Default modules to use with the VM
	importSystem.LoadModule(regex.Regex)
	//
	virtualMachine.LoadFeature(
		importSystem.Result(
			importlib.NewRealFileSystem(sitePackagesPath),
			importlib.NewRealFileSystem(currentDir),
		),
	)
}

func program() {
	setupVM()
	for _, filePath := range files {
		fileHandler, openError := os.Open(filePath)
		if openError != nil {
			_, _ = fmt.Fprintf(os.Stderr, openError.Error())
			os.Exit(1)
		}
		content, readingError := io.ReadAll(fileHandler)
		if readingError != nil {
			_, _ = fmt.Fprintf(os.Stderr, readingError.Error())
			os.Exit(1)
		}
		result, success := virtualMachine.ExecuteMain(string(content))
		if !success {
			_, _ = fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", color.RedString("-"), result.TypeName(), result.String)
			os.Exit(1)
		}
	}
}
```

In the future there will be a simpler way to embed it in your application, which shouldn't break the one provided
before.

## Notable Differences

The major difference between **`ruby`** and **`plasma`** is that in the first the last expression in a function will be
returned without specifying the keyboard `return` but in **`plasma`** you should.

Another one will be that function calls, will always need parentheses to be executed, other way their will be evaluated
as objects.

Example:

This example shows a valid **`ruby`** code that returns from a function a string.

```ruby
def hello()
    "Hello World"
end

puts hello
```

But in **`plasma`** you should code it something like:

```ruby
def hello()
    return "Hello World" # Notice that here is used the keyboard "return"
end

println(hello())
```

# Useful references

This where useful references that made this project possible.

- [BNF grammar](https://ruby-doc.org/docs/ruby-doc-bundle/Manual/man-1.4/yacc.html)
- [Syntax Documentation](https://ruby-doc.org/docs/ruby-doc-bundle/Manual/man-1.4/syntax.html)
