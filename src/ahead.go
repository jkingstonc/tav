package src

import (
	"github.com/llir/llvm/ir"
	"time"
)

// compile to exe
func AheadCompile(filename string, source *string) *ir.Module {
	start := time.Now()
	compiler := &Compiler{FileName: filename, Source: source}
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