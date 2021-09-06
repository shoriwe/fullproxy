package Templates

import (
	"fmt"
	"os"
)

func Exit(message string, wantedArguments int, receivedArguments int, args ...interface{}) {
	if receivedArguments < wantedArguments {
		_, _ = fmt.Fprintf(os.Stderr, message, args...)
		os.Exit(-1)
	}
}
