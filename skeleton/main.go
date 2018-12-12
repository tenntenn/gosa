package main

import (
	"fmt"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

type PkgInfo struct {
	Pkg        string
	ImportPath string
}

func run(args []string) error {

	var info PkgInfo

	if len(args) < 2 {
		return errors.New("package must be specified")
	}
	info.Pkg = args[1]

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, gopath := range filepath.SplitList(build.Default.GOPATH) {
		if gopath == "" {
			continue
		}

		src := filepath.Join(gopath, "src")
		if strings.HasPrefix(cwd, src) {
			rel, err := filepath.Rel(src, cwd)
			if err != nil {
				return err
			}
			info.ImportPath = path.Join(filepath.ToSlash(rel), info.Pkg)
			break
		}
	}

	if info.ImportPath == "" {
		return errors.Errorf("%s must be executed in GOPATH", args[0])
	}

	dir := filepath.Join(cwd, info.Pkg)
	if err := os.Mkdir(dir, 0777); err != nil {
		return err
	}

	src, err := os.Create(filepath.Join(dir, info.Pkg+".go"))
	if err != nil {
		return err
	}
	defer src.Close()
	if err := srcTempl.Execute(src, info); err != nil {
		return err
	}

	test, err := os.Create(filepath.Join(dir, info.Pkg+"_test.go"))
	if err != nil {
		return err
	}
	defer test.Close()
	if err := testTempl.Execute(test, info); err != nil {
		return err
	}

	testdata := filepath.Join(dir, "testdata", "src", "a")
	if err := os.MkdirAll(testdata, 0777); err != nil {
		return err
	}

	adotgo, err := os.Create(filepath.Join(testdata, "a.go"))
	if err != nil {
		return err
	}
	defer adotgo.Close()
	if err := adotgoTempl.Execute(adotgo, info); err != nil {
		return err
	}

	return nil
}
