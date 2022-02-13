package main

import (
	"github.com/Warashi/wrapmsg"

	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(wrapmsg.Analyzer) }
