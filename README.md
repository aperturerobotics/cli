# cli

[![Run Tests](https://github.com/aperturerobotics/cli/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/aperturerobotics/cli/actions/workflows/tests.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/aperturerobotics/cli.svg)](https://pkg.go.dev/github.com/aperturerobotics/cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/aperturerobotics/cli)](https://goreportcard.com/report/github.com/aperturerobotics/cli)
[![codecov](https://codecov.io/gh/aperturerobotics/cli/branch/main/graph/badge.svg)](https://codecov.io/gh/aperturerobotics/cli)

`aperturerobotics/cli` is a **fork** of the popular `urfave/cli` v2 package for building command line apps in Go. The goal is to enable developers to write fast and distributable command line applications in an expressive way, while minimizing dependencies and maximizing compatibility.

Key differences from `urfave/cli`:

1.  **Reflection-Free:** All features relying on `reflect` have been removed. This makes the library suitable for environments where reflection is undesirable or restricted, potentially improving performance and reducing binary size.
2.  **Selective v3 Backports:** Some ergonomic improvements from `urfave/cli` v3 have been incorporated:
    *   `cli.App` has been renamed to `cli.Command` for better semantic clarity, especially in applications with subcommands. The top-level application is now simply the root `Command`.
    *   Action handlers (`Action`, `Before`, `After`, etc.) now accept `context.Context` as their first argument, enabling easier integration with context-aware Go applications.
3.  **Stability:** We strive to maintain backward compatibility and avoid breaking changes.

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
	"os"

	"github.com/aperturerobotics/cli"
)

func main() {
	cmd := &cli.Command{
		Name:  "greet",
		Usage: "say hello",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Value:   "world",
				Usage:   "who to greet",
				EnvVars: []string{"GREET_NAME"},
			},
		},
		Action: func(ctx *cli.Context) error {
			name := ctx.String("name")
			fmt.Printf("Hello %s!\n", name)
			return nil
		},
		Subcommands: []*cli.Command{
			{
				Name:  "add",
				Usage: "add a task to the list",
				Action: func(ctx *cli.Context) error {
					fmt.Println("added task: ", ctx.Args().First())
					return nil
				},
			},
			{
				Name:  "complete",
				Usage: "complete a task on the list",
				Action: func(ctx *cli.Context) error {
					fmt.Println("completed task: ", ctx.Args().First())
					return nil
				},
			},
		},
	}

	if err := cmd.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
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
