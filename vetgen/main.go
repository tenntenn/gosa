package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var g VetGen
	flag.StringVar(&g.ToolName, "name", "", "vet tool name")
	flag.Parse()
	g.ExeName = os.Args[0]
	g.Args = flag.Args()

	if err := g.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
