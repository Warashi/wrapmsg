package wrapmsg

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/ssa"
)

const doc = "wrapmsg is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "wrapmsg",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		buildssa.Analyzer,
	},
}

var (
	builtssa  *buildssa.SSA
	inspected *inspector.Inspector
	posMap    map[token.Pos]ast.Node
)

func printIndent(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Print("\t")
	}
}

type walker struct {
	stack []string
}

func (w *walker) push(n string) {
	w.stack = append(w.stack, n)
}

func (w *walker) pop() {
	w.stack = w.stack[:len(w.stack)-1]
}

func (w *walker) contains(n string) bool {
	for _, s := range w.stack {
		if s == n {
			return true
		}
	}
	return false
}

func getIdentName(v interface{ Pos() token.Pos }) []string {
	ident, ok := posMap[v.Pos()].(*ast.Ident)
	switch v := v.(type) {
	case *ssa.Slice:
		if ok {
			fmt.Println("Slice:", ident.Name)
		}
		return nil
	case *ssa.Alloc:
		if ok {
			fmt.Println("Alloc:", ident.Name)
		}
		return nil
	case *ssa.IndexAddr:
		if ok {
			fmt.Println("IndexAddr:", ident.Name)
		}
		return nil
	case *ssa.Store:
		if ok {
			fmt.Println("Store:", ident.Name)
		}
		return nil
	case *ssa.ChangeInterface:
		if ok {
			fmt.Println("ChangeInterface:", ident.Name)
		}
		return nil
	case *ssa.Call:
		if ok {
			fmt.Println("Call:", ident.Name)
		}
		break
	case *ssa.UnOp:
		if ok {
			fmt.Println("UnOp:", ident.Name)
		}
		break
	case *ssa.Parameter:
		fmt.Println("Parameter:", v.Object().Name())
		return []string{v.Object().Name()}
	}

	if !ok {
		return nil
	}

	return []string{ident.Name}
}

func (w *walker) walkValue(depth int, v ssa.Value) ([]string, bool) {
	if w.contains(v.Name()) {
		return nil, false
	}
	w.push(v.Name())
	defer w.pop()

	printIndent(depth)
	fmt.Printf("%[1]v\t%[1]T\n", v)
	org := v
	switch v := v.(type) {
	case *ssa.Function:
	case *ssa.Const:
	case *ssa.Slice:
		for _, v := range GetOperands(v) {
			if r, ok := w.walkValue(depth+1, v); ok {
				return append(r, getIdentName(org)...), true
			}
		}
	case *ssa.Alloc:
		for _, v := range *v.Referrers() {
			if r, ok := w.walkInstructions(depth+1, v); ok {
				return append(r, getIdentName(org)...), true
			}
		}
	case *ssa.IndexAddr:
		if r, ok := w.walkValue(depth+1, v.X); ok {
			return append(r, getIdentName(org)...), true
		}

		//for _, v := range *v.Referrers() {
		//	if r, ok := w.walkInstructions(depth+1, v); ok {
		//		return append(r, getIdentName(org)...), true
		//	}
		//}
	case *ssa.ChangeInterface:
		for _, v := range GetOperands(v) {
			if r, ok := w.walkValue(depth+1, v); ok {
				return append(r, getIdentName(org)...), true
			}
		}
	case *ssa.Call:
		for _, v := range GetOperands(v) {
			if r, ok := w.walkValue(depth+1, v); ok {
				return append(r, getIdentName(org)...), true
			}
		}
	case *ssa.UnOp:
		if r, ok := w.walkValue(depth+1, v.X); ok {
			return append(r, getIdentName(org)...), true
		}
		return getIdentName(v), true
	case *ssa.Parameter:
		return getIdentName(v), true
	}
	return nil, false
}

func (w *walker) walkInstructions(depth int, v ssa.Instruction) ([]string, bool) {
	printIndent(depth)
	fmt.Printf("%[1]v\t%[1]T\n", v)
	org := v
	switch v := v.(type) {
	case *ssa.Slice:
		for _, v := range GetOperands(v) {
			if r, ok := w.walkValue(depth+1, v); ok {
				return append(r, getIdentName(org)...), true
			}
		}
	case *ssa.Alloc:
		for _, v := range *v.Referrers() {
			if r, ok := w.walkInstructions(depth+1, v); ok {
				return append(r, getIdentName(org)...), true
			}
		}
	case *ssa.IndexAddr:
		if r, ok := w.walkValue(depth+1, v.X); ok {
			return append(r, getIdentName(org)...), true
		}

		//for _, v := range *v.Referrers() {
		//	if r, ok := w.walkInstructions(depth+1, v); ok {
		//		return append(r, getIdentName(org)...), true
		//	}
		//}
	case *ssa.Store:
		for _, v := range GetOperands(v) {
			if r, ok := w.walkValue(depth+1, v); ok {
				return append(r, getIdentName(org)...), true
			}
		}
	case *ssa.ChangeInterface:
		for _, v := range GetOperands(v) {
			if r, ok := w.walkValue(depth+1, v); ok {
				return append(r, getIdentName(org)...), true
			}
		}
	case *ssa.Call:
		for _, v := range GetOperands(v) {
			if r, ok := w.walkValue(depth+1, v); ok {
				return append(r, getIdentName(org)...), true
			}
		}
	case *ssa.UnOp:
		if r, ok := w.walkValue(depth+1, v.X); ok {
			return append(r, getIdentName(org)...), true
		}
		return getIdentName(v), true
	}
	return nil, false
}

func buildPosMap() {
	posMap = make(map[token.Pos]ast.Node)
	inspected.Preorder(nil, func(node ast.Node) {
		for i := node.Pos(); i <= node.End(); i++ {
			posMap[i] = node
		}
	})
}

func run(pass *analysis.Pass) (interface{}, error) {
	builtssa = pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	inspected = pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	buildPosMap()
	for _, f := range builtssa.SrcFuncs {
		fmt.Printf("%[1]s\t%[2]T\t%[2]v\n", f.Name(), f)
		for _, b := range f.Blocks {
			fmt.Printf("\t%[1]s\t%[2]T\t%[2]v\n", b.Comment, b)
			for _, instr := range b.Instrs {
				switch v := instr.(type) {
				case *ssa.Call:

					if f := v.Common().StaticCallee(); f == nil || f.Name() == "Errorf" {
						continue
					}

					w := new(walker)
					if r, ok := w.walkInstructions(2, v); ok {
						fmt.Println(strings.Join(r, "."))
						break
					}
				}
			}
		}
	}
	return nil, nil
}
