
# Scanner design

## Background

A static analysis is a function that inspects a project of code and reports a set of diagnostics

- Analyzer: describes an analysis function and its options.
- Pass: A Pass provides information to the Run function that applies a specific analyzer. The `Run` function in Analyzer requires on this struct. If we
One pass may depend on the result computed by another.

## Analyzer

The primary type in the API is Analyzer. An Analyzer statically describes an analysis function: its name, documentation, flags, relationship to other analyzers, and of course, its logic.

To define an analysis, a user declares a (logically constant) variable of type Analyzer

# Caching

The runner caches facts, directives and diagnostics in a
content-addressable cache that is designed after Go's own cache.
Additionally, it makes use of Go's export data.

This cache not only speeds up repeat runs, it also reduces peak
memory usage. When we've analyzed a project/package, we cache the results
and drop them from memory. When a dependent needs any of this
information, or when analysis is complete and we wish to render the
results, the data gets loaded from disk again.

## Improvements

## Parallelism

- If the analyzer depends on other analyzer, the runner can build a graph structure to resolve the analyzer dependencies.
- Actions are executed in parallel where the dependency graph allows. Overall parallelism is bounded by a semaphore, sized according to GOMAXPROCS. Each concurrently processed project/package takes up a
token, as does each analyzer - but a project/project/package can always execute at
least one analyzer, using the project/package's token.
Depending on the overall shape of the graph, there may be GOMAXPROCS
packages running a single analyzer each, a single package running
GOMAXPROCS analyzers, or anything in between.

Total memory consumption grows roughly linearly with the number of
CPUs, while total execution time is inversely proportional to the
number of CPUs. Overall, parallelism is affected by the shape of
the dependency graph. A lot of inter-connected packages will see
less parallelism than a lot of independent packages.

## Caching

The runner can use a tree-based structure for analyzing, which can be helpful. This feature allows the runner to support the AST for parsed language.

This cache feature not only speeds up repeat runs, it also reduces peak memory usage. When the runner has analyzed a package, it caches the results and drops them from memory. When a dependent package needs any of this information, or when analysis is complete and the results need to be rendered, the data gets loaded from disk again.

By keeping data only in memory when it is immediately needed and not retained for possible future uses, the runner can reduce memory usage. This approach trades increased CPU usage for reduced memory usage. A single dependency may be loaded many times over, but it greatly reduces peak memory usage. An arbitrary amount of time may pass between analyzing a dependency and its dependent, during which other packages will be processed. This helps to keep the peak memory usage low.
