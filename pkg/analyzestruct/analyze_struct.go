package analyzestruct

import (
	"go/ast"
	"go/token"
	"sort"
	"strings"
)

func FieldsMissingDescription(structNode *ast.StructType) ([]string, error) {
	var nodes []ast.Node
	var lastFieldEnd token.Pos
	ast.Inspect(structNode, func(m ast.Node) bool {
		switch y := m.(type) {
		case *ast.Field:
			nodes = append(nodes, y)
			lastFieldEnd = y.End()
		case *ast.Comment:
			if y.Pos() == lastFieldEnd+1 {
				return true
			}
			nodes = append(nodes, y)
		}
		return true
	})

	sortNodes(nodes)

	return findFieldsMissingDesc(nodes), nil
}

func sortNodes(nodes []ast.Node) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Pos() < nodes[j].Pos()
	})
}

func findFieldsMissingDesc(nodes []ast.Node) []string {
	missing := make([]string, 0)
	comments := make([]ast.Comment, 0)
	for _, node := range nodes {
		switch nodeType := node.(type) {
		case *ast.Field:
			if len(nodeType.Names) == 0 {
				continue
			}

			field := nodeType.Names[0].Name
			if !fieldDescExists(field, comments) {
				missing = append(missing, field)
			}
			comments = make([]ast.Comment, 0)
		case *ast.Comment:
			comments = append(comments, *nodeType)
		}
	}

	return missing
}

func fieldDescExists(field string, comments []ast.Comment) bool {
	for _, comment := range comments {
		if strings.HasPrefix(comment.Text, "// "+field) {
			return true
		}
	}

	return false
}
