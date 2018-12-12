package funcbody

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
)

type FuncBody struct {
	Body   *ast.BlockStmt
	Parent ast.Node
}

type Finder struct {
	Fset      *token.FileSet
	TypesInfo *types.Info
}

func (f *Finder) Find(expr ast.Expr) *FuncBody {
	if !f.isFunc(expr) {
		return nil
	}

	switch expr := expr.(type) {
	case *ast.FuncLit:
		return &FuncBody{
			Parent: expr,
			Body:   expr.Body,
		}
	case *ast.Ident:
		obj := expr.Obj
		if obj == nil {
			pos := f.Fset.Position(expr.Pos())
			if filepath.Base(pos.Filename) == "a.go" && pos.Line == 13 {
				fmt.Println(f.defs(expr))
			}
			if id := f.defs(expr); id != nil {
				obj = id.Obj
			}
		}

		if obj == nil {
			return nil
		}

		switch decl := obj.Decl.(type) {
		case *ast.FuncDecl:
			return &FuncBody{
				Parent: decl,
				Body:   decl.Body,
			}
		case *ast.AssignStmt:
			for i := range decl.Lhs {
				if body := f.Find(decl.Rhs[i]); body != nil {
					return body
				}
			}
		}
		return nil
	case *ast.CallExpr:
		return f.Find(expr.Fun)
	case *ast.SelectorExpr:
		return f.Find(expr.Sel)
	}

	return nil
}

func (f *Finder) isFunc(expr ast.Expr) bool {
	_, ok := f.TypesInfo.Types[expr].Type.(*types.Signature)
	return ok
}

func (f *Finder) defs(ident *ast.Ident) *ast.Ident {
	for id := range f.TypesInfo.Defs {
		if f.identEquals(ident, id) {
			return id
		}
	}
	return nil
}

func (f *Finder) identEquals(id1, id2 *ast.Ident) bool {
	return id1 == id2 || (id1.Name == id2.Name &&
		f.TypesInfo.ObjectOf(id1).Parent() == f.TypesInfo.ObjectOf(id2).Parent())
}
