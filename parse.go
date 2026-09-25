package cli

import (
	"flag"
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

// parseInterspersed parses flags placed before, between, or after the
// positional arguments, as in "create name --flag". A "--" terminator keeps
// every later argument positional. The positional arguments become set.Args().
func parseInterspersed(set *flag.FlagSet, ip iterativeParser, args []string, shellComplete bool) error {
	var positional []string
	for {
		if err := parseIter(set, ip, args, shellComplete); err != nil {
			return err
		}

		// The parser stops at the first positional argument or after "--".
		// The unparsed arguments are always a suffix of args.
		rest := set.Args()
		parsed := len(args) - len(rest)
		if len(rest) == 0 || (parsed != 0 && args[parsed-1] == "--") {
			positional = append(positional, rest...)
			break
		}
		positional = append(positional, rest[0])
		args = rest[1:]
	}
	return set.Parse(append([]string{"--"}, positional...))
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
