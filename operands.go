package wrapmsg

import (
	"golang.org/x/tools/go/ssa"
)

type operander interface {
	Operands([]*ssa.Value) []*ssa.Value
}

func getOperands(o operander) []ssa.Value {
	ops := o.Operands(nil)
	ret := make([]ssa.Value, 0, len(ops))
	for i := range ops {
		if ops[i] != nil && *ops[i] != nil {
			ret = append(ret, *ops[i])
		}
	}
	return ret
}
