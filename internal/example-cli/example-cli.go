// minimal example CLI used for binary size checking

package main

import (
	"github.com/aperturerobotics/cli"
)

func main() {
	_ = (&cli.App{}).Run([]string{""})
}
