package src

import (
	"os"

	"github.com/llir/llvm/ir/types"
)

const (
	WARNING  = 0x0
	CRITICAL = 0x1

	TYPE_SYM_TABLE uint32 = 0x0
	TYPE_U8        uint32 = 0x1
	TYPE_I8        uint32 = 0x2
	TYPE_U16       uint32 = 0x3
	TYPE_I16       uint32 = 0x4
	TYPE_U32       uint32 = 0x5
	TYPE_I32       uint32 = 0x6
	TYPE_F32       uint32 = 0x7
	TYPE_U64       uint32 = 0x8
	TYPE_I64       uint32 = 0x9
	TYPE_F64       uint32 = 0xA
	TYPE_BOOL      uint32 = 0xB
	TYPE_STRUCT    uint32 = 0xC
	TYPE_INSTANCE  uint32 = 0xD // Instance of a struct
	TYPE_STRING    uint32 = 0xE
	TYPE_FN        uint32 = 0xF
	TYPE_ANY       uint32 = 0x10
	TYPE_NULL      uint32 = 0x11
)

type File struct {
	Filename string
	Source   *string
}

type TavValue struct {
	Int    int64
	Float  float64
	String []byte
	Bool   bool
	Any    interface{}
}

type TavType struct {
	Type        uint32
	Instance    string // store the identifier of the instance we are referencing
	Indirection int8
	RetType     *TavType // used for function calls
}

func NewTavType(Typ uint32, Instance string, Indirection int8, RetType *TavType) TavType {
	return TavType{
		Type:        Typ,
		Instance:    Instance,
		Indirection: Indirection,
		RetType:     RetType,
	}
}

func (TavType TavType) IsInt() bool {
	return TavType.Type == TYPE_I8 || TavType.Type == TYPE_I16 || TavType.Type == TYPE_I32 || TavType.Type == TYPE_I64
}

func (TavType TavType) IsFloat() bool {
	return TavType.Type == TYPE_F32 || TavType.Type == TYPE_F64
}

type Compiler struct {
	File *File
}

// report an error, the compiler will decide what to do given the severity
func (compiler *Compiler) Report(severity uint8, reporter *Reporter, errCode uint32, msg string) {
	switch severity {
	case WARNING:
		compiler.Warning(reporter, errCode, msg)
	case CRITICAL:
		compiler.Critical(reporter, errCode, msg)
	}
}

// report a warning to the compiler. the compiler will continue and this will not effect the output
func (compiler *Compiler) Warning(reporter *Reporter, errCode uint32, msg string) {
	Log("WARNING", "line", reporter.GetLine(), "pos", reporter.GetIndent(), "in", compiler.File.Filename, "\n")
	Log(msg)
	Log("\n")
}

// report a critical error to the compiler. the compiler will exit from this point as it cannot continue
func (compiler *Compiler) Critical(reporter *Reporter, errCode uint32, msg string) {
	Log("CRITICAL ERROR", reporter.GetLine(), ":", reporter.GetIndent(), "in", compiler.File.Filename, "\n")
	Log(reporter.ReportLine())
	Log(msg)
	Log("\n")
	os.Exit(2)
}

func ConvertType(tavType TavType, SymTable *SymTable) types.Type {
	switch tavType.Type {
	case TYPE_BOOL:
		if tavType.Indirection > 0 {
			return types.I1Ptr
		}
		return types.I1
	case TYPE_I8:
		if tavType.Indirection > 0 {
			return types.I8Ptr
		}
		return types.I8
	case TYPE_I16:
		if tavType.Indirection > 0 {
			return types.I16Ptr
		}
		return types.I16
	case TYPE_I32:
		if tavType.Indirection > 0 {
			return types.I32Ptr
		}
		return types.I32
	case TYPE_I64:
		if tavType.Indirection > 0 {
			return types.I64Ptr
		}
		return types.I64
	case TYPE_F32:
		if tavType.Indirection > 0 {
			// NOT SURE HOW THIS WORKS
		}
		return types.Float
	case TYPE_F64:
		if tavType.Indirection > 0 {
			// NOT SURE HOW THIS WORKS
		}
		return types.Double
	case TYPE_STRING:
		return types.I8Ptr
	case TYPE_INSTANCE:
		v := SymTable.Get(tavType.Instance).Value
		return v.(types.Type)
	}
	return types.Void
}

