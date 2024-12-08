package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/jtbonhomme/go-nilcheck"
)

func main() {
	singlechecker.Main(nilcheck.Analyzer)
}
