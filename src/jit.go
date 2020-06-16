package src

import (
	"time"
)

// JIT compile by walking the AST
// Return the value as a go value (will need to be casted if doing compile time JIT)
func JITCompile(File *File) interface{} {
	start := time.Now()
	compiler := &Compiler{File: File}
	tokens := Lex(compiler)
	tokens = ProcessDirectives(compiler, tokens)
	ast := Parse(compiler, tokens)
	ast = Check(compiler, ast)
	optimized := Optimize(compiler, ast)
	result := Interpret(compiler, optimized)
	end := time.Since(start)
	Log("compilation took ", end.Seconds(), "seconds")
	return result
}
