package opcodes

const (
	Push byte = iota
	Pop
	IdentifierAssign
	SelectorAssign
	Label
	Jump
	IfJump
	Return
	DeleteIdentifier
	DeleteSelector
	Defer
	NewFunction
	NewClass
	Call
	NewArray
	NewTuple
	NewHash
	Identifier
	Integer
	Float
	String
	Bytes
	True
	False
	None
	Selector
	Super
)

var OpCodes = map[byte]string{
	Push:             "Push",
	Pop:              "Pop",
	IdentifierAssign: "IdentifierAssign",
	SelectorAssign:   "SelectorAssign",
	Label:            "Label",
	Jump:             "Jump",
	IfJump:           "IfJump",
	Return:           "Return",
	DeleteIdentifier: "DeleteIdentifier",
	DeleteSelector:   "DeleteSelector",
	Defer:            "Defer",
	NewFunction:      "NewFunction",
	NewClass:         "NewClass",
	Call:             "Call",
	NewArray:         "NewArray",
	NewTuple:         "NewTuple",
	NewHash:          "NewHash",
	Identifier:       "Identifier",
	Integer:          "Integer",
	Float:            "Float",
	String:           "String",
	Bytes:            "Bytes",
	True:             "True",
	False:            "False",
	None:             "None",
	Selector:         "Selector",
	Super:            "Super",
}
