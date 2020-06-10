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
	tokens := src.Lex(compiler)
	ast := src.Parse(tokens)
	optimized := src.Optimize(ast)
	src.Generate(optimized)

}
