package analyzestruct_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/vaguecoder/fieldescription/pkg/analyzestruct"
)

const filename = `src.go`

func readerToASTNode(inputFile io.Reader) (ast.Node, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, inputFile, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func findStructNodes(node ast.Node) []ast.StructType {
	structs := make([]ast.StructType, 0)
	ast.Inspect(node, func(n ast.Node) bool {
		if sn, ok := n.(*ast.StructType); ok {
			structs = append(structs, *sn)
		}
		return true
	})

	return structs
}

func TestAnalyzeStruct(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:    "No-Struct",
			input:   `package main`,
			want:    []string{},
			wantErr: false,
		},
		{
			name: "No-Fields",
			input: `package main

			type MyStruct struct {}`,
			want:    []string{},
			wantErr: false,
		},
		{
			name: "No-Comments",
			input: `package main

			type MyStruct struct {
				ID int
				Name string
			}`,
			want:    []string{`ID`, `Name`},
			wantErr: false,
		},
		{
			name: "All-Fields-Compliant",
			input: `package main

			type MyStruct struct {
				// ID is identity
				ID int
				// Name is name
				Name string
			}`,
			want:    []string{},
			wantErr: false,
		},
		{
			name: "All-Fields-Compliant-Multi-Line-Comments",
			input: `package main

			type MyStruct struct {
				// ID is identity
				// which is the ID
				ID int
				// Name is name
				Name string
			}`,
			want:    []string{},
			wantErr: false,
		},
		{
			name: "All-Fields-Compliant-In-Line-Comments",
			input: `package main

			type MyStruct struct {
				// ID is identity
				// which is the ID
				ID int // ID comment
				// Name is name
				Name string // Name comment
			}`,
			want:    []string{},
			wantErr: false,
		},
		{
			name: "Multiple-Structs-All-Compliant",
			input: `package main

			type MyStruct struct {
				// ID is identity
				// which is the ID
				ID int // ID comment
				// Name is name
				Name string // Name comment
			}
			
			type YourStruct struct {
				// ID is identity
				ID int
				// Name is name
				Name string
			}`,
			want:    []string{},
			wantErr: false,
		},
		{
			name: "First-Field-Non-Compliant",
			input: `package mypkg
			
			// MyStruct is my struct
			type MyStruct struct {
				// off you must go
				F1 uint
				// ID is identity
				ID int
				// Name is name
				Name string
			}`,
			want:    []string{"F1"},
			wantErr: false,
		},
		{
			name: "Middle-Field-Non-Compliant",
			input: `package mypkg
			
			// MyStruct is my struct
			type MyStruct struct {
				// ID is identity
				ID int
				// off you must go
				F1 uint
				// Name is name
				Name string
			}`,
			want:    []string{"F1"},
			wantErr: false,
		},
		{
			name: "Last-Field-Non-Compliant",
			input: `package mypkg
			
			// MyStruct is my struct
			type MyStruct struct {
				// ID is identity
				ID int
				// Name is name
				Name string
				// off you must go
				F1 uint
			}`,
			want:    []string{"F1"},
			wantErr: false,
		},
		{
			name: "Multiple-Fields-Non-Compliant",
			input: `package mypkg
			
			// MyStruct is my struct
			type MyStruct struct {
				// off you must go
				F1 uint
				// ID is identity
				// which is the ID
				ID int // ID comment
				// Name is name
				Name string // Name comment
				// off you must go
				Age uint
				Height string
			}`,
			want:    []string{"F1", "Age", "Height"},
			wantErr: false,
		},
		{
			name: "Multiple-Structs-One-Non-Compliant",
			input: `package main

			type MyStruct struct {
				// ID is identity
				// which is the ID
				ID int // ID comment
				// Name is name
				Name string // Name comment
			}
			
			type YourStruct struct {
				ID int
				// Name is name
				Name string
			}`,
			want:    []string{`ID`},
			wantErr: false,
		},
		{
			name: "Multiple-Structs-All-Non-Compliant",
			input: `package main

			type MyStruct struct {
				// is identity
				// which is the ID
				ID int // ID comment
				// Name is name
				Name string // Name comment
			}
			
			type YourStruct struct {
				ID int
				// Name is name
				Name string
			}`,
			want:    []string{`ID`, `ID`},
			wantErr: false,
		},
		{
			name: "Multiple-Structs-All-Non-Compliant-All-Fields",
			input: `package main

			type MyStruct struct {
				// is identity
				// which is the ID
				ID int // ID comment
				Name string // Name comment
			}
			
			type YourStruct struct {
				ID int
				// is name
				Name string
			}`,
			want:    []string{`ID`, `Name`, `ID`, `Name`},
			wantErr: false,
		},
		{
			name: "Incomplete-Struct",
			input: `package main

			type MyStruct struct {
				// ID is identity
				ID int
				// Name is name
				Name string`,
			want:    []string{},
			wantErr: true,
		},
		{
			name: "Composition-All-Compliant",
			input: `package main

			type MyStruct struct {
				// ID is identity
				ID int
				// Name is name
				Name string
				YourStruct
			}
			
			type YourStruct struct {
				// ID is identity
				ID int
				// Name is name
				Name string
			}`,
			want:    []string{},
			wantErr: false,
		},
		{
			name: "Composition-External-Type-All-Compliant",
			input: `package main

			type MyStruct struct {
				// ID is identity
				ID int
				// Name is name
				Name string
				SomeStruct
			}`,
			want:    []string{},
			wantErr: false,
		},
		{
			name: "Composition-Non-Compliant",
			input: `package main

			type MyStruct struct {
				// ID is identity
				ID int
				// is name
				Name string
				YourStruct
			}`,
			want:    []string{`Name`},
			wantErr: false,
		},
		{
			name: "Unexported-Fields",
			input: `package main
		
			type MyStruct struct {
				// id is identity
				id int
				name string
			}`,
			want:    []string{`name`},
			wantErr: false,
		},
		{
			name: "Case-Mismatch",
			input: `package main
		
			type MyStruct struct {
				// ID is identity
				id int
				// Name is name
				name string
				// age in number
				Age uint
			}`,
			want:    []string{`id`, `name`, `Age`},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			node, err := readerToASTNode(reader)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ReaderToASTNode() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			missing := make([]string, 0)
			for _, structNode := range findStructNodes(node) {
				got, err := analyzestruct.FieldsMissingDescription(&structNode)
				if (err != nil) != tt.wantErr {
					t.Errorf("AnalyzeStruct() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				missing = append(missing, got...)
			}

			if !reflect.DeepEqual(missing, tt.want) {
				t.Errorf("AnalyzeStruct() = %v, want %v", missing, tt.want)
				return
			}
		})
	}
}
