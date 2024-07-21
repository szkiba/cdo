# Example: dependency

This example shows how to specify dependencies between tasks.

Let's say the test run generates the coverage data, and the coverage tool creates a coverage report from the coverage data.

## prepare - Do some preparation

```sh
echo "task: prepare"
```

## test - Run tests

Run tests and collect coverage data.

```sh
echo "task: test"
```

## coverage - Report test coverage

Create a test coverage report from previously collected coverage data.

Requires
: prepare, test

```sh
echo "task: coverage"
```

## Execute

Since the `coverage` task specified the `prepare` and `test` tasks as dependencies with the **Requires** directive, the `prepare` and `test` tasks will be executed before the `coverage` task.

```
cdo coverage
```

Output:

```text
task: prepare
task: test
task: coverage
```