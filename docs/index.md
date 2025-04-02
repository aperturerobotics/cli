# Welcome to aperturerobotics/cli

[![Go Reference](https://pkg.go.dev/badge/github.com/aperturerobotics/cli.svg)](https://pkg.go.dev/github.com/aperturerobotics/cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/aperturerobotics/cli)](https://goreportcard.com/report/github.com/aperturerobotics/cli)

`aperturerobotics/cli` is a powerful **fork** of the popular `urfave/cli` v2 package, designed for building command-line applications in Go with a focus on simplicity and performance.

Key differences from `urfave/cli`:

1.  **Slim and Reflection-Free:**
    *   Removed `reflect` usage for smaller binaries and better performance.
    *   Tinygo compatible.
    *   Removed documentation generators.
    *   Removed altsrc package to focus on CLI handling only.
2.  **Stability:** Try to maintain backward compatibility as much as possible.

Documentation:

- [Getting Started](./getting-started/)
- [Examples](./examples/)

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
```

This example demonstrates basic command functionality, including automatic help text generation, flag parsing (with environment variable support), and subcommand routing. You can easily extend this foundation by adding more flags, subcommands, and complex actions. See the [full getting started guide](./getting-started/) for more details.

### Supported platforms

cli is tested against multiple versions of Go on Linux, and against the latest
released version of Go on OS X and Windows. This project uses GitHub Actions
for builds. To see our currently supported go versions and platforms, look at
the [github workflow configuration](https://github.com/aperturerobotics/cli/blob/main/.github/workflows/tests.yml).
