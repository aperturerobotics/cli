package cli

import (
	"flag"
	"slices"
	"strings"
)

type iterativeParser interface {
	newFlagSet() (*flag.FlagSet, error)
	useShortOptionHandling() bool
}

// To enable short-option handling (e.g., "-it" vs "-i -t") we have to
// iteratively catch parsing errors. This way we achieve LR parsing without
// transforming any arguments. Otherwise, there is no way we can discriminate
// combined short options from common arguments that should be left untouched.
// Pass `shellComplete` to continue parsing options on failure during shell
// completion when, the user-supplied options may be incomplete.
func parseIter(set *flag.FlagSet, ip iterativeParser, args []string, shellComplete bool) error {
	for {
		err := set.Parse(args)
		if !ip.useShortOptionHandling() || err == nil {
			if shellComplete {
				return nil
			}
			return err
		}

		trimmed, trimErr := flagFromError(err)
		if trimErr != nil {
			return err
		}

		// regenerate the initial args with the split short opts
		argsWereSplit := false
		for i, arg := range args {
			// skip args that are not part of the error message
			if name := strings.TrimLeft(arg, "-"); name != trimmed {
				continue
			}

			// if we can't split, the error was accurate
			shortOpts := splitShortOptions(set, arg)
			if len(shortOpts) == 1 {
				return err
			}

			// swap current argument with the split version
			// do not include args that parsed correctly so far as it would
			// trigger Value.Set() on those args and would result in
			// duplicates for slice type flags
			args = append(shortOpts, args[i+1:]...)
			argsWereSplit = true
			break
		}

		// This should be an impossible to reach code path, but in case the arg
		// splitting failed to happen, this will prevent infinite loops
		if !argsWereSplit {
			return err
		}
	}
}

// flagScope is an ancestor command's parsed flag set and its flag definitions.
type flagScope struct {
	set   *flag.FlagSet
	flags []Flag
}

// parseArgs parses args into set and leaves the positional arguments in
// set.Args().
//
// With interspersed set, flags may come before, between, or after the
// positional arguments, as in "create name --flag", and a "--" terminator
// keeps every later argument positional. Otherwise parsing stops at the first
// positional argument, which names a subcommand with its own flags.
//
// A flag the command does not define is parsed into the nearest of ancestors
// that defines it, so "parent child --parent-flag" sets the parent's flag.
func parseArgs(set *flag.FlagSet, ip iterativeParser, args []string, interspersed bool, ancestors []flagScope, shellComplete bool) error {
	var positional []string
	for {
		if err := parseIter(set, ip, args, shellComplete); err != nil {
			// The parser applied every flag before the undefined one.
			rest, ok, ancestorErr := parseAncestorFlag(err, args, ancestors)
			if ancestorErr != nil {
				return ancestorErr
			}
			if !ok {
				return err
			}
			args = rest
			continue
		}

		// The parser stops at the first positional argument or after "--".
		// The unparsed arguments are always a suffix of args.
		rest := set.Args()
		parsed := len(args) - len(rest)
		if !interspersed || len(rest) == 0 || (parsed != 0 && args[parsed-1] == "--") {
			positional = append(positional, rest...)
			break
		}
		positional = append(positional, rest[0])
		args = rest[1:]
	}
	return set.Parse(append([]string{"--"}, positional...))
}

// parseAncestorFlag parses the flag err reports as undefined, with its value,
// into the nearest ancestor scope that defines it. Returns the arguments after
// the flag, or ok false when err names no flag an ancestor defines.
func parseAncestorFlag(err error, args []string, ancestors []flagScope) ([]string, bool, error) {
	name, nameErr := flagFromError(err)
	if nameErr != nil {
		return nil, false, nil
	}
	i := slices.IndexFunc(args, func(arg string) bool {
		return flagArgName(arg) == name
	})
	if i < 0 {
		return nil, false, nil
	}
	for _, scope := range ancestors {
		f := scope.set.Lookup(name)
		if f == nil {
			continue
		}

		// Parse the flag alone, or with the next argument as its value.
		n := 1
		if !strings.Contains(args[i], "=") && !isBoolFlag(f) && i+1 < len(args) {
			n = 2
		}
		scopeArgs := scope.set.Args()
		if err := scope.set.Parse(args[i : i+n]); err != nil {
			return nil, true, err
		}
		if err := scope.set.Parse(append([]string{"--"}, scopeArgs...)); err != nil {
			return nil, true, err
		}
		if err := normalizeFlags(scope.flags, scope.set); err != nil {
			return nil, true, err
		}
		return args[i+n:], true, nil
	}
	return nil, false, nil
}

// flagArgName returns the flag name an argument such as "--name=value" sets,
// or "" when the argument is not a flag.
func flagArgName(arg string) string {
	if len(arg) < 2 || arg[0] != '-' {
		return ""
	}
	name := strings.TrimPrefix(arg[1:], "-")
	name, _, _ = strings.Cut(name, "=")
	return name
}

// isBoolFlag reports whether the flag takes no value argument.
func isBoolFlag(f *flag.Flag) bool {
	bf, ok := f.Value.(interface{ IsBoolFlag() bool })
	return ok && bf.IsBoolFlag()
}

const providedButNotDefinedErrMsg = "flag provided but not defined: -"

// flagFromError tries to parse a provided flag from an error message. If the
// parsing fials, it returns the input error and an empty string
func flagFromError(err error) (string, error) {
	errStr := err.Error()
	trimmed := strings.TrimPrefix(errStr, providedButNotDefinedErrMsg)
	if errStr == trimmed {
		return "", err
	}
	return trimmed, nil
}

func splitShortOptions(set *flag.FlagSet, arg string) []string {
	shortFlagsExist := func(s string) bool {
		for _, c := range s[1:] {
			if f := set.Lookup(string(c)); f == nil {
				return false
			}
		}
		return true
	}

	if !isSplittable(arg) || !shortFlagsExist(arg) {
		return []string{arg}
	}

	separated := make([]string, 0, len(arg)-1)
	for _, flagChar := range arg[1:] {
		separated = append(separated, "-"+string(flagChar))
	}

	return separated
}

func isSplittable(flagArg string) bool {
	return strings.HasPrefix(flagArg, "-") && !strings.HasPrefix(flagArg, "--") && len(flagArg) > 2
}
