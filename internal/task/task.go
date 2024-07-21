package task

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type Task struct {
	Name     string
	Short    string
	Long     string
	Script   []byte
	Requires [][]string
}

func Load(taskdefs []byte) (map[string]*Task, error) {
	parser := newParser()
	reader := text.NewReader(taskdefs)
	root := parser.Parse(reader).OwnerDocument()
	builder := newBuilder(taskdefs)

	err := ast.Walk(root, builder.walk)
	if err != nil {
		return nil, err
	}

	return builder.build()
}

func newParser() parser.Parser { //nolint:ireturn
	const (
		defListPriority = 101
		defDescPriority = 102
	)

	return parser.NewParser(
		parser.WithBlockParsers(parser.DefaultBlockParsers()...),
		parser.WithBlockParsers(
			util.Prioritized(extension.NewDefinitionListParser(), defListPriority),
			util.Prioritized(extension.NewDefinitionDescriptionParser(), defDescPriority),
		),
		parser.WithInlineParsers(parser.DefaultInlineParsers()...),
		parser.WithParagraphTransformers(parser.DefaultParagraphTransformers()...),
	)
}
