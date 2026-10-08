package cli

import (
	"flag"
	"fmt"
	"slices"
	"sort"
	"strings"
)

// Command is a subcommand for a cli.App.
type Command struct {
	// Name is the name of the command.
	Name string
	// Aliases are other names for the command.
	Aliases []string
	// Usage is a short description of the usage of this command.
	Usage string
	// UsageText is custom text for the USAGE section of help.
	UsageText string
	// Description is a longer explanation of how the command works.
	Description string
	// Args reports whether this command supports arguments.
	Args bool
	// ArgsUsage is a short description of the arguments of this command.
	ArgsUsage string
	// Category is the category the command is part of.
	Category string
	// BashComplete is called when checking for bash command completions.
	BashComplete BashCompleteFunc
	// Before runs after every flag on the path to the selected command is
	// parsed, and before the Before of any subcommand. A non-nil error stops
	// the run before any subcommand hook or Action.
	Before BeforeFunc
	// After runs after the selected command and every subcommand After hook
	// finish. It runs even when the Action panics.
	After AfterFunc
	// Action is called when this command is invoked.
	Action ActionFunc
	// OnUsageError is called when a usage error occurs.
	OnUsageError OnUsageErrorFunc
	// Subcommands are the child commands.
	Subcommands []*Command
	// Flags are the flags to parse.
	Flags          []Flag
	flagCategories FlagCategories
	// SkipFlagParsing treats all flags as normal arguments.
	SkipFlagParsing bool
	// HideHelp hides the built-in help command and help flag.
	HideHelp bool
	// HideHelpCommand hides the built-in help command but keeps the help
	// flag. It is ignored when HideHelp is set.
	HideHelpCommand bool
	// Hidden hides this command from help and completion.
	Hidden bool
	// UseShortOptionHandling lets the user combine several single-character
	// bool flags into one, as in "foobar -ov" for "foobar -o -v".
	UseShortOptionHandling bool

	// HelpName is the full name of the command for help. It defaults to the
	// full command name, including parent commands.
	HelpName        string
	commandNamePath []string

	// CustomHelpTemplate is the text/template for the command help topic.
	CustomHelpTemplate string

	// categories contains the categorized commands, populated on app startup.
	categories CommandCategories

	// isRoot reports whether this is the special root command.
	isRoot bool

	separator separatorSpec
}

// Commands is a list of commands.
type Commands []*Command

// CommandsByName sorts commands by name.
type CommandsByName []*Command

// Len returns the number of commands.
func (c CommandsByName) Len() int {
	return len(c)
}

// Less orders commands lexicographically by name.
func (c CommandsByName) Less(i, j int) bool {
	return lexicographicLess(c[i].Name, c[j].Name)
}

// Swap swaps two commands.
func (c CommandsByName) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// FullName returns the full name of the command. For a subcommand it includes
// the parent commands.
func (c *Command) FullName() string {
	if c.commandNamePath == nil {
		return c.Name
	}
	return strings.Join(c.commandNamePath, " ")
}

// Command returns the subcommand with the given name or alias, or nil.
func (c *Command) Command(name string) *Command {
	for _, sub := range c.Subcommands {
		if sub.HasName(name) {
			return sub
		}
	}
	return nil
}

// setup adds the help command and flag, categorizes the subcommands, and
// fills their help names before the command parses its arguments.
func (c *Command) setup(ctx *Context) {
	// Add the built-in help command and flag.
	if c.Command(helpName) == nil && !c.HideHelp {
		if !c.HideHelpCommand {
			c.Subcommands = append(c.Subcommands, ctx.App.helpCmd)
		}
	}
	if !c.HideHelp && HelpFlag != nil {
		c.appendFlag(ctx.App.helpFlag())
	}

	// Inherit app options and categorize the subcommands for help.
	if ctx.App.UseShortOptionHandling {
		c.UseShortOptionHandling = true
	}
	c.categories = newCommandCategories()
	for _, command := range c.Subcommands {
		c.categories.AddCommand(command.Category, command)
	}
	sort.Sort(c.categories.(*commandCategories))

	// Give each subcommand its help name and flag separator.
	for _, scmd := range c.Subcommands {
		if scmd.HelpName == "" {
			scmd.HelpName = fmt.Sprintf("%s %s", c.HelpName, scmd.Name)
		}
		scmd.separator = c.separator
	}

	// Complete flags by default.
	if c.BashComplete == nil {
		c.BashComplete = DefaultCompleteWithFlags(c)
	}
}

// Run parses arguments, selects the subcommand they name, and runs it. Every
// flag on the path to the selected command is parsed before any Before hook
// runs, so a hook sees a flag given after a subcommand, as in "app sub
// --app-flag". The hooks then run from the root down.
func (c *Command) Run(cCtx *Context, arguments ...string) error {
	return c.run(cCtx, nil, arguments)
}

// commandStage is a command on the path to the selected command, with the
// context its flags were parsed into.
type commandStage struct {
	cmd  *Command
	cCtx *Context
}

