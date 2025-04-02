---
search:
  boost: 2
---

When your application has numerous subcommands, organizing them into categories can significantly improve the clarity of the help output. You can assign a category to a command by setting its `Category` field. Commands with the same category will be grouped together in the help text.

<!-- {
  "output": ".*COMMANDS:\\n.*noop[ ]*\\n.*\\n[ ]*template:\\n[ ]*add[ ]*\\n[ ]*remove.*"
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
		Commands: []*cli.Command{
			{
				Name: "noop",
			},
			{
				Name:     "add",
				Category: "template",
			},
			{
				Name:     "remove",
				Category: "template",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

Will include:

```
COMMANDS:
  noop

  template:
    add
    remove
```
