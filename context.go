package wrapmsg

import (
	"context"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ast/inspector"
)

type contextKey int

const (
	_ contextKey = iota
	passKey
	ssaKey
	inspectorKey
	posMapKey
)

func getPass(ctx context.Context) *analysis.Pass {
	return ctx.Value(passKey).(*analysis.Pass)
}

func getSSA(ctx context.Context) *buildssa.SSA {
	return ctx.Value(ssaKey).(*buildssa.SSA)
}

func getInspector(ctx context.Context) *inspector.Inspector {
	return ctx.Value(inspectorKey).(*inspector.Inspector)
}

func getPosMap(ctx context.Context) map[token.Pos][]ast.Node {
	return ctx.Value(posMapKey).(map[token.Pos][]ast.Node)
}
