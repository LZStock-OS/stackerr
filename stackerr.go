package stackerr

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
)

// Version is the current version of the library.
const Version = "v0.1.0"

// Configuration holds the configuration for the stackerr library.
type Configuration struct {
	// StripSequences is a list of strings that, if found in a function name,
	// will cause the function name to be truncated from that point.
	// Default: []string{"(0", "({"}
	StripSequences []string

	// PathReplacements is a map of path prefixes to replace in the file path.
	// Default: derived from GOPATH and GOROOT.
	PathReplacements map[string]string

	// MaxStackDepth is the maximum number of stack frames to include.
	// Default: 10
	MaxStackDepth int

	// Output is the writer where panic recovery logs are written.
	// If nil, no logging is done. Default: os.Stderr
	Output io.Writer
}

// defaultStripSequences defines the default sequences to strip from function names.
var defaultStripSequences = []string{"(0", "({"}

// Config is the global configuration for the library.
// You can modify this variable directly to configure the library.
var Config = Configuration{
	StripSequences:   defaultStripSequences,
	PathReplacements: GetDefaultPathReplacements(),
	MaxStackDepth:    10,
	Output:           os.Stderr,
}

// callStackErr wraps an error with a stack trace.
type callStackErr struct {
	Err   error
	Stack string
}

func (e callStackErr) Error() string {
	return e.Err.Error()
}

func (e callStackErr) GetStack() string {
	return e.Stack
}

func (e callStackErr) Unwrap() error {
	return e.Err
}

// GetStack retrieves the stack trace from an error if it is a callStackErr.
// It returns an empty string if the error does not contain a stack trace.
func GetStack(err error) string {
	if err == nil {
		return ""
	}
	if cse, ok := err.(callStackErr); ok {
		return cse.Stack
	}
	return ""
}

// formatStackTrace parses and cleans up the raw stack trace from debug.Stack().
func formatStackTrace(stack string) string {
	lines := strings.Split(stack, "\n")
	var frames []string

	// Skip the first few lines which are usually the debug.Stack() call itself
	for i := 1; i < len(lines)-1; i += 2 {
		funcLine := strings.TrimSpace(lines[i])
		fileLine := strings.TrimSpace(lines[i+1])

		// Clean up function arguments to keep the output readable
		for _, seq := range Config.StripSequences {
			if idx := strings.Index(funcLine, seq); idx != -1 {
				funcLine = funcLine[:idx]
				break
			}
		}

		parts := strings.SplitN(fileLine, ":", 2)
		if len(parts) != 2 {
			continue
		}

		filePath := strings.TrimSpace(parts[0])
		lineInfo := strings.TrimSpace(parts[1])
		lineParts := strings.Split(lineInfo, " ")
		lineNumber := lineParts[0]

		// Apply path replacements
		for oldPath, newPath := range Config.PathReplacements {
			if oldPath != "" && oldPath != newPath {
				filePath = strings.Replace(filePath, oldPath, newPath, 1)
			}
		}

		frame := fmt.Sprintf("File: %s:%s Func: %s", filePath, lineNumber, funcLine)
		frames = append(frames, frame)
	}

	// Limit stack trace depth to avoid overly huge logs
	depth := Config.MaxStackDepth
	if depth <= 0 {
		depth = 10 // Fallback safety
	}

	start := len(frames) - depth
	if start < 0 {
		start = 0
	}

	var buf bytes.Buffer
	for i, frame := range frames[start:] {
		buf.WriteString(frame)
		if i < len(frames[start:])-1 {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

// Recover should be used in defer to handle panics and wrap errors with stack traces.
// Usage: defer stackerr.Recover(&err)
func Recover(err *error) {
	// 1. Handle Panic
	if r := recover(); r != nil {
		var panicError error
		var formattedStack string

		// Check if it's already a callStackErr (from ThrowPanic)
		if cse, ok := r.(callStackErr); ok {
			panicError = cse.Err
			formattedStack = cse.Stack
		} else {
			switch v := r.(type) {
			case error:
				panicError = v
			case string:
				panicError = fmt.Errorf("panic: %s", v)
			default:
				panicError = fmt.Errorf("panic: %v", v)
			}

			stack := string(debug.Stack())
			formattedStack = formatStackTrace(stack)
		}

		// Assign the panic error to the return value
		*err = callStackErr{
			Err:   panicError,
			Stack: formattedStack,
		}

		// Logging
		if Config.Output != nil {
			fmt.Fprintf(Config.Output, "Recovered from panic: %v\nStack: %s\n", panicError, formattedStack)
		}
		return
	}

	// 2. Handle standard errors (wrap with stack if not already wrapped)
	if err != nil && *err != nil {
		if _, ok := (*err).(callStackErr); !ok {
			stack := string(debug.Stack())
			formattedStack := formatStackTrace(stack)

			*err = callStackErr{
				Err:   *err,
				Stack: formattedStack,
			}
		}
	}
}

// ThrowPanic wraps an error with a stack trace and panics.
func ThrowPanic(err error) {
	if err != nil {
		message := err.Error()

		stack := string(debug.Stack())
		formattedStack := formatStackTrace(stack)

		// Create the error first
		wrappedErr := callStackErr{
			Err:   fmt.Errorf("%s", message),
			Stack: formattedStack,
		}

		// Log and Panic
		// If logging is desired before panic, we could do it here, but usually panic bubbles up to Recover.
		// However, ThrowPanic explicitly panics with the wrapper.
		panic(wrappedErr)
	}
}

// GetDefaultPathReplacements initializes the default path replacements.
func GetDefaultPathReplacements() map[string]string {
	replacements := map[string]string{}
	if goroot := os.Getenv("GOROOT"); goroot != "" {
		replacements[goroot] = "[GOROOT]"
	}
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		replacements[gopath+"/src/"] = "[Proj]"
		replacements[gopath+"/pkg/"] = "[Lib]"
	}
	return replacements
}
