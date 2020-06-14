package src

import (
	"github.com/llir/llvm/ir"
	"time"
)

// implements visitor
type Generator struct {
	Root   AST
	Module *ir.Module
}

func Generate(compiler *Compiler, ast AST) uint8 {
	start := time.Now()
	//LLVMCompile(compiler, ast)

	// create a new LLVM module
	module := ir.NewModule()
	generator := &Generator{
		Root: ast,
		Module: module,
	}

	result := generator.Run()

	end := time.Since(start)
	Log("back end took ", end.Seconds(), "seconds")
	return result
}

func (generator *Generator) Run() uint8 {



	return SUCCESS_COMP
}