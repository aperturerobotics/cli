---
search:
  boost: 2
---

Normally, short options (like `-s`, `-o`, `-m`) are specified individually:

```sh-session
$ cmd -s -o -m "Some message"
```

If you want to allow users to combine multiple short boolean options (and optionally one non-boolean short option at the end), you can enable this behavior by setting `UseShortOptionHandling` to `true` on your `cli.App` or a specific `cli.Command`.

Here's an example:

<!-- {
  "args": ["short", "&#45;som", "Some message"],
  "output": "serve: true\noption: true\nmessage: Some message\n"
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
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			{
				Name:  "short",
				Usage: "complete a task on the list",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "serve", Aliases: []string{"s"}},
					&cli.BoolFlag{Name: "option", Aliases: []string{"o"}},
					&cli.StringFlag{Name: "message", Aliases: []string{"m"}},
				},
				Action: func(cCtx *cli.Context) error {
					fmt.Println("serve:", cCtx.Bool("serve"))
					fmt.Println("option:", cCtx.Bool("option"))
					fmt.Println("message:", cCtx.String("message"))
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

If your program has any number of bool flags such as `serve` and `option`, and
optionally one non-bool flag `message`, with the short options of `-s`, `-o`,
and `-m` respectively, setting `UseShortOptionHandling` will also support the
following syntax:

```sh-session
$ cmd -som "Some message"
```

**Important:** When `UseShortOptionHandling` is enabled, you cannot define flags that use a single dash followed by multiple characters (e.g., `-option`). This syntax becomes ambiguous with combined short options. Standard double-dash flags (e.g., `--option`) remain unaffected.