func InvertPtrType(tavType TavType, direction int8) TavType {
	return TavType{
		Type:        tavType.Type,
		Indirection: tavType.Indirection + direction,
		RetType:     tavType.RetType,
	}
}

// TODO some type of check as to whether the inference join was valid
// infer the type of an expression
// this may be somewhat recursive as we have to infer sub types
// if there are multiple types, we do an inference join (infer the correct type given multiple)
func InferType(expression AST, SymTable *SymTable) TavType {
	switch e := expression.(type) {
	case *VariableAST:
		t := SymTable.Get(e.Identifier.Lexme())
		if t != nil {
			return t.Type
		}
		break
	case *UnaryAST:
		switch e.Operator.Type {
		case ADDR:
			// here we need to cast the type to a pointer
			t := InferType(e.Right, SymTable)
			// increase the indirection count
			t.Indirection += 1
			return t
		case STAR:
			// here we need to cast the type to a pointer
			t := InferType(e.Right, SymTable)
			// increase the indirection count
			t.Indirection -= 1
			return t
		}
	case *LiteralAST:
		return e.Type
	case *ReturnAST:
		return InferType(e.Value, SymTable)
	case *BinaryAST:
		return JoinInfered(InferType(e.Left, SymTable), InferType(e.Right, SymTable))
	case *CallAST:
		t := InferType(e.Caller, SymTable)
		return *t.RetType
	case *StructGetAST:
		// get the name of the struct that we are referencing
		s := InferType(e.Struct, SymTable).Instance
		// now we know that the struct stores its members in a symbol table entry called
		// structname_members
		structSymName := s + "_members"
		// get the symbol table for the struct
		sym := SymTable.Get(structSymName)
		t := sym.Value.(*Scope)
		if t != nil {
			// get the structs member symbol table and access the member
			// finally get the type of the member in that symtable
			return t.Get(e.Member.Lexme()).Type
		}
		break
	case *VarDefAST:
		return e.Type
	case *CastAST:
		return e.TavType
	}
	// this is unreachable (in theory)
	return TavType{}
}

// TODO some way to cast the type if they can be joined
// join 2 infered types and figure out what the next type will be
func JoinInfered(type1, type2 TavType) TavType {
	if type1.Type == type2.Type {
		return type1
	}
	// deal with number types
	if type1.IsInt() && type2.IsInt() {
		return type1
	} else if type1.IsFloat() && type2.IsFloat() {
		return type1
	} else if type1.IsFloat() && type2.IsInt() {
		return type1
	} else if type1.IsInt() && type2.IsFloat() {
		return type2
	}
	Assert(false, "cant join infered types")
	return type1
}

// check if 2 types are able to be cast
func Compatible(t1, t2 TavType) bool {
	// TODO fix this so it accomodates strings as they are pointers
	//// 2 types are immediately not compatible if they have different pointer indirections
	//if t1.Indirection != t2.Indirection{
	//	return false
	//}
	// check integers
	if (t1.Type == TYPE_I8 || t1.Type == TYPE_I16 || t1.Type == TYPE_I32 || t1.Type == TYPE_I64) &&
		(t2.Type == TYPE_I8 || t2.Type == TYPE_I16 || t2.Type == TYPE_I32 || t2.Type == TYPE_I64) {
		return true
		// check floats
	} else if (t1.Type == TYPE_F32 || t1.Type == TYPE_F64) && (t2.Type == TYPE_F32 || t2.Type == TYPE_F64) {
		return true
	} else if (t1.Type == TYPE_I8 || t1.Type == TYPE_I16 || t1.Type == TYPE_I32 || t1.Type == TYPE_I64) &&
		(t2.Type == TYPE_F32 || t2.Type == TYPE_F64) {
		return true
	} else if (t1.Type == TYPE_F32 || t1.Type == TYPE_F64) &&
		(t2.Type == TYPE_I8 || t2.Type == TYPE_I16 || t2.Type == TYPE_I32 || t2.Type == TYPE_I64) {
		return true
	}

	return false
}

// cast an expression to a certian type
// TODO actually check if the cast is valid, this is very hacky
// TODO also this probably shouldn't modify the expression, it may lead to some errors in code gen
// TODO also make this work with other variables
// as this is called in the checker
func Cast(tavType TavType, expression AST) bool {
	//switch e:=expression.(type){
	//case *LiteralAST:
	//	if Compatible(tavType, e.Type){
	//		e.Type = tavType
	//
	//		return true
	//	}
	//}
	return false
}
