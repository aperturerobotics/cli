package cli

import (
	"io"
	"sync"
	"testing"
)

// TestConcurrentAppsWithHelpFlag runs two Apps in parallel. Each app installs
// and parses its own copy of the built-in help flag; the copies must not share
// the package-level HelpFlag, whose fields every Apply call mutates.
func TestConcurrentAppsWithHelpFlag(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			app := &App{
				Name:      "test",
				Writer:    io.Discard,
				ErrWriter: io.Discard,
			}
			if err := app.Run([]string{"test", "--help"}); err != nil {
				t.Errorf("Run returned unexpected error: %v", err)
			}
		}()
	}
	wg.Wait()
}

// TestConcurrentAppsWithSubcommands runs Apps with nested subcommands in
// parallel. Setup installs the built-in help command in each command and
// writes its help name, so the Apps must not share one instance.
func TestConcurrentAppsWithSubcommands(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			app := &App{
				Name:      "test",
				Writer:    io.Discard,
				ErrWriter: io.Discard,
				Commands: []*Command{{
					Name: "sub",
					Subcommands: []*Command{{
						Name:   "leaf",
						Action: func(*Context) error { return nil },
					}},
				}},
			}
			for _, args := range [][]string{
				{"test", "sub", "leaf"},
				{"test", "help", "sub"},
				{"test", "sub", "help"},
				{"test", "sub", "--help"},
			} {
				if err := app.Run(args); err != nil {
					t.Errorf("Run(%v) returned unexpected error: %v", args, err)
				}
			}
		}()
	}
	wg.Wait()
}
