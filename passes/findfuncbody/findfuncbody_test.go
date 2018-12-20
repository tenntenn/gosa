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
	finder := pass.ResultOf[findfuncbody.Analyzer].(*funcbody.Finder)

	nodeFilter := []ast.Node{
		(*ast.FuncLit)(nil),
		(*ast.Ident)(nil),
		(*ast.UnaryExpr)(nil),
		(*ast.StarExpr)(nil),
		(*ast.CallExpr)(nil),
		(*ast.IndexExpr)(nil),
		(*ast.SelectorExpr)(nil),
	}

	type key struct {
		filename string
		line     int
	}

	type value struct {
		body *funcbody.FuncBody
		expr ast.Expr
	}

	found := map[key]*value{}
	inspect.Preorder(nodeFilter, func(n ast.Node) {

		expr, ok := n.(ast.Expr)
		if !ok {
			return
		}

		body := finder.Find(expr)
		pos := pass.Fset.Position(expr.Pos())
		key := key{line: pos.Line, filename: pos.Filename}
		if body != nil && found[key] == nil {
			found[key] = &value{
				expr: expr,
				body: body,
			}
		}
	})

	for _, v := range found {
		bodyPos := pass.Fset.Position(v.body.Parent.Pos())
		filename := filepath.Base(bodyPos.Filename)
		pass.Reportf(v.expr.Pos(), "body is %s:%d", filename, bodyPos.Line)
	}

	return nil, nil
}

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "a")
}
