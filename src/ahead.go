package src

import (
	"time"

	"github.com/llir/llvm/ir"
)

// compile to exe
func AheadCompile(File *File) *ir.Module {
	start := time.Now()
	compiler := &Compiler{File: File}
	tokens := Lex(compiler)
	tokens = ProcessDirectives(compiler, tokens)
	ast := Parse(compiler, tokens)
	ast = Check(compiler, ast)
	optimized := Optimize(compiler, ast)
	result := Generate(compiler, optimized)
	end := time.Since(start)
	Log("compilation took ", end.Seconds(), "seconds")
	return result
}
