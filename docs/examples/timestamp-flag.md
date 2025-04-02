---
search:
  boost: 2
---

The `TimestampFlag` allows users to provide date and time values as command-line arguments. You must specify the expected format using the `Layout` field. The layout string follows the rules of Go's `time.Parse` function (refer to the [`time` package documentation](https://golang.org/pkg/time/#Parse) for details on defining layouts).

<!-- {
  "args": ["&#45;&#45;meeting", "2019-08-12T15:04:05"],
  "output": "2019\\-08\\-12 15\\:04\\:05 \\+0000 UTC"
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
			&cli.TimestampFlag{Name: "meeting", Layout: "2006-01-02T15:04:05"},
		},
		Action: func(cCtx *cli.Context) error {
			fmt.Printf("%s", cCtx.Timestamp("meeting").String())
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

In this example the flag could be used like this:

```sh-session
$ myapp --meeting 2019-08-12T15:04:05
```

If the specified `Layout` does not include timezone information, the parsed time will be in UTC by default. You can specify a different default timezone (like the system's local time) using the `Timezone` field:

```go
package main

import (
	"log"
	"os"
	"time"

	"github.com/aperturerobotics/cli"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.TimestampFlag{Name: "meeting", Layout: "2006-01-02T15:04:05", Timezone: time.Local},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

(time.Local contains the system's local time zone.)

Side note: quotes may be necessary around the date depending on your layout (if
you have spaces for instance)