// run parses this command's arguments, then descends into the selected
// subcommand with this command appended to ancestors. The selected command
// runs the hooks of every stage and its Action.
func (c *Command) run(cCtx *Context, ancestors []commandStage, arguments []string) (err error) {
	// Prepare a subcommand; the app prepares the root.
	if !c.isRoot {
		c.setup(cCtx)
		if err := checkDuplicatedCmds(c); err != nil {
			return err
		}
	}

	// Parse this command's flags.
	a := args(arguments)
	set, err := c.parseFlags(&a, cCtx)
	cCtx.flagSet = set
	if checkCompletions(cCtx) {
		return nil
	}

	// Report a usage error through the hook, or print it with help.
	if err != nil {
		if c.OnUsageError != nil {
			err = c.OnUsageError(cCtx, err, !c.isRoot)
			cCtx.App.handleExitCoder(cCtx, err)
			return err
		}
		_, _ = fmt.Fprintf(cCtx.App.Writer, "%s %s\n\n", "Incorrect Usage:", err.Error())
		if cCtx.App.Suggest {
			if suggestion, err := c.suggestFlagFromError(err, ""); err == nil {
				fmt.Fprintf(cCtx.App.Writer, "%s", suggestion)
			}
		}
		if !c.HideHelp {
			if c.isRoot {
				_ = ShowAppHelp(cCtx)
			} else {
				_ = ShowCommandHelp(cCtx.parentContext, c.Name)
			}
		}
		return err
	}

	// Answer help and version requests, then check required flags.
	if checkHelp(cCtx) {
		return helpCommand.Action(cCtx)
	}
	if c.isRoot && !cCtx.App.HideVersion && checkVersion(cCtx) {
		ShowVersion(cCtx)
		return nil
	}
	if cerr := cCtx.checkRequiredFlags(c.Flags); cerr != nil {
		_ = helpCommand.Action(cCtx)
		return cerr
	}

	// Select the subcommand named by the first argument, falling back to the
	// app's default command.
	var cmd *Command
	args := cCtx.Args()
	if args.Present() {
		name := args.First()
		cmd = c.Command(name)
		if cmd == nil {
			hasDefault := cCtx.App.DefaultCommand != ""
			isFlagName := checkStringSliceIncludes(name, cCtx.FlagNames())

			var isDefaultSubcommand, defaultHasSubcommands bool
			if hasDefault {
				dc := cCtx.App.Command(cCtx.App.DefaultCommand)
				defaultHasSubcommands = len(dc.Subcommands) > 0
				for _, dcSub := range dc.Subcommands {
					if checkStringSliceIncludes(name, dcSub.Names()) {
						isDefaultSubcommand = true
						break
					}
				}
			}

			if isFlagName || (hasDefault && (defaultHasSubcommands && isDefaultSubcommand)) {
				// Use the default command only when it was prepended. Compare
				// contents because Args may be implemented by a pointer.
				argsWithDefault := cCtx.App.argsWithDefaultCommand(args)
				if !slices.Equal(args.Slice(), argsWithDefault.Slice()) {
					cmd = cCtx.App.rootCommand.Command(argsWithDefault.First())
				}
			}
		}
	} else if c.isRoot && cCtx.App.DefaultCommand != "" {
		if dc := cCtx.App.Command(cCtx.App.DefaultCommand); dc != c {
			cmd = dc
		}
	}

	// Descend into the selected subcommand with this command as a stage.
	stages := append(ancestors, commandStage{cmd: c, cCtx: cCtx})
	if cmd != nil {
		newcCtx := NewContext(cCtx.App, nil, cCtx)
		newcCtx.Command = cmd
		return cmd.run(newcCtx, stages, cCtx.Args().Slice())
	}

	// Run the hooks of every stage around this command's Action.
	if c.Action == nil {
		c.Action = helpCommand.Action
	}
	return runStages(stages, func() error {
		err := c.Action(cCtx)
		cCtx.App.handleExitCoder(cCtx, err)
		return err
	})
}

// runStages runs the Before hook and flag actions of each stage from the root
// down, then action. The After hook of each stage reached runs in reverse
// order, even when a later stage or action fails or panics.
func runStages(stages []commandStage, action func() error) (err error) {
	for _, stage := range stages {
		cmd, cCtx := stage.cmd, stage.cCtx

		// Run the After hook on the way out.
		if cmd.After != nil && !cCtx.shellComplete {
			defer func() {
				afterErr := cmd.After(cCtx)
				if afterErr != nil {
					cCtx.App.handleExitCoder(cCtx, err)
					if err != nil {
						err = newMultiError(err, afterErr)
					} else {
						err = afterErr
					}
				}
			}()
		}

		// Run the Before hook, then the flag actions.
		if cmd.Before != nil && !cCtx.shellComplete {
			if beforeErr := cmd.Before(cCtx); beforeErr != nil {
				cCtx.App.handleExitCoder(cCtx, beforeErr)
				return beforeErr
			}
		}
		if err := runFlagActions(cCtx, cmd.Flags); err != nil {
			return err
		}
	}
	return action()
}

