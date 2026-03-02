package stackerr_test

import (
	"errors"
	"fmt"
	"os"

	"github.com/LZStock-OS/stackerr"
)

// ExampleRecover demonstrates how to use Recover to catch panics and wrap them with a stack trace.
func ExampleRecover() {
	// Disable default logging to stderr for this example to keep output clean
	stackerr.Config.Output = nil
	defer func() { stackerr.Config.Output = os.Stderr }() // Restore default

	// Function that might panic
	riskyFunction := func() (err error) {
		// Defer Recover to catch panics and populate err
		defer stackerr.Recover(&err)

		// Simulating a panic
		panic("something went wrong")
	}

	err := riskyFunction()
	if err != nil {
		fmt.Printf("Error caught: %v\n", err)
		// Access the stack trace if needed
		if stack := stackerr.GetStack(err); stack != "" {
			fmt.Println("Stack trace captured.")
		}
	}

	// Output:
	// Error caught: panic: something went wrong
	// Stack trace captured.
}

// ExampleThrowPanic demonstrates how to use ThrowPanic to panic with a stack trace.
func ExampleThrowPanic() {
	// Disable default logging
	stackerr.Config.Output = nil
	defer func() { stackerr.Config.Output = os.Stderr }()

	// Function that throws a panic
	thrower := func() (err error) {
		defer stackerr.Recover(&err)
		
		// ThrowPanic wraps the error and panics
		stackerr.ThrowPanic(errors.New("critical failure"))
		return nil
	}

	err := thrower()
	if err != nil {
		fmt.Printf("Error caught: %v\n", err)
	}

	// Output:
	// Error caught: critical failure
}

// Example_configuration demonstrates how to configure the library.
func Example_configuration() {
	// Configure global settings
	stackerr.Config.StripSequences = []string{"(0", "({", "(*"}
	stackerr.Config.MaxStackDepth = 5
	// Disable logging for this example to avoid output mismatch
	stackerr.Config.Output = nil

	var err error
	func() {
		defer stackerr.Recover(&err)
		panic("custom config panic")
	}()

	// Reset config for other tests
	stackerr.Config.Output = os.Stderr
	
	fmt.Println("Done")
	
	// Output:
	// Done
}
