package src

import "hash/fnv"

const (
	SYM_MACRO  uint8 = 0x0
	SYM_STRUCT uint8 = 0x1

	// the symbol is private to the file it is declared in
	ATTRIB_PRIVATE uint8 = 0x1 << 0
	// the symbol is exposed to other files
	ATTRIB_EXPOSED uint8 = 0x1 << 1
	// the symbol is runnable at compile time (used for functions)
	ATTRIB_DOABLE uint8 = 0x1 << 2
)

// a symbol is identified by a type and an attribute
type Symbol struct {
	Type        uint8
	Attribuites uint8
}

// keep a record of symbol identifiers along with their type and attribute
type SymTable struct {
	Symbols map[uint32]Symbol
}

func (symTable *SymTable) Add(identifier string, symType, attributes uint8) {
	h := hash(identifier)
	_, ok := symTable.Symbols[h]
	Assert(!ok, "symbol already exists in symbol table", identifier)
	symTable.Symbols[h] = Symbol{Type: symType, Attribuites: attributes}
}

func (symTable *SymTable) Get(identifier string) Symbol {
	h := hash(identifier)
	sym, ok := symTable.Symbols[h]
	Assert(ok, "symbol couldn't be found in the symbol table", identifier)
	return sym
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
