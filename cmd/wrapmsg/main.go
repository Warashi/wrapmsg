package main

import (
	"wrapmsg"

	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(wrapmsg.Analyzer) }
