package src

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type FuncEvaluator struct {
	Func *ir.Func
	Args []value.Value
}

// JIT compile using LLVM and return the result
func JITCompile(filename string, source *string) (value.Value, uint8) {
	// first compile down to a LLVM module
	module := AheadCompile(filename, source)
	// then execute the module
	result, err := JIT(module)
	return result, err
}

// Compile to LLVM and JIT compile
func JIT(module *ir.Module) (value.Value, uint8) {
	return constant.NewNull(types.I1Ptr), SUCCESS_JIT
}

func NewEvaluator(Func *ir.Func, Args ...value.Value) *FuncEvaluator{
	return &FuncEvaluator{
		Func: Func,
		Args: Args,
	}
}

func (evaluator *FuncEvaluator) Eval() value.Value {
	return constant.NewNull(types.I1Ptr)
}

