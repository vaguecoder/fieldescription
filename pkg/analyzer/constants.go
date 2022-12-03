package analyzer

import "fmt"

const (
	project     = `fieldescription`
	description = `Linter to check for all struct fields comments`
	usage       = `fieldescription <file.go>
	fieldescription ./<path-to-package>/<file.go>
	fieldescription <file1.go> <file2.go>
	fieldescription ./...
	fieldescription ./<path-to-package>/...
	`
)

var (
	documentation = fmt.Sprintf("%s\n\n", project)
)
