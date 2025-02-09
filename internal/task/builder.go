package task

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"

	"github.com/iancoleman/strcase"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

type builder struct {
	tasks      []*Task
	task       *Task
	level      int
	source     []byte
	startIndex int
	endIndex   int
	term       string
	options    map[string]string
}

func newBuilder(source []byte) *builder {
	b := new(builder)

	b.source = source

	return b
}

func (b *builder) add() {
	if b.task == nil {
		return
	}

	getopts(b.task, b.options)

	b.tasks = append(b.tasks, b.task)

	b.task = nil
	b.options = nil
}

func (b *builder) build() (map[string]*Task, error) {
	b.add()

	tasks := make(map[string]*Task, len(b.tasks))

	for _, task := range b.tasks {
		tasks[task.Name] = task
	}

	lookup := func(name string) (bool, [][]string) {
		task, has := tasks[name]
		if !has {
			return false, nil
		}

		return true, task.Requires
	}

	for _, task := range b.tasks {
		for _, req := range task.Requires {
			visited := map[string]struct{}{task.Name: {}}
			if err := checkdep(req[0], lookup, visited); err != nil {
				return nil, err
			}
		}
	}

	return tasks, nil
}

func checkdep(name string, lookup func(string) (bool, [][]string), visited map[string]struct{}) error {
	if _, done := visited[name]; done {
		return fmt.Errorf("%w: %s", errRequiresCycle, name)
	}

	found, deps := lookup(name)
	if !found {
		return fmt.Errorf("%w: %s", errMissingTask, name)
	}

	visited[name] = struct{}{}

	for _, dep := range deps {
		if err := checkdep(dep[0], lookup, visited); err != nil {
			return err
		}
	}

	return nil
}

func extractNameShort(contents []byte) (bool, string, string) {
	const fieldNum = 2

	fields := bytes.SplitN(contents, separator, fieldNum)

	if len(fields) != fieldNum {
		return false, "", ""
	}

	return true,
		strcase.ToKebab(string(bytes.TrimSpace(fields[0]))),
		string(bytes.TrimSpace(fields[1]))
}

func (b *builder) handleDefinitionList(node ast.Node, entering bool) {
	if asDefinitionList(node, entering) != nil {
		b.term = ""

		return
	}

	if desc := asDefinitionDescription(node, entering); desc != nil {
		if len(b.term) != 0 {
			b.options[b.term] = string(desc.Text(b.source))

			b.term = ""
		}

		return
	}

	if term := asDefinitionTerm(node, entering); term != nil && b.task != nil {
		b.term = string(term.Text(b.source))
	}
}

func (b *builder) handleHeading(node ast.Node, entering bool) bool {
	heading := asHeading(node, entering)
	if heading != nil {
		if heading.Level <= b.level {
			b.add()
		}

		contents := extractBlock(heading.Lines(), b.source)

		match, name, short := extractNameShort(contents)
		if match {
			b.add()
			b.task = &Task{Name: name, Short: short}
			b.startIndex = heading.Lines().At(0).Start
			b.level = heading.Level
			b.options = make(map[string]string)
		}

		return true
	}

	return false
}

func (b *builder) handleCodeBlock(node ast.Node, entering bool) error {
	if b.task == nil {
		return nil
	}

	fcb := asFencedCodeBlock(node, entering)
	if entering || fcb == nil || fcb.Info == nil {
		return nil
	}

	const fencePrefixLen = 3

	b.endIndex = fcb.Info.Segment.Start - fencePrefixLen

	found, script, err := extractScript(fcb, b.source)
	if !found || err != nil {
		return err
	}

	b.task.Script = script

	b.task.Long = string(bytes.TrimSpace(b.source[b.startIndex:b.endIndex]))

	b.add()

	b.task = nil
	b.options = nil

	return nil
}

func (b *builder) walk(node ast.Node, entering bool) (ast.WalkStatus, error) {
	b.handleDefinitionList(node, entering)

	if b.handleHeading(node, entering) {
		return ast.WalkContinue, nil
	}

	err := b.handleCodeBlock(node, entering)

	return ast.WalkContinue, err
}

var reInfo = regexp.MustCompile(`\s*(\w+)\s*(.*)\s*`)

var separator = []byte{' ', '-', ' '} //nolint:gochecknoglobals

func asHeading(node ast.Node, entering bool) *ast.Heading {
	if entering || node.Kind() != ast.KindHeading {
		return nil
	}

	if heading, ok := node.(*ast.Heading); ok {
		return heading
	}

	return nil
}

func asFencedCodeBlock(node ast.Node, entering bool) *ast.FencedCodeBlock {
	if entering || node.Kind() != ast.KindFencedCodeBlock {
		return nil
	}

	if fcb, ok := node.(*ast.FencedCodeBlock); ok {
		return fcb
	}

	return nil
}

func extractScript(fcb *ast.FencedCodeBlock, source []byte) (bool, []byte, error) {
	lang, err := extractInfo(fcb, source)
	if err != nil {
		return false, nil, err
	}

	if lang != "bash" && lang != "sh" {
		return false, nil, nil
	}

	return true, extractBlock(fcb.Lines(), source), nil
}

func extractBlock(lines *text.Segments, source []byte) []byte {
	var buff bytes.Buffer

	for i := 0; i < lines.Len(); i++ {
		seg := lines.At(i)

		buff.Write(seg.Value(source))
	}

	return buff.Bytes()
}

func extractInfo(fcb *ast.FencedCodeBlock, source []byte) (string, error) {
	if fcb.Info == nil {
		return "", nil
	}

	return parseInfo(fcb.Info.Text(source))
}

func parseInfo(text []byte) (string, error) {
	all := reInfo.FindSubmatch(text)
	if all == nil {
		return "", nil
	}

	var (
		lang string
		err  error
	)

	if len(all) > 1 {
		lang = string(all[1])
	}

	if len(all) <= reInfo.NumSubexp() {
		return lang, nil
	}

	return lang, err
}

func asDefinitionTerm(node ast.Node, entering bool) *east.DefinitionTerm {
	if entering || node.Kind() != east.KindDefinitionTerm {
		return nil
	}

	if term, ok := node.(*east.DefinitionTerm); ok {
		return term
	}

	return nil
}

func asDefinitionDescription(node ast.Node, entering bool) *east.DefinitionDescription {
	if entering || node.Kind() != east.KindDefinitionDescription {
		return nil
	}

	if desc, ok := node.(*east.DefinitionDescription); ok {
		return desc
	}

	return nil
}

func asDefinitionList(node ast.Node, entering bool) *east.DefinitionList {
	if entering || node.Kind() != east.KindDefinitionList {
		return nil
	}

	if list, ok := node.(*east.DefinitionList); ok {
		return list
	}

	return nil
}

var (
	errRequiresCycle = errors.New("requires cycle")
	errMissingTask   = errors.New("missing task")
)
