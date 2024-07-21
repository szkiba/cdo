package cmd

import (
	"strings"

	"github.com/spf13/pflag"
)

func token2flag(src []string, fenv, ffile *pflag.Flag) []string {
	args := make([]string, 0, len(src))

	long := func(flag *pflag.Flag) string { return "--" + flag.Name }
	short := func(flag *pflag.Flag) string { return "-" + flag.Shorthand }

	isFile := false
	isEnv := false

	for _, arg := range src {
		if isFile {
			isFile = false
		} else {
			if arg == short(ffile) || arg == long(ffile) {
				isFile = true
			} else if arg[0] == '@' {
				args = append(args, long(ffile), arg[1:])

				continue
			}
		}

		if isEnv {
			isEnv = false
		} else {
			if arg == short(fenv) || arg == long(fenv) {
				isEnv = true
			} else if strings.ContainsRune(arg, '=') {
				args = append(args, long(fenv))
			}
		}

		args = append(args, arg)
	}

	return args
}
