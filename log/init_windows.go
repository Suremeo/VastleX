package log

import (
	"golang.org/x/sys/windows"
	"os"
)

// init is a hack to make sure colors will work on windows :).
func init() {
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32
	_ = windows.GetConsoleMode(stdout, &originalMode)
	_ = windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}
