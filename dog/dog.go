package dog

import (
	"fmt"
	"go/token"
	"strings"

	"golang.org/x/tools/go/ssa"
)

func In(fset *token.FileSet, pos token.Pos, file string, line int) bool {
	f := fset.File(pos)
	return strings.HasSuffix(f.Name(), file) && f.Line(pos) == line
}

func DumpReferrers(node ssa.Node, limit int) {
	dumpReferrers(node, limit, 0)
}

func dumpReferrers(node ssa.Node, limit, n int) {
	refs := node.Referrers()
	if limit == n || refs == nil {
		return
	}

	for _, ref := range *refs {
		fmt.Printf("%s%T\n", strings.Repeat(" ", n), ref)
		switch ref := ref.(type) {
		case ssa.Node:
			dumpReferrers(ref, limit, n+1)
		}
	}
}
