package shell

import (
	"bytes"
	"context"
	"os"

	"github.com/szkiba/cdo/internal/environ"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func Run(task string, args []string, script []byte, dir string, env environ.Environ) error {
	file, _ := syntax.NewParser().Parse(bytes.NewReader(script), task)
	params := []string{"-e", "--"}
	params = append(params, args...)
	runner, _ := interp.New(
		interp.StdIO(os.Stdin, os.Stdout, os.Stdout),
		interp.Params(params...),
		interp.Env(env),
		interp.Dir(dir),
		interp.ExecHandlers(busyBoxHandler(dir, env)),
	)

	return runner.Run(context.TODO(), file)
}
