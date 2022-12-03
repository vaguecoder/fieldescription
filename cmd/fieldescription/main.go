package main

import (
	"github.com/vaguecoder/fieldescription/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	fieldDescAnalyzer := analyzer.NewAnalyzer()
	singlechecker.Main(fieldDescAnalyzer)
}
