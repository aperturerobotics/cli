---
search:
  boost: 2
---

Similar to the help flag, the library provides built-in support for a version flag. By default, `-v` or `--version` (defined by `cli.VersionFlag`) triggers the display of the application's version string (set in `App.Version`). The `cli.VersionPrinter` function handles the printing, and then the application exits.

#### Customization

You can customize the version flag behavior:

1.  **Change the Flag:** Assign a different `cli.Flag` implementation to the `cli.VersionFlag` variable to change which flag(s) trigger the version display.

<!-- {
  "args": ["&#45;&#45print-version"],
  "output": "partay version v19\\.99\\.0"
} -->
```go
package main

import (
	"os"

	"github.com/aperturerobotics/cli"
)

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "print-version",
		Aliases: []string{"V"},
		Usage:   "print only the version",
	}

	app := &cli.App{
		Name:    "partay",
		Version: "v19.99.0",
	}
	app.Run(os.Args)
}
```

2.  **Customize Output:** Replace the `cli.VersionPrinter` function to control how the version information is formatted and printed. This is useful for including additional details like build revision numbers.

<!-- {
  "args": ["&#45;&#45version"],
  "output": "version=v19\\.99\\.0 revision=fafafaf"
} -->
```go
package main

import (
	"fmt"
	"os"

	"github.com/aperturerobotics/cli"
)

var (
	Revision = "fafafaf"
)

func main() {
	cli.VersionPrinter = func(cCtx *cli.Context) {
		fmt.Printf("version=%s revision=%s\n", cCtx.App.Version, Revision)
	}

	app := &cli.App{
		Name:    "partay",
		Version: "v19.99.0",
	}
	app.Run(os.Args)
}
```
