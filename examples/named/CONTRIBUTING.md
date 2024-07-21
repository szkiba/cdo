# Example: named

This example shows how to use named parameters in tasks.


## build - Build an executable

This task builds the executable binary for the OS and CPU specified as parameters (of course not, it just displays the parameters). Both parameters have a default (`linux` operating system and `amd64` processor)

```bash
echo "operating system: ${os:-linux}"
echo "processor: ${cpu:-amd64}"
```

The `os` and `cpu` parameters can be specified as variable assignments on the command line:

```shell
cdo build os=windows cpu=arm64
```

