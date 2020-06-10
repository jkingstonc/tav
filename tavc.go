package main

import (
	"os"
	"tav/src"
)

const (
	// build and emit executable
	BUILD = 0x0
)

func main() {
	src.Log("tav v_a_0_1")
	args := os.Args[1:]
	if args[0] == "build" {
		build()
	}
}

func build() {
	compiler := &src.Compiler{}
	source := src.Preprocess(compiler, "//this is a comment\nhello world")
	tokens := src.Lex(compiler, source)
	ast := src.Parse(compiler, tokens)
	optimized := src.Optimize(compiler, ast)
	src.Generate(compiler, optimized)
}
