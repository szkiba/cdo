# How to contribute to cdo

Thank you for your interest in contributing to cdo!

Before you begin, make sure to familiarize yourself with the [Code of Conduct](CODE_OF_CONDUCT.md). If you've previously contributed to other open source project, you may recognize it as the classic [Contributor Covenant](https://contributor-covenant.org/).

## Filing issues

Don't be afraid to file issues! Nobody can fix a bug we don't know exists, or add a feature we didn't think of.

1. **Ensure the bug was not already reported** by searching on GitHub under [Issues](https://github.com/szkiba/cdo/issues).

2. If you're unable to find an open issue addressing the problem, [open a new one](https://github.com/szkiba/cdo/issues/new). Be sure to include a **title and clear description**, as much relevant information as possible.


The worst that can happen is that someone closes it and points you in the right direction.

## Contributing code

If you'd like to contribute code to cdo, this is the basic procedure.

1. Find an issue you'd like to fix. If there is none already, or you'd like to add a feature, please open one, and we can talk about how to do it. Out of respect for your time, please start a discussion regarding any bigger contributions in a GitHub Issue **before** you get started on the implementation.

2. Create a fork and open a feature branch - `feature/my-cool-feature` is the classic way to name these, but it really doesn't matter.

3. Create a pull request!

4. We will discuss implementation details until everyone is happy, then a maintainer will merge it.

## Typical tasks

This section describes the typical tasks of contributing code.

If the [cdo](https://github.com/szkiba/cdo) tool is installed, the tasks can be easily executed.

<!-- #region lint -->
### lint - Run the linter

The `golangci-lint` tool is used for static analysis of the source code.
It is advisable to run it before committing the changes.

```bash
golangci-lint run
```
<!-- #endregion lint -->

[lint]: <#lint---run-the-linter>

### readme - Update the README.md

In order to keep README.md up to date, some parts of it are updated from other files. For example, the task definition examples are updated from the `CONTRIBUTING.md` file using the [mdcode] tool.

```bash
mdcode update
```

[mdcode]: <https://github.com/szkiba/mdcode>

### test - Run the tests

```bash
go test -count 1 -race -coverprofile=coverage.txt ./...
```

[test]: <#test---run-the-tests>

### coverage - View the test coverage report

Requires
: [test]

```bash
go tool cover -html=coverage.txt
```

### build - Build the executable binary

This is the easiest way to create an executable binary (although the release process uses the `goreleaser` tool to create release versions).

```bash
go build -ldflags="-w -s" -o build/cdo .
```

[build]: <#build---build-the-executable-binary>

### snapshot - Creating an executable binary with a snapshot version

The goreleaser command-line tool is used during the release process. During development, it is advisable to create binaries with the same tool from time to time.

```bash
rm -f build/cdo
goreleaser build --snapshot --clean --single-target -o build/cdo
```

[snapshot]: <#snapshot---creating-an-executable-binary-with-a-snapshot-version>

### clean - Delete the build directory

```bash
rm -rf build
```

<!-- #region ci -->

### ci - Run all ci-relevant tasks

Run all the tasks that will run in the Continuous Integration system.

Requires
: [lint], [test], [build], [snapshot]

<!-- #endregion ci -->

### makefile - Generate Makefile

Generate Makefile from CONTRIBURIG.md task definitions.

```bash
cdo --makefile Makefile
```