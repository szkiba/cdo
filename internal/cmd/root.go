package cmd

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/szkiba/cdo/internal/environ"
	"github.com/szkiba/cdo/internal/shell"
	"github.com/szkiba/cdo/internal/task"
)

//nolint:gochecknoglobals
var (
	appname = "cdo"
	version = "dev"
)

func findReadme() (string, string, error) {
	const base = "README.md"

	abs, err := filepath.Abs(base)
	if err != nil {
		return "", "", err
	}

	for dir := filepath.Dir(abs); ; dir = filepath.Dir(dir) {
		filename := filepath.Clean(filepath.Join(dir, base))
		if _, err := os.Stat(filename); err == nil {
			return filename, dir, nil
		}

		if dir[len(dir)-1] == filepath.Separator {
			break
		}
	}

	return "", "", os.ErrNotExist
}

func findContributing() (string, string, error) {
	const base = "CONTRIBUTING.md"

	abs, err := filepath.Abs(base)
	if err != nil {
		return "", "", err
	}

	for dir := filepath.Dir(abs); ; dir = filepath.Dir(dir) {
		files := []string{
			filepath.Clean(filepath.Join(dir, base)),
			filepath.Clean(filepath.Join(dir, "docs", base)),
			filepath.Clean(filepath.Join(dir, ".github", base)),
		}

		for _, filename := range files {
			if _, err := os.Stat(filename); err == nil {
				return filename, dir, nil
			}
		}

		if dir[len(dir)-1] == filepath.Separator {
			break
		}
	}

	return "", "", os.ErrNotExist
}

func findDefinitions() (string, string, error) {
	filename, dir, err := findContributing()
	if errors.Is(err, os.ErrNotExist) {
		filename, dir, err = findReadme()
	}

	return filename, dir, err
}

//go:embed help.txt
var help string

func newCommand() *cobra.Command {
	root := &cobra.Command{
		Use:               appname + " [flags] [task]",
		Version:           version,
		Short:             "Markdown-based task runner for contributors",
		Long:              strings.ReplaceAll(help, "cdo", appname),
		SilenceUsage:      true,
		SilenceErrors:     true,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	root.SetHelpCommand(&cobra.Command{Hidden: true})
	root.Flags().BoolP("help", "h", false, "Print usage")
	root.Flags().BoolP("version", "V", false, "Print version")
	root.SetUsageTemplate(usageTemplate)

	return root
}

func New(args []string) (*cobra.Command, error) {
	env := environ.New(os.Environ())
	flagenv := environ.New(nil)

	filename, dir, err := findDefinitions()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	root := newCommand()
	root.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		if err := env.Load(dir); err != nil {
			return err
		}

		env.Override(flagenv)

		return nil
	}

	flags := root.PersistentFlags()

	flags.VarP(&flagenv, "env", "e", "Set environment variable(s)")
	flags.StringVarP(&filename, "file", "f", filename, "Task definitions file")

	args = token2flag(args, flags.Lookup("env"), flags.Lookup("file"))

	root.SetArgs(args)

	flags.ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{UnknownFlags: true}
	flags.SetOutput(io.Discard)

	err = flags.Parse(args)
	if err != nil && !errors.Is(err, pflag.ErrHelp) {
		return nil, err
	}

	if fflag := flags.Lookup("file"); fflag.Changed {
		fflag.DefValue = filename
		dir = filepath.Dir(filename)
	}

	if len(filename) != 0 {
		if err := addCommands(root, env, filename, dir); err != nil {
			return nil, err
		}
	} else {
		root.RunE = runNoFile
	}

	return root, nil
}

func runNoFile(cmd *cobra.Command, _ []string) error {
	if err := cmd.Help(); err != nil {
		return err
	}

	fmt.Fprintln(cmd.ErrOrStderr())

	return errNoFile
}

func runRequires(task *task.Task, cmd *cobra.Command) error {
	for _, req := range task.Requires {
		rcmd, rargs, err := cmd.Root().Find(req)
		if err != nil {
			return err
		}

		if err := rcmd.RunE(rcmd, rargs); err != nil {
			return err
		}
	}

	return nil
}

func addCommands(cmd *cobra.Command, env environ.Environ, filename string, dir string) error {
	taskdefs, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		return err
	}

	tasks, err := task.Load(taskdefs)
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		return fmt.Errorf("%w in %s", errNoTasks, filename)
	}

	for _, task := range tasks {
		sub := &cobra.Command{
			Use:                task.Name,
			Short:              task.Short,
			Long:               task.Long,
			FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		}

		if len(task.Script) != 0 || len(task.Requires) != 0 {
			sub.RunE = func(cmd *cobra.Command, args []string) error {
				if err := runRequires(task, cmd); err != nil {
					return err
				}

				if len(task.Script) != 0 {
					return shell.Run(cmd.Name(), args, task.Script, dir, env)
				}

				return nil
			}
		}

		sub.Flags().BoolP("help", "h", false, "Print usage")

		cmd.AddCommand(sub)

		sub.Long = strings.Replace(sub.Long, task.Name, sub.CommandPath(), 1)
	}

	return nil
}

var (
	errNoTasks = errors.New("no task definitions")
	errNoFile  = errors.New("no task definition file found, use the --file flag to specify one")
)

const usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [task]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Tasks:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Tasks:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [task] --help" for more information about a task.{{end}}
`
