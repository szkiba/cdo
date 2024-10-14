# ｃｄｏ

**Markdown-based task runner for contributors**

The `cdo` (/siː duː/) command line tool allows contributors to perform routine tasks conveniently. The markdown code blocks defining the tasks can be specified in the `CONTRIBUTING.md` or `README.md` files by default. Tasks are executed using a portable, embedded `bash`-like shell.

The goal is to help contributors in performing routine tasks in a uniform, portable way. It is not intended to replace complex, programmable task runners (grunt, gulp, etc.).

Here is a `lint` task definition from cdo's own [`CONTRIBUTING.md`](CONTRIBUTING.md):

<!-- Do not edit the following code block, it is updated from CONTRIBUTING.md -->
~~~markdown file=CONTRIBUTING.md region=lint
### lint - Run the linter

The `golangci-lint` tool is used for static analysis of the source code.
It is advisable to run it before committing the changes.

```bash
golangci-lint run
```
~~~

To execute the `lint` task above:

```bash
cdo lint
```

List of available tasks with a short descriptions:

```bash
cdo
```

Display the long description of a specific task:

```bash
cdo lint --help
```

### Features

- the contributing documentation becomes task definitions
- the existing `CONTRIBUTING.md`/`README.md` is used for task definition
- human-readable alternative to make/Makefile (for simple tasks)
- dependencies can be specified for tasks
- portable, bash-like embedded shell
- [BusyBox support](#busybox) for non-built-in commands for portability
- the tasks can be executed even without the `cdo`
- Makefile can be generated from the task definitions

## Install

Precompiled binaries can be downloaded and installed from the [Releases](https://github.com/szkiba/cdo/releases) page.

If you have a go development environment (not required to use `cdo`), the installation can also be done with the following command:

```
go install github.com/szkiba/cdo@latest
```

## Examples

The [examples](examples) directory contains examples of how to use cdo. Each example is a subdirectory in which the CONTRIBUTING.md file contains the task definitions.

As an example, cdo's own task definitions can also be used, which can be found in the [`CONTRIBUTING.md`](CONTRIBUTING.md) file in the [Typical tasks](CONTRIBUTING.md#typical-tasks) section.


## Defining tasks

### Task definition file

`cdo` looks for task definitions in the `CONTRIBUTING.md` and `README.md` files by default.

First, the `CONTRIBUTING.md` file is searched recursively from the current directory to the root directory. If `CONTRIBUTING.md` exists in the given directory or its `docs/` subdirectory, the search will stop and this file will be used as the task definition file.

If `CONTRIBUTING.md` is not found, the `README.md` file is searched recursively from the current directory to the root directory. If `README.md` exists in the given directory, the search will stop and this file will be used as the task definition file.

Any other markdown file can be used for task definitions using the `-f/--file` flag. No search will be performed, the exact location of the task definition file must be specified.

### Tasks

The structure of the task definition file is relatively loose, basically determined by the content of the contribution documentation. For example, it is not necessary to put the task definitions under a special section. There can be task definitions both under **Submit an issue** and **Contribute code** sections (or under any other section).

#### Name and short description

The task definition must include a heading element with an appropriate format. The heading level doesn't matter. The heading element must contain the ` - ` (space, hyphen, space) separator character sequence. The separator character sequence divides the heading into two parts: the first part is the name of the task, and the second part is a short description of the task.

The task definition optionally contains a code block with the language `bash` (or `sh`).

For example:

~~~markdown
### readme - Update the README.md

```bash
./tools/update-readme
```
~~~

The name of the task will be `readme`, the short description will be `Update the README.md` and the `./tools/update-readme` program will run when the task is executed.

#### Commands

The code block containing the task definition is executed as a shell script with an embedded bash-like shell. You can use the usual bash control statements (`if`, `for`) and variable substitutions. Since the script is executed by an embedded shell, it will work the same way on all operating systems. Of course, the external commands used in the script (`grep`, `find`, `curl`) must be available, otherwise an execution error will occur.

Check [examples/shell](examples/shell/CONTRIBUTING.md) for more information on advanced shell features.

#### Help

The long description of the task can be displayed using the `-h/--help` flag. The long description of the task consists of the heading element and the markdown text between the heading element and the code block.

For example:

~~~markdown
### readme - Update the README.md

In order to keep README.md up to date, some parts of it are updated from other files.
For example, the task definition examples are updated from the `CONTRIBUTING.md` file using the [mdcode] tool.

```bash
./tools/update-readme
```
~~~

The long description will be as follows:

~~~
readme - Update the README.md

In order to keep README.md up to date, some parts of it are updated from other files.
For example, the task definition examples are updated from the `CONTRIBUTING.md` file using the [mdcode] tool.
~~~

### Variables

Shell variable substitutions can be used in tasks. In addition to simple variable substitutions, more sophisticated forms can be used, such as:

```bash
echo ${FOO:-bar}
# if FOO has no value, "bar" will be printed
```

Positional arguments are accessed in the usual way:

```bash
echo $*
# all positional arguments will be printed
echo $1
# the first positional argument will be printed
```

The value of the variables can be set in the task definition itself or by using the `-e/--env` flag or in dotenv files (`.env`, `.env.local`).

#### Named parameters

It is possible to assign values ​​to variables in any part of the command line by simple assignment. This can also be considered as passing parameters by name to the task.

```bash
cdo build os=linux cpu=amd64
# The os and cpu variables will be passed to the build task
```

Check [examples/named](examples/named/CONTRIBUTING.md) for more information.

#### Dotenv files

If a file called `.env` and/or `.env.local` exists in the directory containing the task definition file, it will be read and variables defined in it will be available in every task. If a variable is assigned a value in both the `.env` and `.env.local` files, the value assigned in `.env.local` will be used. Since `.env.local` is conveniently included in `.gitignore`, it can be used for local settings.

Lines in the dotenv file contain variable value assignments separated by an equal sign (`=`). The `export` keyword can optionally be used at the beginning of the line. The hashmark (`#`) character can be used for comments.

```sh
# I am a comment and that is OK
SOME_VAR=someval
FOO=BAR # comments at line end are OK too
export BAR=BAZ
```

Check [examples/dotenv](examples/dotenv/CONTRIBUTING.md) for more information.

### Dependencies

Tasks can have one or more other tasks as dependencies. The execution of the dependencies precedes the execution of the task.

Dependencies can be specified using a markdown definition list. A definition term with the name Requires must be created. In the subsequent definition description, you must enter the names of the dependency tasks, separated by commas.

The example below specifies two dependencies: test and build:

```markdown
Requires
: test, build
```

Unfortunately, GitHub does not support the markdown definition list, so the definition term and definition description are rendered in one line separated by a colon:

Requires
: test, build

Among cdo's own tasks, the `ci` task is a good example of specifying dependencies:

~~~markdown file=CONTRIBUTING.md region=ci

### ci - Run all ci-relevant tasks

Run all the tasks that will run in the Continuous Integration system.

Requires
: [lint], [test], [build], [snapshot]

~~~

In this example, it can also be seen that the dependency can also be a markdown link to the corresponding task.

Check [examples/dependency](examples/dependency/CONTRIBUTING.md) for more information on dependency support.

### BusyBox

If there is a [`busybox`](https://www.busybox.net/) command in the search path, the non-shell built-in commands used in the tasks (such as `find`, `dirname`, `sort`) are executed as subcommands of `busybox` command (if busybox supports the command). So where these commands are not available, only the `busybox` command needs to be installed (eg [BusyBox for Windows](https://frippery.org/busybox/))

Busybox commands have limited functionality, but this functionality is available on all platforms. In order to support contributors using Windows, it is advisable to use the limited functionality of these commands. In this way only busybox needs to be installed on Windows. The author of the task definitions should therefore install busybox even if the Linux operating system is used. This is because busybox commands with limited functionality will be used during the creation/testing of the task definition.

Check [examples/busybox](examples/busybox/CONTRIBUTING.md) for more information on busybox support.

### Makefile

`Makefile` can be generated from task definitions using the `-m/--makefile` flag. The generated `Makefile` can be executed without `cdo`.

If the `make` command is issued without parameters, a short help is displayed about the targets that can be used.

It is important to note that the `Makefile` will not use cdo's embedded shell, but the `bash` shell.

