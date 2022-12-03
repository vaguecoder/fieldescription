package analyzer

import (
	"go/ast"
	"strings"

	"github.com/vaguecoder/fieldescription/pkg/analyzestruct"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

func (f *FieldDescriptionAnalyzer) analyzeFieldDesc(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeTypes := []ast.Node{
		(*ast.StructType)(nil),
		(*ast.Ident)(nil),
	}

	insp.Preorder(nodeTypes, f.traverse(pass))

	return nil, nil
}

func (f *FieldDescriptionAnalyzer) traverse(pass *analysis.Pass) func(node ast.Node) {
	return func(node ast.Node) {
		switch typedNode := node.(type) {
		case *ast.StructType:
			var structName string
			if f.lastEnd+1 == node.Pos() {
				structName = f.lastStruct
			}

			missing, err := analyzestruct.FieldsMissingDescription(typedNode)
			if err != nil {
				pass.Reportf(node.Pos(), "failed to analyze struct with identifier %q: %v", structName, err)
				return
			}

			if len(missing) > 0 {
				missingFieldsStr := strings.Join(missing, ", ")
				pass.Reportf(node.Pos(), "missing field description for fields for struct with identifier %q: %s", structName, missingFieldsStr)
			}
		case *ast.Ident:
			f.lastStruct = typedNode.Name
			f.lastEnd = typedNode.End()
		}
	}
}