// newFlagSet returns an empty flag set for the command's flags.
func (c *Command) newFlagSet() (*flag.FlagSet, error) {
	return flagSet(c.Name, c.Flags, c.separator)
}

// useShortOptionHandling reports whether combined short flags are allowed.
func (c *Command) useShortOptionHandling() bool {
	return c.UseShortOptionHandling
}

// suggestFlagFromError returns a "did you mean" suggestion for the unknown flag
// in err, looking in the named subcommand when command is set. It returns err
// when there is no suggestion.
func (c *Command) suggestFlagFromError(err error, command string) (string, error) {
	// Find the unknown flag and the flags it may have meant.
	flag, parseErr := flagFromError(err)
	if parseErr != nil {
		return "", err
	}
	flags := c.Flags
	hideHelp := c.HideHelp
	if command != "" {
		cmd := c.Command(command)
		if cmd == nil {
			return "", err
		}
		flags = cmd.Flags
		hideHelp = hideHelp || cmd.HideHelp
	}

	// Format the closest match.
	suggestion := SuggestFlag(flags, flag, hideHelp)
	if len(suggestion) == 0 {
		return "", err
	}
	return fmt.Sprintf(SuggestDidYouMeanTemplate, suggestion) + "\n\n", nil
}

// parseFlags parses the command's arguments. A command without subcommands
// accepts flags after its arguments, and a flag defined only by an ancestor
// command sets the ancestor's flag.
func (c *Command) parseFlags(args Args, cCtx *Context) (*flag.FlagSet, error) {
	// Build the flag set; with flag parsing skipped every argument is
	// positional.
	set, err := c.newFlagSet()
	if err != nil {
		return nil, err
	}
	if c.SkipFlagParsing {
		return set, set.Parse(append([]string{"--"}, args.Tail()...))
	}

	// Parse the arguments, letting ancestor flags apply after this command.
	var ancestors []flagScope
	for ctx := cCtx.parentContext; ctx != nil; ctx = ctx.parentContext {
		if ctx.flagSet != nil && ctx.Command != nil {
			ancestors = append(ancestors, flagScope{set: ctx.flagSet, flags: ctx.Command.Flags})
		}
	}
	err = parseArgs(set, c, args.Tail(), !c.hasSubcommands(cCtx.App.helpCmd), ancestors, cCtx.shellComplete)
	if err != nil {
		return nil, err
	}
	if err := normalizeFlags(c.Flags, set); err != nil {
		return nil, err
	}
	return set, nil
}

// hasSubcommands reports whether the command has a subcommand other than
// the app's help command, so its first argument may name a command with its
// own flags.
func (c *Command) hasSubcommands(help *Command) bool {
	for _, sub := range c.Subcommands {
		if sub != help {
			return true
		}
	}
	return false
}

// Names returns the names including short names and aliases.
func (c *Command) Names() []string {
	return append([]string{c.Name}, c.Aliases...)
}

// HasName reports whether name is the command's name or an alias.
func (c *Command) HasName(name string) bool {
	return slices.Contains(c.Names(), name)
}

// VisibleCategories returns the categories that contain a visible command.
func (c *Command) VisibleCategories() []CommandCategory {
	ret := []CommandCategory{}
	for _, category := range c.categories.Categories() {
		if len(category.VisibleCommands()) > 0 {
			ret = append(ret, category)
		}
	}
	return ret
}

// VisibleCommands returns the subcommands that are not hidden.
func (c *Command) VisibleCommands() []*Command {
	var ret []*Command
	for _, command := range c.Subcommands {
		if !command.Hidden {
			ret = append(ret, command)
		}
	}
	return ret
}

// VisibleFlagCategories returns the visible flag categories with their flags.
func (c *Command) VisibleFlagCategories() []VisibleFlagCategory {
	if c.flagCategories == nil {
		c.flagCategories = newFlagCategoriesFromFlags(c.Flags)
	}
	return c.flagCategories.VisibleCategories()
}

// VisibleFlags returns the flags that are not hidden.
func (c *Command) VisibleFlags() []Flag {
	return visibleFlags(c.Flags)
}

// appendFlag adds fl unless the command already has it.
func (c *Command) appendFlag(fl Flag) {
	if !hasFlag(c.Flags, fl) && !hasFlagNamed(c.Flags, fl.Names()[0]) {
		c.Flags = append(c.Flags, fl)
	}
}

// hasCommand reports whether commands contains command.
func hasCommand(commands []*Command, command *Command) bool {
	return slices.Contains(commands, command)
}

// checkDuplicatedCmds returns an error when two subcommands of parent share a
// name or alias.
func checkDuplicatedCmds(parent *Command) error {
	seen := make(map[string]struct{})
	for _, c := range parent.Subcommands {
		for _, name := range c.Names() {
			if _, exists := seen[name]; exists {
				return fmt.Errorf("parent command [%s] has duplicated subcommand name or alias: %s", parent.Name, name)
			}
			seen[name] = struct{}{}
		}
	}
	return nil
}
