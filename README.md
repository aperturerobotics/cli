# cli

[![Run Tests](https://github.com/aperturerobotics/cli/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/aperturerobotics/cli/actions/workflows/tests.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/aperturerobotics/cli.svg)](https://pkg.go.dev/github.com/aperturerobotics/cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/aperturerobotics/cli)](https://goreportcard.com/report/github.com/aperturerobotics/cli)
[![codecov](https://codecov.io/gh/aperturerobotics/cli/branch/main/graph/badge.svg)](https://codecov.io/gh/aperturerobotics/cli)

`aperturerobotics/cli` is a **fork** of the popular `urfave/cli` v2 package for building command line apps in Go. The goal is to enable developers to write fast and distributable command line applications in an expressive way, while minimizing dependencies and maximizing compatibility.

Key differences from `urfave/cli`:

1.  **Slim and Reflection-Free:**
    *   Removed `reflect` usage for smaller binaries and better performance.
    *   Tinygo compatible.
    *   Removed documentation generators.
    *   Removed altsrc package to focus on CLI handling only.
2.  **Stability:** Try to maintain backward compatibility as much as possible.

## Installation

Using this package requires a working Go environment. [See the install instructions for Go](http://golang.org/doc/install.html).

Go Modules are required when using this package. [See the go blog guide on using Go Modules](https://blog.golang.org/using-go-modules).

```sh
go get github.com/aperturerobotics/cli
```

## Getting Started

Here's a simple example to get you started:

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aperturerobotics/cli"
)

func main() {
	app := &cli.App{
		Name:  "greet",
		Usage: "a simple greeter application",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Value:   "world",
				Usage:   "who to greet",
				EnvVars: []string{"GREET_NAME"},
			},
		},
		// The action for the root command (optional)
		Action: func(ctx *cli.Context) error {
			name := ctx.String("name")
			fmt.Printf("Hello %s!\n", name)
			return nil
		},
		// Define subcommands
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "add a task to the list",
				// Action for the 'add' subcommand
				Action: func(ctx *cli.Context) error {
					fmt.Println("added task: ", ctx.Args().First())
					return nil
				},
			},
			{
				Name:  "complete",
				Usage: "complete a task on the list",
				// Action for the 'complete' subcommand
				Action: func(ctx *cli.Context) error {
					fmt.Println("completed task: ", ctx.Args().First())
					return nil
				},
			},
		},
	}

	// Run the application
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// Try running this:
// GREET_NAME=everyone ./greet --name someone add some-task
// ./greet complete --help
```

Running this provides basic command functionality, including help text generation, flag parsing, environment variable handling, and subcommand routing. You can easily add more flags, subcommands, and complex actions.

## Documentation

Full documentation and examples are available in the [`./docs`](./docs) directory and online at <https://cli.aperture.app>.

*   [Getting Started](./docs/getting-started.md)
*   [Examples](./docs/examples/)

## License

This fork retains the original MIT license from `urfave/cli`. See [`LICENSE`](./LICENSE).
