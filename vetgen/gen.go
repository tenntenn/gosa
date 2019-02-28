package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type PkgInfo struct {
	Pkg        string
	ImportPath string
}

type VetGen struct {
	ToolName string
	ExeName  string
	Args     []string
}

func (g *VetGen) Run() error {
	if len(g.Args[0]) == 0 {
		return errors.New("subcommand must be specified: init, add")
	}
	cmd, args := g.Args[0], g.Args[1:]

	switch cmd {
	case "init":
		return g.init(args)
	case "add":
		return g.add(args)
	}

	return fmt.Errorf("unsuported command %s", cmd)
}

func (g *VetGen) init(args []string) error {
	if len(args) == 0 {
		return errors.New("directory must be specified")
	}
	dir := args[0]

	info, err := os.Stat(dir)
	switch {
	case os.IsNotExist(err):
		if err := os.MkdirAll(dir, 0777); err != nil {
			return err
		}
	case err != nil:
		return err
	case !info.IsDir():
		return fmt.Errorf("%s is already exist but it is not directory", dir)
	default:
		var yorn string
		fmt.Printf("%s is already exsit, overwrite?[y/N] >", dir)
		fmt.Scan(&yorn)
		switch strings.ToLower(yorn) {
		default:
			return nil
		case "y", "yes":
			// noop
		}
	}

	mainFile := filepath.Join(dir, "main.go")
	if err := g.generate(mainFile, nil); err != nil {
		return err
	}

	return nil
}

func (g *VetGen) add(args []string) error {

	if len(args) == 0 {
		return errors.New("import path must be specified")
	}

	pkgs, err := g.importedPkgs()
	if err != nil {
		return err
	}

	pkgs = append(pkgs, &PkgInfo{
		Pkg:        path.Base(args[0]),
		ImportPath: args[0],
	})

	if err := g.generate("main.go", pkgs); err != nil {
		return err
	}

	return nil
}

func (g *VetGen) importedPkgs() ([]*PkgInfo, error) {

	srcMain, err := os.Open("main.go")
	if err != nil {
		return nil, err
	}
	defer srcMain.Close()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "main.go", srcMain, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var pkgs []*PkgInfo
	for _, importSpec := range f.Imports {
		var comment string
		if importSpec.Comment != nil {
			comment = strings.TrimSpace(importSpec.Comment.Text())
		}

		if comment == "add by vetgen" {
			importPath, err := strconv.Unquote(importSpec.Path.Value)
			if err != nil {
				return nil, err
			}
			pkgs = append(pkgs, &PkgInfo{
				Pkg:        path.Base(importPath),
				ImportPath: importPath,
			})
		}
	}

	return pkgs, nil
}

func (g *VetGen) generate(mainFile string, pkgs []*PkgInfo) error {

	var pkgset []*PkgInfo
	done := map[string]bool{}
	for _, p := range pkgs {
		if !done[p.ImportPath] {
			done[p.ImportPath] = true
			pkgset = append(pkgset, p)
		}
	}

	var buf bytes.Buffer
	if err := srcTempl.Execute(&buf, pkgset); err != nil {
		return err
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	if err := os.Remove(mainFile); err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := ioutil.WriteFile(mainFile, src, 0644); err != nil {
		return err
	}

	return nil
}
