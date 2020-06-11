package src

const (	
	// emit an exe
	AHEAD_COMPILE uint8 = 0x0
	// JIT compile using LLVM execution engine (used for compile time execution)
	JIT_COMPILE   uint8 = 0x1

	SUCCESS_COMP  uint8 = 0x0
	FAIL_COMP     uint8 = 0x1
)

// compile to exe
func AheadCompile(source string) uint8{
	compiler := &Compiler{}
	modSrc := Preprocess(compiler, source)
	tokens := Lex(compiler, modSrc)
	ast := Parse(compiler, tokens)
	optimized := Optimize(compiler, ast)
	return Generate(compiler, optimized)
}

// JIT compile using LLVM and return the result
func JITCompile(source string) (uint8, string) {
	compiler := &Compiler{}
	modSrc := Preprocess(compiler, source)
	tokens := Lex(compiler, modSrc)
	ast := Parse(compiler, tokens)
	optimized := Optimize(compiler, ast)
	return JIT(compiler, optimized)
}