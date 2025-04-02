---
search:
  boost: 2
---

The library automatically generates help text for your application, commands, and subcommands. By default, this help text is displayed when the user provides the `-h` or `--help` flag (defined by `cli.HelpFlag`). When this flag is detected, the help text is printed, and the application exits.

#### Customization

You have several ways to customize the generated help output:

1.  **Modify Templates:** The default Go text templates used for generation (`cli.AppHelpTemplate`, `cli.CommandHelpTemplate`, `cli.SubcommandHelpTemplate`) are exported variables. You can modify them directly, for example, by appending extra information or completely replacing them with your own template strings.
2.  **Replace Help Printer:** For complete control over rendering, you can replace the `cli.HelpPrinter` function. This function receives the output writer, the template string, and the data object (like `*cli.App` or `*cli.Command`) and is responsible for executing the template or generating help in any way you choose.

Here are examples of these customization techniques:

<!-- {
  "output": "Ha HA.  I pwnd the help!!1"
} -->
```go
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/aperturerobotics/cli"
)

func main() {
	// EXAMPLE: Append to an existing template
	cli.AppHelpTemplate = fmt.Sprintf(`%s

WEBSITE: http://awesometown.example.com

SUPPORT: support@awesometown.example.com

`, cli.AppHelpTemplate)

	// EXAMPLE: Override a template
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
`

	// EXAMPLE: Replace the `HelpPrinter` func
	cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
		fmt.Println("Ha HA.  I pwnd the help!!1")
	}

	(&cli.App{}).Run(os.Args)
}
```

You can also change the flag used to trigger the help display (instead of the default `-h/--help`) by assigning a different `cli.Flag` implementation to the `cli.HelpFlag` variable:

<!-- {
  "args": ["&#45;&#45halp"],
  "output": "haaaaalp.*HALP"
} -->
```go
package main

import (
	"os"

	"github.com/aperturerobotics/cli"
)

func main() {
	cli.HelpFlag = &cli.BoolFlag{
		Name:    "haaaaalp",
		Aliases: []string{"halp"},
		Usage:   "HALP",
		EnvVars: []string{"SHOW_HALP", "HALPPLZ"},
	}

	(&cli.App{}).Run(os.Args)
}
```
