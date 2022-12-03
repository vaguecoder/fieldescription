package analyzer

import (
	"flag"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

type FieldDescriptionAnalyzer struct {
	lastStruct string
	lastEnd    token.Pos
}

func NewAnalyzer() *analysis.Analyzer {
	fdAnalyzer := &FieldDescriptionAnalyzer{}
	return &analysis.Analyzer{
		Name:             project,
		Doc:              documentation,
		Flags:            fdAnalyzer.newFlagSet(),
		Run:              fdAnalyzer.analyzeFieldDesc,
		RunDespiteErrors: false,
		Requires:         []*analysis.Analyzer{inspect.Analyzer},
		ResultType:       nil,
		FactTypes:        []analysis.Fact{},
	}
}

func (f *FieldDescriptionAnalyzer) newFlagSet() flag.FlagSet {
	fs := flag.NewFlagSet("exhaustruct flags", flag.PanicOnError)

	return *fs
}
