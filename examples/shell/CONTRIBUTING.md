# Example: shell

This example shows how to use advanced shell features such as function definition, `if/endif` statement or `for/do/done` loop.

## files - Create markdown file list

Creating a Markdown file list with links. A pointless task, it only demonstrates the use of advanced shell features.

```sh
# Yes, functions can be defined and control statements can also be used.
markdownFileList() {
    if [[ -n $1 ]]; then
      cd $1
    fi
    for i in *; do
      echo "- [$i](${1:-.}/$i)"
    done
}

# List the contents of the current directory 
markdownFileList
echo "---"
# List the contents of the directory passed as a parameter
markdownFileList ../..
```

The output:

```markdown
- [CONTRIBUTING.md](./CONTRIBUTING.md)
---
- [CODE_OF_CONDUCT.md](../../CODE_OF_CONDUCT.md)
- [CONTRIBUTING.md](../../CONTRIBUTING.md)
- [LICENSE](../../LICENSE)
- [README.md](../../README.md)
- [build](../../build)
- [coverage.txt](../../coverage.txt)
- [examples](../../examples)
- [go.mod](../../go.mod)
- [go.sum](../../go.sum)
- [internal](../../internal)
- [main.go](../../main.go)
```