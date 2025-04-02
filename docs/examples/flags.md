---
search:
  boost: 2
---

Flags provide ways for users to modify the behavior of your command-line application. Defining and accessing flags is straightforward.

<!-- {
  "output": "Hello Nefertiti"
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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "lang",
				Value: "english",
				Usage: "language for the greeting",
			},
		},
		Action: func(cCtx *cli.Context) error {
			name := "Nefertiti"
			if cCtx.NArg() > 0 {
				name = cCtx.Args().Get(0)
			}
			if cCtx.String("lang") == "spanish" {
				fmt.Println("Hola", name)
			} else {
				fmt.Println("Hello", name)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

You can also bind a flag directly to a variable in your code using the `Destination` field. The flag's value will be automatically parsed and stored in the specified variable. If a default `Value` is also set for the flag, the `Destination` variable will be initialized to this default before command-line arguments are parsed.

<!-- {
  "output": "Hello someone"
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
	var language string

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "lang",
				Value:       "english",
				Usage:       "language for the greeting",
				Destination: &language,
			},
		},
		Action: func(cCtx *cli.Context) error {
			name := "someone"
			if cCtx.NArg() > 0 {
				name = cCtx.Args().Get(0)
			}
			if language == "spanish" {
				fmt.Println("Hola", name)
			} else {
				fmt.Println("Hello", name)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

See the [Go Reference](https://pkg.go.dev/github.com/aperturerobotics/cli) for a full list of available flag types.

For boolean flags (`BoolFlag`), you can use the `Count` field to track how many times a flag is provided. This is useful for flags like `-v` for verbosity.

> Note: To support combining short boolean flags like `-vvv`, you must set `UseShortOptionHandling: true` on your `App` or `Command`. See the [Combining Short Options](./combining-short-options/) example for details.

<!-- {
  "args": ["&#45;&#45;f", "&#45;&#45;f", "&#45;fff",  "&#45;f"],
  "output": "count 6"
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
	var count int

	app := &cli.App{
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "foo",
				Usage:       "foo greeting",
				Aliases:     []string{"f"},
				Count: &count,
			},
		},
		Action: func(cCtx *cli.Context) error {
			fmt.Println("count", count)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

#### Placeholder Values

You can indicate a placeholder for a flag's value directly within the `Usage` string. This helps clarify what kind of value the flag expects in the help output. Placeholders are denoted using backticks (` `` `).

For example this:

<!-- {
  "args": ["&#45;&#45;help"],
  "output": "&#45;&#45;config FILE, &#45;c FILE"
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
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

Will result in help output like:

```
--config FILE, -c FILE   Load configuration from FILE
```

Note that only the first placeholder is used. Subsequent back-quoted words will
be left as-is.

#### Alternate Names

Flags can have multiple names, often a longer descriptive name and a shorter alias. You can define aliases using the `Aliases` field, which accepts a slice of strings.

<!-- {
  "args": ["&#45;&#45;help"],
  "output": "&#45;&#45;lang value, &#45;l value.*language for the greeting.*default: \"english\""
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
			&cli.StringFlag{
				Name:    "lang",        // Primary name
				Aliases: []string{"l"}, // Alternate names
				Value:   "english",
				Usage:   "language for the greeting",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

That flag can then be set with `--lang spanish` or `-l spanish`. Providing both forms (e.g., `--lang spanish -l spanish`) in the same command invocation will result in an error.

#### Multiple Values per Single Flag

Slice flags allow users to specify a flag multiple times, collecting all provided values into a slice. Available types include:

- `Int64SliceFlag`
- `IntSliceFlag`
- `StringSliceFlag`

<!-- {
  "args": ["&#45;&#45;greeting", "Hello", "&#45;&#45;greeting", "Hola"],
  "output": "Hello, Hola"
} -->
```go
package main
import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aperturerobotics/cli"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  "greeting",
				Usage: "Pass multiple greetings",
			},
		},
		Action: func(cCtx *cli.Context) error {
			fmt.Println(strings.Join(cCtx.StringSlice("greeting"), `, `))
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

To pass multiple values, the user repeats the flag, e.g., `--greeting Hello --greeting Hola`.

#### Grouping

You can associate a category for each flag to group them together in the help output, e.g:

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
				Name:     "silent",
				Aliases:  []string{"s"},
				Usage:    "no messages",
				Category: "Miscellaneous:",
			},
			&cli.BoolFlag{
				Name:     "perl-regexp",
				Aliases:  []string{"P"},
				Usage:    "PATTERNS are Perl regular expressions",
				Category: "Pattern selection:",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

Will result in help output like:

```
GLOBAL OPTIONS:
   Miscellaneous:

   --silent, -s  no messages (default: false)

   Pattern selection:

   --perl-regexp, -P  PATTERNS are Perl regular expressions (default: false)
```

#### Ordering

Flags for the application and commands are shown in the order they are defined.
However, it's possible to sort them from outside this library by using `FlagsByName`
or `CommandsByName` with `sort`.

For example this:

<!-- {
  "args": ["&#45;&#45;help"],
  "output": ".*Load configuration from FILE\n.*Language for the greeting.*"
} -->
```go
package main

import (
	"log"
	"os"
	"sort"

	"github.com/aperturerobotics/cli"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "lang",
				Aliases: []string{"l"},
				Value:   "english",
				Usage:   "Language for the greeting",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "complete",
				Aliases: []string{"c"},
				Usage:   "complete a task on the list",
				Action: func(*cli.Context) error {
					return nil
				},
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a task to the list",
				Action: func(*cli.Context) error {
					return nil
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

Will result in help output like:

```
--config FILE, -c FILE  Load configuration from FILE
--lang value, -l value  Language for the greeting (default: "english")
```

#### Values from the Environment

You can specify environment variables that can provide default values for flags using the `EnvVars` field (a slice of strings).

<!-- {
  "args": ["&#45;&#45;help"],
  "output": "language for the greeting.*APP_LANG"
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
			&cli.StringFlag{
				Name:    "lang",
				Aliases: []string{"l"},
				Value:   "english",
				Usage:   "language for the greeting",
				EnvVars: []string{"APP_LANG"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

If `EnvVars` contains multiple variable names, the library uses the value of the first environment variable found to be set.

<!-- {
  "args": ["&#45;&#45;help"],
  "output": "language for the greeting.*LEGACY_COMPAT_LANG.*APP_LANG.*LANG"
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
			&cli.StringFlag{
				Name:    "lang",
				Aliases: []string{"l"},
				Value:   "english",
				Usage:   "language for the greeting",
				EnvVars: []string{"LEGACY_COMPAT_LANG", "APP_LANG", "LANG"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

If a flag's `Value` field is not explicitly set, but a corresponding environment variable from `EnvVars` is found, the environment variable's value will be used as the default and shown in the help text.

#### Values from Files

Similarly to environment variables, you can specify a file path from which to read a flag's default value using the `FilePath` field.

<!-- {
  "args": ["&#45;&#45;help"],
  "output": "password for the mysql database"
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
			&cli.StringFlag{
				Name:     "password",
				Aliases:  []string{"p"},
				Usage:    "password for the mysql database",
				FilePath: "/etc/mysql/password",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

Note that default values sourced from environment variables (`EnvVars`) take precedence over those sourced from files (`FilePath`). See the full precedence order below.

#### Required Flags

You can enforce that a flag must be provided by the user by setting its `Required` field to `true`. If a required flag is missing, the application will print an error message and exit.

Take for example this app that requires the `lang` flag:

<!-- {
  "error": "Required flag \"lang\" not set"
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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "lang",
				Value:    "english",
				Usage:    "language for the greeting",
				Required: true,
			},
		},
		Action: func(cCtx *cli.Context) error {
			output := "Hello"
			if cCtx.String("lang") == "spanish" {
				output = "Hola"
			}
			fmt.Println(output)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

If the app is run without the `lang` flag, the user will see the following message

```
Required flag "lang" not set
```

#### Default Values for Help Output

In cases where a flag's default value (`Value` field) is dynamic or complex to represent (e.g., a randomly generated port), you can specify custom text to display as the default in the help output using the `DefaultText` field. This text is purely for display purposes and doesn't affect the actual default value used by the application.

For example this:

<!-- {
  "args": ["&#45;&#45;help"],
  "output": "&#45;&#45;port value"
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
			&cli.IntFlag{
				Name:        "port",
				Usage:       "Use a randomized port",
				Value:       0,
				DefaultText: "random",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

Will result in help output like:

```
--port value  Use a randomized port (default: random)
```

#### Precedence

The order of precedence for determining a flag's value is as follows (from highest to lowest):

1.  Value provided on the command line by the user.
2.  Value from the first set environment variable listed in `EnvVars`.
3.  Value read from the file specified in `FilePath`.
4.  Default value specified in the `Value` field of the flag definition.

#### Flag Actions

You can attach an `Action` function directly to a flag definition. This function will be executed *after* the flag's value has been parsed from the command line, environment, or file. This is useful for performing validation or other logic specific to that flag. The action receives the `*cli.Context` and the parsed value of the flag.

<!-- {
  "args": ["&#45;&#45;port","70000"],
  "error": "Flag port value 70000 out of range[0-65535]"
} -->
```go
package main

import (
	"log"
	"os"
	"fmt"

	"github.com/aperturerobotics/cli"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				Usage:       "Use a randomized port",
				Value:       0,
				DefaultText: "random",
				Action: func(ctx *cli.Context, v int) error {
					if v >= 65536 {
						return fmt.Errorf("Flag port value %v out of range[0-65535]", v)
					}
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

Will result in help output like:

```
Flag port value 70000 out of range[0-65535]
```
