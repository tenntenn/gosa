package funcbody

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/tenntenn/gosa/dog"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
)

type FuncBody struct {
	Body   *ast.BlockStmt
	Parent ast.Node
}

type Finder struct {
	Fset      *token.FileSet
	Files     []*ast.File
	TypesInfo *types.Info
	Pkg       *types.Package
	SSA       *ssa.Package
}

func (f *Finder) fileByPos(pos token.Pos) *ast.File {
	for _, file := range f.Files {
		if file.Pos() <= pos && pos <= file.End() {
			return file
		}
	}
	return nil
}

func (f *Finder) Find(expr ast.Expr) *FuncBody {
	file := f.fileByPos(expr.Pos())
	path, exact := astutil.PathEnclosingInterval(file, expr.Pos(), expr.End())
	if exact {
		return f.FindWithPath(expr, path)
	}
	return nil
}

func (f *Finder) FindWithPath(expr ast.Expr, path []ast.Node) *FuncBody {
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
		switch obj := f.TypesInfo.ObjectOf(expr).(type) {
		case *types.Func: // named function or method
			fun := f.SSA.Prog.FuncValue(obj)
			return f.findByValue(fun)
		case *types.Var: // variable or field
			v, _ := f.SSA.Prog.VarValue(obj, f.SSA, path)
			return f.findByValue(v)
		}
	case *ast.CallExpr:
		return f.Find(expr.Fun)
	case *ast.SelectorExpr:
		path = append([]ast.Node{expr.Sel}, path...)
		return f.FindWithPath(expr.Sel, path)
	}

	return nil
}

func (f *Finder) isFunc(expr ast.Expr) bool {
	_, ok := f.TypesInfo.TypeOf(expr).(*types.Signature)
	return ok
}

func (f *Finder) findByValue(v ssa.Value) *FuncBody {
	switch v := v.(type) {
	case *ssa.Function:
		return f.fromNode(v.Syntax())
	case *ssa.UnOp: // pointer
		if v.Op == token.MUL {
			return f.findByValue(v.X)
		}
	case *ssa.FieldAddr:
		/*
			p := f.analyzePtr(v.X)
			for _, l := range p.PointsTo().Labels() {
				fmt.Printf("%#v\n", l.Value())
			}
			fmt.Println()
		*/
		//return f.findByValue(v.X)
	case *ssa.Field:
		dog.DumpReferrers(v, -1)
		//return f.findByValue(v.X)
	case *ssa.Alloc:
		//p := f.analyzePtr(v)
		//for _, l := range p.PointsTo().Labels() {
		//	fmt.Printf("%T\n", l.Value())
		//}
		//fmt.Println()
	}
	return nil
}

func (f *Finder) fromNode(n ast.Node) *FuncBody {
	switch n := n.(type) {
	case *ast.FuncLit:
		return &FuncBody{
			Parent: n,
			Body:   n.Body,
		}
	case *ast.FuncDecl:
		return &FuncBody{
			Parent: n,
			Body:   n.Body,
		}
	}
	return nil
}

func (f *Finder) fromTypesObj(obj types.Object, path []ast.Node) *FuncBody {
	switch obj := obj.(type) {
	case *types.Func: // named function or method
		fun := f.SSA.Prog.FuncValue(obj)
		return f.findByValue(fun)
	case *types.Var: // variable or field
		v, _ := f.SSA.Prog.VarValue(obj, f.SSA, path)
		return f.findByValue(v)
	}
	return nil
}

func (f *Finder) findStore(v ssa.Value, refs *[]ssa.Instruction) *ssa.Store {
	if refs == nil {
		return nil
	}

	for _, ref := range *refs {
		switch ref := ref.(type) {
		case *ssa.Store:
			if ref.Addr == v {
				return ref
			}
			return f.findStore(v, ref.Referrers())
		case *ssa.DebugRef:
			fmt.Printf("store %#v\n", ref)
		case ssa.Node:
			return f.findStore(v, ref.Referrers())
		}
	}

	return nil
}

func (f *Finder) findFromStoreRefs(refs *[]ssa.Instruction) *FuncBody {
	if refs == nil {
		return nil
	}

	for _, ref := range *refs {
		switch ref := ref.(type) {
		case *ssa.Store:
			switch v := ref.Val.(type) {
			case *ssa.Function:
				switch n := v.Syntax().(type) {
				case *ast.FuncLit:
					return &FuncBody{
						Parent: n,
						Body:   n.Body,
					}
				case *ast.FuncDecl:
					return &FuncBody{
						Parent: n,
						Body:   n.Body,
					}
				}
			}
		}
	}
	return nil
}

func (f *Finder) analyzePtr(v ssa.Value) (ptr pointer.Pointer) {
	config := &pointer.Config{
		Mains: []*ssa.Package{f.SSA},
	}
	config.AddQuery(v)
	result, err := pointer.Analyze(config)
	if err != nil {
		return
	}
	return result.Queries[v]
}
