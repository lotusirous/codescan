
# Scanner design

A static analysis is a function that inspects a package of codebase and reports a set of diagnostics
"checker" runs the analysis and reports the potential vulnerability. For example the `cryptoleak` rule reports the potential cryptography keys leak in the codebase.

**An `Analyzer`** statically describes an analysis function: its name, description, severity, and of course, its logic.

To define an analysis, a user declares a (logically constant) variable of type Analyzer. Here is a typical example from one of the analyzers in the checker/rules subdirectory:

```go
package cryptoleak

var Analyzer = &analysis.Analyzer{
 Name: "cryptoleak",
 Meta: analysis.Meta{
  Description: "Leak the cryptography keys",
  Severity:    "HIGH",
 },
 Run: run,
}
```

An rule set is defined by a named rule and analysis function.

```go
var allRules = map[string]*analysis.Analyzer{
 "G402": cryptoleak.Analyzer,
}
```

**A Pass** represents a single unit of work in the analysis process. It provides information to the Analyzer's Run function regarding the package being analyzed. Additionally, it also holds the parsed Abstract Syntax Tree (AST) of a codebase, which allows analyzers to extract and analyze the code further.

**The Runner** handles the execution of the set of rules. It also provides parallel execution and resolves dependencies for the Analyzer, if any.

**The Scanner**: Since code scanning involves two types of analysis: static and dynamic.  organizes the rules by type and performs the analysis according to the designated rule set.

## Improvements

**AST**: This project focuses on simple code scanning and does not include AST analysis. This is because AST analysis requires parsing a wide range of programming languages, such as JavaScript, Python, Java, etc.

**Dependency Analyzer with DAG**: The analyzers can depend on other analyzers and may require the results from parent analyzers. To handle this, the execution utilizes a Directed Acyclic Graph (DAG) to resolve and run the analyzers. The Runner need to be designed to support both linear and parallel execution
