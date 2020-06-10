package src

type Preprocessor struct {
	SymTable *SymTable
}

func Preprocess(compiler *Compiler, source string) string {
	preprocessor := &Preprocessor{
		SymTable: &SymTable{
			Symbols: make(map[uint32]Symbol),
		},
	}
	preprocessor.SymTable.Add("test", SYM_STRUCT, 0)
	return source
}
