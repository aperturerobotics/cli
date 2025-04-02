---
search:
  boost: 2
---

By default, `App.Run` does not call `os.Exit`. If your application's `Action` (or subcommand actions) returns `nil`, the process exits with code `0`. To exit with a specific non-zero code, return an error that implements the `cli.ExitCoder` interface. The `cli.Exit` helper function is provided for convenience. If using `cli.MultiError`, the exit code will be determined by the first `cli.ExitCoder` found within the wrapped errors.

Here's an example using `cli.Exit`:
<!-- {
  "error": "Ginger croutons are not in the soup"
} -->
```go
package main

import (
	"log"
	"os"

	"github.com/aperturerobotics/cli"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "ginger-crouton",
				Usage: "is it in the soup?",
			},
		},
		Action: func(ctx *cli.Context) error {
			if !ctx.Bool("ginger-crouton") {
				return cli.Exit("Ginger croutons are not in the soup", 86)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```
