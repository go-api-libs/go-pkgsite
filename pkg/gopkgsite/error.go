package gopkgsite

import (
	"fmt"
	"strings"
)

func (e *Error) Error() string {
	b := &strings.Builder{}
	fmt.Fprintf(b, "%d: %s", e.Code, e.Message)
	if len(e.Fixes) > 0 {
		b.WriteString("\n\nFixes:")
		for _, fix := range e.Fixes {
			fmt.Fprintf(b, "\n- %s", fix)
		}
	}

	if len(e.Candidates) > 0 {
		b.WriteString("\n\nCandidates:")
		for _, c := range e.Candidates {
			fmt.Fprintf(b, "\n- pkg: %s, module: %s", c.PackagePath, c.ModulePath)
		}
	}

	return b.String()
}
