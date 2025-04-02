---
search:
  boost: 2
---

Accessing command-line arguments (values passed after the command and flags) is straightforward using the `Args` method on `cli.Context`. Here's an example:

<!-- {
  "output": "Hello \""
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
		Action: func(cCtx *cli.Context) error {
			fmt.Printf("Hello %q", cCtx.Args().Get(0))
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```
