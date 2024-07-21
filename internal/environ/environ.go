package environ

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"mvdan.cc/sh/v3/expand"
)

type Environ map[string]string

func New(lines []string) Environ {
	e := make(Environ)

	e.parse(lines)

	return e
}

//nolint:exhaustruct
func (e Environ) Get(name string) expand.Variable {
	value, has := e[name]
	if !has {
		return expand.Variable{Kind: expand.Unset}
	}

	return expand.Variable{Exported: true, Kind: expand.String, Str: value}
}

//nolint:exhaustruct
func (e Environ) Each(fn func(name string, vr expand.Variable) bool) {
	for key, value := range e {
		if !fn(key, expand.Variable{Exported: true, Kind: expand.String, Str: value}) {
			return
		}
	}
}

func (e Environ) String() string {
	return ""
}

func (e Environ) Set(line string) error {
	dict, err := godotenv.Unmarshal(line)
	if err != nil {
		return err
	}

	for key, value := range dict {
		e[key] = value
	}

	return nil
}

func (e Environ) Type() string {
	return "name=value"
}

func (e Environ) parse(lines []string) {
	for _, line := range lines {
		if idx := strings.Index(line, "="); idx >= 0 {
			e[line[:idx]] = line[idx+1:]
		}
	}
}

func (e Environ) Load(dir string) error {
	if err := e.loadDotenv(filepath.Join(dir, ".env")); err != nil {
		return err
	}

	return e.loadDotenv(filepath.Join(dir, ".env.local"))
}

func (e Environ) Override(env Environ) {
	for key, value := range env {
		e[key] = value
	}
}

func (e Environ) loadDotenv(filename string) error {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	all, err := godotenv.Read(filename)
	if err != nil {
		return err
	}

	for key, value := range all {
		e[key] = value
	}

	return nil
}
