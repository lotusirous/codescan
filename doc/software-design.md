
# Design

## Background

A static analysis is a function that inspects a package of Go code and reports a set of diagnostics
"checker" runs the analysis and reports the potential vulnerability.

- `analyzer`: describes an analysis function and its options.
- `pass` provides information to the `Run` function that applies a specific analyzer. The `Run` function in Analyzer requires might requires on another analyzer.

## Analyzer

An `Analyzer` statically describes an analysis function: its name, documentation, and its logic.

```go
 package cryptoleak

 var Analyzer = &analysis.Analyzer{
  Name: "cryptoleak",
  Doc:  "check for private_key and public_key in the codebase",
  Run:  run,
  ...
 }
```

## Improvements

### AST

### Dependency analyzer with DAG

### Caching
