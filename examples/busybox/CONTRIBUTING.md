# Example: busybox

This example shows how to use busybox commands in tasks.

If there is a [`busybox`](https://www.busybox.net/) command in the search path, the non-shell built-in commands used in the task (`find`, `dirname`, `sort`) are executed as subcommands of `busybox` (if busybox supports the command). So where these commands are not available, only the `busybox` command needs to be installed (eg [BusyBox for Windows](https://frippery.org/busybox/)).

## list - List doc directories

List the directories containing markdown files.
A useless task to demonstrate how to use busybox commands.

The base directory can be specified as the first positional parameter.
The default base directory is `../..`

```bash
find ${1:-../..} -name '*.md' | while read name; do
dirname $name
done | sort -u
```

