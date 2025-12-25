# Agent Instructions

## General Expectations
- Favor interfaces and struct-based implementations over package-level functions.
- Make code easy to test by injecting dependencies (e.g., `io.Reader`, `io.Writer`, or interfaces).
- Add tests for all new functionality and keep behavior covered.

## Platform Requirements
- Target Linux and macOS only.
- Do not implement or document Windows support.
