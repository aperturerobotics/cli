---
search:
  boost: 2
---

Getting started with `aperturerobotics/cli` is incredibly simple. A functional command-line application can be created with just a single line of code in your `main()` function.

<!-- {
  "args": ["&#45;&#45;help"],
  "output": "A new cli application"
} -->
```go
package main

import (
	"os"

	"github.com/aperturerobotics/cli"
)

func main() {
	(&cli.App{}).Run(os.Args)
}
```

While this app runs and displays basic help text, it doesn't perform any actions yet. Let's enhance it by adding a name, usage description, and an action to execute:

<!-- {
  "output": "boom! I say!"
} -->
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
		Name:  "boom",
		Usage: "make an explosive entrance",
		Action: func(*cli.Context) error {
			fmt.Println("boom! I say!")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

Running this improved version provides a more useful application with clear help text and a defined action.

## Adding Flags

Let's make our application more interactive by adding a flag. Flags allow users to pass options to the command. We'll add a `--name` flag to specify who to greet.

<!-- {
  "args": ["&#45;&#45;name", "Alice"],
  "output": "boom! Hello Alice!"
} -->
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
		Name:  "boom",
		Usage: "make an explosive entrance",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Value: "world", // Default value
				Usage: "who to greet",
			},
		},
		Action: func(cCtx *cli.Context) error {
			name := cCtx.String("name")
			fmt.Printf("boom! Hello %s!\n", name)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

Now you can run `go run main.go --name Bob` or simply `go run main.go` to use the default value "world". The help output (`go run main.go --help`) will also automatically include information about the new flag.

## Adding Subcommands

For more complex applications, you might want different actions grouped under subcommands (like `git commit` or `docker ps`). Let's add a `greet` subcommand to our `boom` app.

<!-- {
  "args": ["greet", "&#45;&#45;name", "Carol"],
  "output": "Hello Carol!"
} -->
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
		Name:  "boom",
		Usage: "make an explosive entrance",
		Commands: []*cli.Command{
			{
				Name:  "greet",
				Usage: "say hello",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Value: "world",
						Usage: "who to greet",
					},
				},
				Action: func(cCtx *cli.Context) error {
					name := cCtx.String("name")
					fmt.Printf("Hello %s!\n", name)
					return nil
				},
			},
			// Add more subcommands here
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

Now, the main action is performed by running `go run main.go greet --name Dave`. Running `go run main.go --help` will show the available subcommands.

This tutorial covers the basics of creating a CLI application, adding flags, and organizing functionality with subcommands using `aperturerobotics/cli`. Explore the [Examples](./examples/) section for more advanced use cases.
