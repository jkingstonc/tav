package src

import "time"

func Generate(compiler *Compiler, ast *AST) uint8 {
	start := time.Now()

	LLVMCompile(compiler, ast)

	end := time.Since(start)
	Log("back end took ", end.Seconds(), "ms")
	return SUCCESS_COMP
}
