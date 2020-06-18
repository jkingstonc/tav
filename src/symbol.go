package src

const (
	// the symbol is private to the file it is declared in
	ATTRIB_PRIVATE uint8 = 0x1 << 0
	// the symbol is exposed to other files
	ATTRIB_EXPOSED uint8 = 0x1 << 1
	// the symbol is runnable at compile time (used for functions)
	ATTRIB_DOABLE uint8 = 0x1 << 2
)

// a symbol is identified by a type and an attribute
type Symbol struct {
	Identifier string
	Type       TavType
	Value      interface{} // used for value checks etc
}

type Scope struct {
	Identifier string
	// Store a reference to the parent so we can look up scopes for variable declerations
	Parent  *Scope
	Symbols []*Symbol
}

// keep a record of symbol identifiers along with their type and attribute
type SymTable struct {
	CurrentScope *Scope
}

func NewSym(Identifier string, Type TavType, Value interface{}) *Symbol {
	return &Symbol{
		Identifier: Identifier,
		Type:       Type,
		Value:      Value,
	}
}

func NewScope(parent *Scope, identifier string) *Scope {
	return &Scope{
		Identifier: identifier,
		Parent:     parent,
		Symbols:    nil,
	}
}

func (Scope *Scope) Add(symbol *Symbol) {
	Scope.Symbols = append(Scope.Symbols, symbol)
}

func (Scope *Scope) Get(identifier string) *Symbol {
	for _, sym := range Scope.Symbols {
		if sym.Identifier == identifier {
			return sym
		}
	}
	if Scope.Parent != nil {
		return Scope.Parent.Get(identifier)
	}
	return nil
}

func NewSymTable() *SymTable {
	return &SymTable{CurrentScope: NewScope(nil, "root")}
}

// enter a new scope in the symbol table
func (SymTable *SymTable) NewScope(Identifier string) {
	// create a new scope
	newScope := NewScope(SymTable.CurrentScope, Identifier)
	// add the scope to the current symbol table
	newScopeSym := NewSym(Identifier, NewTavType(TYPE_SYM_TABLE, "", 0, nil), newScope)
	SymTable.CurrentScope.Add(newScopeSym)
	// update the current scope
	SymTable.CurrentScope = newScope
}

// return from the scope in the symbol table
func (SymTable *SymTable) PopScope() {
	SymTable.CurrentScope = SymTable.CurrentScope.Parent
}

// add a symbol to the table and retrieve the integer id
func (SymTable *SymTable) Add(identifier string, tavType TavType, value interface{}) {
	SymTable.CurrentScope.Add(NewSym(identifier, tavType, value))
}

// get the symbol value given an id
func (SymTable *SymTable) GetLocal(identifier string) *Symbol {
	for _, sym := range SymTable.CurrentScope.Symbols {
		if sym.Identifier == identifier {
			return sym
		}
	}
	return nil
}

// get the symbol value given an id
func (SymTable *SymTable) Get(identifier string) *Symbol {
	return SymTable.CurrentScope.Get(identifier)
}
