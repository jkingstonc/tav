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
	Identifier  string
	Type        TavType
	Attribuites uint8
	Value       interface{}		// used for value checks etc
}

// keep a record of symbol identifiers along with their type and attribute
type SymTable struct {
	// used so we can keep track of the scope of variables
	Parent   *SymTable
	Symbols  []*Symbol
}

// create a new symbol table
func NewSymTable(parent *SymTable) *SymTable {
	return &SymTable{
		Parent: parent,
	}
}

// enter a new scope in the symbol table
func (symTable *SymTable) NewScope() *SymTable{

	// create a new table
	newTable := NewSymTable(symTable)
	// add it to the current scope
	symTable.Add("", TavType{
		Type:   TYPE_SYM_TABLE,
	},0, newTable)
	// return the new table
	return newTable
}

// return from the scope in the symbol table
func (symTable *SymTable) PopScope() *SymTable{
	parent := symTable.Parent
	parent.RemoveByValue(symTable)
	return parent
}

// add a symbol to the table and retrieve the integer id
func (symTable *SymTable) Add(identifier string, symType TavType, attributes uint8, value interface{}) {
	symTable.Symbols = append(symTable.Symbols, &Symbol{
		Identifier:  identifier,
		Type:        symType,
		Attribuites: attributes,
		Value:       value,
	})
}


// get the symbol value given an id
func (symTable *SymTable) GetLocal(identifier string) *Symbol {
	for _, sym := range symTable.Symbols{
		if sym.Identifier == identifier {
			return sym
		}
	}
	return nil
}

// get the symbol value given an id
func (symTable *SymTable) Get(identifier string) *Symbol {
	for _, sym := range symTable.Symbols{
		if sym.Identifier == identifier {
			return sym
		}
	}
	if symTable.Parent != nil {
		return symTable.Parent.Get(identifier)
	}
	return nil
}

// get the symbol value given an id
func (symTable *SymTable) RemoveByID(identifier string){
	for i, sym := range symTable.Symbols{
		if sym.Identifier == identifier{
			symTable.Symbols = append(symTable.Symbols[:i], symTable.Symbols[i+1:]...)
			return
		}
	}
	Assert(true, "symbol couldn't be removed (doesn't exist!)")
}

// get the symbol value given an id
func (symTable *SymTable) RemoveByValue(value interface{}){
	for i, sym := range symTable.Symbols{
		if sym.Value == value{
			symTable.Symbols = append(symTable.Symbols[:i], symTable.Symbols[i+1:]...)
			return
		}
	}
	Assert(true, "symbol couldn't be removed (doesn't exist!)")
}


// hash function to hash a string
func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
