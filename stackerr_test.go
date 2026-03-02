package stackerr_test

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/LZStock-OS/stackerr"
)

// Helper to reset config after each test
func resetConfig() {
	stackerr.Config.StripSequences = []string{"(0", "({"}
	stackerr.Config.PathReplacements = stackerr.GetDefaultPathReplacements()
	stackerr.Config.MaxStackDepth = 10
	stackerr.Config.Output = os.Stderr
}

func TestRecover(t *testing.T) {
	defer resetConfig()
	var err error

	func() {
		defer stackerr.Recover(&err)
		panic("boom")
	}()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	stack := stackerr.GetStack(err)
	if stack == "" {
		t.Fatal("expected stack trace")
	}
	if err.Error() != "panic: boom" {
		t.Errorf("expected panic: boom, got %v", err)
	}
	if !strings.Contains(stack, "TestRecover") {
		t.Errorf("stack trace should contain TestRecover")
	}
}

func TestThrowPanic(t *testing.T) {
	defer resetConfig()
	var err error

	func() {
		defer stackerr.Recover(&err)
		stackerr.ThrowPanic(errors.New("throw error"))
	}()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	stack := stackerr.GetStack(err)
	if stack == "" {
		t.Fatal("expected stack trace")
	}
	if err.Error() != "throw error" {
		t.Errorf("expected throw error, got %v", err)
	}
	if !strings.Contains(stack, "TestThrowPanic") {
		t.Errorf("stack trace should contain TestThrowPanic")
	}
}

func TestMaxStackDepth(t *testing.T) {
	defer resetConfig()
	stackerr.Config.MaxStackDepth = 2

	var err error
	func() {
		defer stackerr.Recover(&err)
		panic("depth test")
	}()

	stack := stackerr.GetStack(err)
	if stack == "" {
		t.Fatal("expected stack trace")
	}

	lines := strings.Split(strings.TrimSpace(stack), "\n")
	if len(lines) > 2 {
		t.Errorf("expected at most 2 stack frames, got %d\nStack:\n%s", len(lines), stack)
	}
}

func TestPathReplacement(t *testing.T) {
	defer resetConfig()

	// We'll replace the path to the test file with something unique
	// Since we don't know the absolute path easily in test without runtime.Caller,
	// let's just use a replacement that likely matches part of the path.
	// We know the file is stackerr_test.go.
	// But the library does simple string replacement on the file path.

	// Let's assume the stack trace contains "stackerr_test.go".
	// We want to replace "stackerr" with "replaced_lib".
	stackerr.Config.PathReplacements = map[string]string{
		"stackerr": "replaced_lib",
	}

	var err error
	func() {
		defer stackerr.Recover(&err)
		panic("path test")
	}()

	stack := stackerr.GetStack(err)
	if !strings.Contains(stack, "replaced_lib") {
		// It's possible the file path in stack trace is absolute and doesn't contain "stackerr" if run in a weird way,
		// but usually it should be there.
		// However, in the sandbox environment, path is /Users/universetennis/Code/go/src/stackerr/stackerr_test.go
		// So "stackerr" is definitely in there.
		t.Errorf("Stack trace expected to contain 'replaced_lib', got:\n%s", stack)
	}
}

func TestCustomOutput(t *testing.T) {
	defer resetConfig()
	var buf bytes.Buffer
	stackerr.Config.Output = &buf

	var err error
	func() {
		defer stackerr.Recover(&err)
		panic("log test")
	}()

	output := buf.String()
	if !strings.Contains(output, "Recovered from panic: panic: log test") {
		t.Errorf("expected log output to contain panic message, got: %s", output)
	}
}

func TestRecoverStandardError(t *testing.T) {
	defer resetConfig()
	var err error = errors.New("standard error")

	func() {
		defer stackerr.Recover(&err)
	}()

	stack := stackerr.GetStack(err)
	if stack == "" {
		t.Fatal("expected stack trace")
	}
	if err.Error() != "standard error" {
		t.Errorf("expected error message 'standard error', got '%s'", err.Error())
	}
}

func TestRecoverNilError(t *testing.T) {
	defer resetConfig()
	var err error

	func() {
		defer stackerr.Recover(&err)
	}()

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}
