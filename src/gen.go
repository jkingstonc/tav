package src

import "time"

// implements visitor
type Generator struct {
	Root   AST
}

func Generate(compiler *Compiler, ast AST) uint8 {
	start := time.Now()

	LLVMCompile(compiler, ast)

	end := time.Since(start)
	Log("back end took ", end.Seconds(), "seconds")
	return SUCCESS_COMP
}
