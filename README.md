# stackerr

[![Go Report Card](https://goreportcard.com/badge/github.com/LZStock-OS/stackerr)](https://goreportcard.com/report/github.com/LZStock-OS/stackerr)
[![Go Doc](https://pkg.go.dev/badge/github.com/LZStock-OS/stackerr.svg)](https://pkg.go.dev/github.com/LZStock-OS/stackerr)
[![License](https://img.shields.io/github/license/LZStock-OS/stackerr)](LICENSE)
[![Release](https://img.shields.io/github/v/release/LZStock-OS/stackerr)](https://github.com/LZStock-OS/stackerr/releases)

`stackerr` is a lightweight Go library for panic recovery and error handling with customizable stack traces. It simplifies debugging by capturing clean, readable stack traces when panics occur or when errors are explicitly thrown.

## Features

- **Panic Recovery**: Easily recover from panics and wrap them into errors with stack traces.
- **Clean Stack Traces**: Automatically strips noisy arguments and simplifies file paths (e.g., replacing `$GOPATH` with `[Proj]`).
- **Customizable**: Configure which strings to strip from function names, define custom path replacements, and set stack depth limits.
- **Error Wrapping**: Supports standard Go error wrapping (`Unwrap()`).

## Installation

```bash
go get github.com/LZStock-OS/stackerr
```

## Usage

### Basic Recovery

Use `stackerr.Recover` in a `defer` statement to catch panics and populate an error return variable.

```go
package main

import (
	"fmt"
	"github.com/LZStock-OS/stackerr"
)

func riskyOperation() (err error) {
	defer stackerr.Recover(&err)

	// Simulating a panic
	panic("something went wrong")
}

func main() {
	if err := riskyOperation(); err != nil {
		fmt.Printf("Caught error: %v\n", err)
		
		if stack := stackerr.GetStack(err); stack != "" {
			fmt.Println("Stack Trace:")
			fmt.Println(stack)
		}
	}
}
```

### Throwing Errors with Stack Traces

Use `stackerr.ThrowPanic` to explicitly panic with a wrapped error containing the current stack trace.

```go
func doSomething() {
    err := someInternalFunction()
    if err != nil {
        // This will panic and be caught by a defer stackerr.Recover up the chain
        stackerr.ThrowPanic(fmt.Errorf("operation failed: %w", err))
    }
}
```

### Configuration

You can customize the library's behavior by modifying the global `stackerr.Config` variable. It is recommended to do this in your `init()` function or main setup.

```go
import (
    "os"
    "github.com/LZStock-OS/stackerr"
)

func init() {
    // Customize string sequences to strip from function names
    stackerr.Config.StripSequences = []string{"(0", "({", "(*"}
    
    // Set maximum stack depth (default: 10)
    stackerr.Config.MaxStackDepth = 20
    
    // Customize path replacements to hide sensitive paths or make logs shorter
    stackerr.Config.PathReplacements = map[string]string{
        "/home/user/go/src/": "[Src]",
        "github.com/my/project/": "", 
    }
    
    // Redirect log output (default: os.Stderr)
    // Set to nil to disable automatic logging on recovery
    stackerr.Config.Output = os.Stdout 
}
```

## API Reference

### `func Recover(err *error)`

Catches panics, wraps them in a `callStackErr`, and assigns it to `*err`. It also logs the panic and stack trace to the configured output.

### `func ThrowPanic(err error)`

Wraps the given error with the current stack trace and panics. Use this when you want to abort execution but preserve context for the recovery handler.

### `func GetStack(err error) string`
Retrieves the stack trace from an error if it was created by `Recover` or `ThrowPanic`. Returns an empty string otherwise.

## License

MIT