# go-nilcheck

A Go linter that detects functions accepting pointer arguments but failing to check for `nil` values. This helps enforce best practices in Go programs, ensuring safer and more robust code.

## Table of Contents
- [go-nilcheck](#go-nilcheck)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Features](#features)
  - [Installation](#installation)
  - [Usage](#usage)
  - [How It Works](#how-it-works)
    - [Analyzer Logic](#analyzer-logic)
    - [Implementation Details](#implementation-details)
  - [Examples](#examples)
    - [Example 1: Function Missing nil Check](#example-1-function-missing-nil-check)
    - [Example 2: Correct Usage](#example-2-correct-usage)
  - [Contributing](#contributing)
  - [License](#license)

## Introduction

The Nil Pointer Linter analyzes Go code to identify functions that:
1. Accept pointer arguments.
2. Do not verify if the pointers are `nil` before using them.

This ensures that potential runtime `nil` dereference errors are avoided, promoting cleaner and safer Go code.

## Features

- Detects functions with pointer arguments.
- Analyzes the function body for `nil` checks on pointer arguments.
- Reports functions missing `nil` checks with specific line numbers and function names.

## Installation

To install the linter, use `go install`:

```bash
go install github.com/jtbonhomme/go-nilcheck@latest
```

## Usage
Run the linter on your Go project using:

```sh
nilcheck ./...
```

This will analyze all .go files in the specified directory (and subdirectories) and report any violations.

## How It Works

### Analyzer Logic

1. **Identify Functions**: The linter inspects the abstract syntax tree (AST) of Go files to identify all function declarations.
1. **Extract Pointer Arguments**: It checks the function's parameters for pointer arguments.
1. **Analyze Function Body**: It searches the function body to ensure nil checks are performed for all pointer arguments.
1. **Report Issues**: If any pointer argument is not checked for nil, the linter reports an issue with the function's name and position.

### Implementation Details

The main logic resides in the `analyzer.go` file:

* `getPointerArguments`: Extracts all pointer arguments from a function's parameters.
* `checkNilInBody`: Ensures the function body contains appropriate nil checks.
* `hasNilCheck`: Recursively inspects expressions to detect nil checks.

The `main.go` file integrates the linter with the `golang.org/x/tools/go/analysis/singlechecker` package, enabling command-line usage.

## Examples

### Example 1: Function Missing nil Check

```go
func ExampleFunc(p *int) {
    *p = 42 // Potential runtime panic if `p` is nil
}
```

**Output:**

```go
main.go:5: Function ExampleFunc has pointer arguments but does not check for nil
```

### Example 2: Correct Usage

```go
func ExampleFunc(p *int) {
    if p != nil {
        *p = 42
    }
}
```

**No issues are reported.**

## Contributing

Contributions are welcome! To contribute:

1. Fork the repository.
1. Create a feature branch.
1. Commit your changes.
1. Submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
