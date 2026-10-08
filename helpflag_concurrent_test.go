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
