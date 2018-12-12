package findfuncbody_test

import (
	"go/ast"
	"path/filepath"
	"testing"

	"github.com/tenntenn/gosa/passes/findfuncbody"
	"github.com/tenntenn/gosa/passes/findfuncbody/funcbody"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		findfuncbody.Analyzer,
	},
	Run: run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	funcbody := pass.ResultOf[findfuncbody.Analyzer].(*funcbody.Finder)

	nodeFilter := []ast.Node{
		(*ast.FuncLit)(nil),
		(*ast.Ident)(nil),
		(*ast.UnaryExpr)(nil),
		(*ast.StarExpr)(nil),
		(*ast.CallExpr)(nil),
		(*ast.IndexExpr)(nil),
		(*ast.SelectorExpr)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {

		expr, ok := n.(ast.Expr)
		if !ok {
			return
		}

		body := funcbody.Find(expr)
		if body != nil {
			pos := pass.Fset.Position(body.Parent.Pos())
			pass.Reportf(expr.Pos(), "body is %s:%d", filepath.Base(pos.Filename), pos.Line)
		}
	})

	return nil, nil
}

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "a")
}
