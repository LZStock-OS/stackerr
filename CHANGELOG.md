# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-03-01

### Added
- Initial release of `stackerr` library.
- `Recover` function to handle panics and wrap errors with stack traces.
- `ThrowPanic` function to explicitly panic with a wrapped error containing the current stack trace.
- Internal `callStackErr` type implementing `error` and `Unwrap` interface.
- Configuration support via `stackerr.Config`:
  - `StripSequences`: Custom strings to strip from function names (e.g., `(0`, `({`).
  - `PathReplacements`: Custom file path cleanup (e.g., replacing `$GOPATH` with `[Proj]`).
  - `MaxStackDepth`: Configurable stack trace depth limit (default: 10).
  - `Output`: Configurable `io.Writer` for panic logs (default: `os.Stderr`).
- Unit tests (`stackerr_test.go`) and usage examples (`example_test.go`).
- GitHub Actions CI workflow for automated testing.
- MIT License.
