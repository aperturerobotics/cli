// Package cli builds command line Go applications from an App and a tree of
// Commands, each with its own Flags, Before, Action and After.
//
// The smallest application is:
//
//	func main() {
//		(&cli.App{}).Run(os.Args)
//	}
//
// One that does something:
//
//	func main() {
//		app := &cli.App{
//			Name:  "greet",
//			Usage: "say a greeting",
//			Action: func(c *cli.Context) error {
//				fmt.Println("Greetings")
//				return nil
//			},
//		}
//		app.Run(os.Args)
//	}
//
// The package is a fork of urfave/cli v2 without reflection, documentation
// generators or the altsrc package, so it builds small binaries and compiles
// under TinyGo.
//
// Parsing differs from urfave/cli in three ways. Running a command first
// parses every flag on the path from the root to the selected command; only
// then do the Before hooks and flag actions run from the root down, followed
// by the Action and the After hooks in reverse order, so a root Before sees a
// root flag given after a subcommand. A flag given to a subcommand but
// defined only by an ancestor is set on the nearest ancestor that defines it.
// A command without subcommands parses flags anywhere among its arguments,
// as in "create name --flag", and "--" keeps every later argument
// positional; a command with subcommands stops at its first argument, which
// names the subcommand.
//
// Each App installs its own copy of the help command, so several Apps may
// run concurrently.
//
// Programs built on this package group related flags in a typed args record
// whose BuildFlags method returns its flags bound to the record's fields, and
// whose Run(*cli.Context) method serves as the command's Action, so a command
// reads typed fields rather than looking flags up by name.
package cli

//go:generate make -C cmd/urfave-cli-genflags run
