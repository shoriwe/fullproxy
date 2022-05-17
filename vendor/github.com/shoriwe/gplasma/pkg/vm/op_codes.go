package vm

const (
	ReturnState uint8 = iota
	BreakState
	RedoState
	ContinueState
	NoState
)

const (
	NewStringOP uint8 = iota
	NewBytesOP
	NewIntegerOP
	NewFloatOP
	GetTrueOP
	GetFalseOP
	GetNoneOP
	NewTupleOP
	NewArrayOP
	NewHashOP

	UnaryOP
	NegateBitsOP
	BoolNegateOP
	NegativeOP

	BinaryOP
	AddOP
	SubOP
	MulOP
	DivOP
	FloorDivOP
	ModOP
	PowOP
	BitXorOP
	BitAndOP
	BitOrOP
	BitLeftOP
	BitRightOP
	AndOP
	OrOP
	XorOP
	EqualsOP
	NotEqualsOP
	GreaterThanOP
	LessThanOP
	GreaterThanOrEqualOP
	LessThanOrEqualOP
	ContainsOP

	IfOP
	IfOneLinerOP
	UnlessOP
	UnlessOneLinerOP

	GetIdentifierOP
	IndexOP
	SelectNameFromObjectOP
	MethodInvocationOP
	AssignIdentifierOP
	AssignSelectorOP
	AssignIndexOP
	BreakOP
	RedoOP
	ContinueOP
	ReturnOP

	ForLoopOP
	WhileLoopOP
	DoWhileLoopOP
	UntilLoopOP

	NOP

	SwitchOP

	LoadFunctionArgumentsOP
	NewFunctionOP
	PushOP
	NewGeneratorOP
	TryOP
	NewModuleOP
	NewClassOP
	NewClassFunctionOP
	RaiseOP
)

var unaryInstructionsFunctions = map[uint8]string{
	NegateBitsOP: NegateBits,
	BoolNegateOP: Negate,
	NegativeOP:   Negative,
}

var binaryInstructionsFunctions = map[uint8][2]string{
	AddOP:                {Add, RightAdd},
	SubOP:                {Sub, RightSub},
	MulOP:                {Mul, RightMul},
	DivOP:                {Div, RightDiv},
	FloorDivOP:           {FloorDiv, RightFloorDiv},
	ModOP:                {Mod, RightMod},
	PowOP:                {Pow, RightPow},
	BitXorOP:             {BitXor, RightBitXor},
	BitAndOP:             {BitAnd, RightBitAnd},
	BitOrOP:              {BitOr, RightBitOr},
	BitLeftOP:            {BitLeft, RightBitLeft},
	BitRightOP:           {BitRight, RightBitRight},
	AndOP:                {And, RightAnd},
	OrOP:                 {Or, RightOr},
	XorOP:                {Xor, RightXor},
	EqualsOP:             {Equals, RightEquals},
	NotEqualsOP:          {NotEquals, RightNotEquals},
	GreaterThanOP:        {GreaterThan, RightGreaterThan},
	LessThanOP:           {LessThan, RightLessThan},
	GreaterThanOrEqualOP: {GreaterThanOrEqual, RightGreaterThanOrEqual},
	LessThanOrEqualOP:    {LessThanOrEqual, RightLessThanOrEqual},
	ContainsOP:           {"029p3847980479087437891734", Contains},
}

var instructionNames = map[uint8]string{
	NewStringOP:  "NewStringOP",
	NewBytesOP:   "NewBytesOP",
	NewIntegerOP: "NewIntegerOP",
	NewFloatOP:   "NewFloatOP",
	GetTrueOP:    "GetTrueOP",
	GetFalseOP:   "GetFalseOP",
	GetNoneOP:    "GetNoneOP",
	NewTupleOP:   "NewTupleOP",
	NewArrayOP:   "NewArrayOP",
	NewHashOP:    "NewHashOP",

	UnaryOP:      "UnaryOP",
	NegateBitsOP: "NegateBitsOP",
	BoolNegateOP: "BoolNegateOP",
	NegativeOP:   "NegativeOP",

	BinaryOP:             "BinaryOP",
	AddOP:                "AddOP",
	SubOP:                "SubOP",
	MulOP:                "MulOP",
	DivOP:                "DivOP",
	FloorDivOP:           "FloorDivOP",
	ModOP:                "ModOP",
	PowOP:                "PowOP",
	BitXorOP:             "BitXorOP",
	BitAndOP:             "BitAndOP",
	BitOrOP:              "BitOrOP",
	BitLeftOP:            "BitLeftOP",
	BitRightOP:           "BitRightOP",
	AndOP:                "AndOP",
	OrOP:                 "OrOP",
	XorOP:                "XorOP",
	EqualsOP:             "EqualsOP",
	NotEqualsOP:          "NotEqualsOP",
	GreaterThanOP:        "GreaterThanOP",
	LessThanOP:           "LessThanOP",
	GreaterThanOrEqualOP: "GreaterThanOrEqualOP",
	LessThanOrEqualOP:    "LessThanOrEqualOP",
	ContainsOP:           "ContainsOP",

	IfOP:             "IfOP",
	IfOneLinerOP:     "IfOneLinerOP",
	UnlessOP:         "UnlessOP",
	UnlessOneLinerOP: "UnlessOneLinerOP",

	GetIdentifierOP:        "GetIdentifierOP",
	IndexOP:                "IndexOP",
	SelectNameFromObjectOP: "SelectNameFromObjectOP",
	MethodInvocationOP:     "MethodInvocationOP",
	AssignIdentifierOP:     "AssignIdentifierOP",
	AssignSelectorOP:       "AssignSelectorOP",
	AssignIndexOP:          "AssignIndexOP",
	BreakOP:                "BreakOP",
	RedoOP:                 "RedoOP",
	ContinueOP:             "ContinueOP",
	ReturnOP:               "ReturnOP",

	ForLoopOP:     "ForLoopOP",
	WhileLoopOP:   "WhileLoopOP",
	DoWhileLoopOP: "DoWhileLoopOP",
	UntilLoopOP:   "UntilLoopOP",

	NOP: "NOP",

	SwitchOP: "SwitchOP",

	LoadFunctionArgumentsOP: "LoadFunctionArgumentsOP",
	NewFunctionOP:           "NewFunctionOP",
	PushOP:                  "PushOP",
	NewGeneratorOP:          "NewGeneratorOP",
	TryOP:                   "TryOP",
	NewModuleOP:             "NewModuleOP",
	NewClassOP:              "NewClassOP",
	NewClassFunctionOP:      "NewClassFunctionOP",
	RaiseOP:                 "RaiseOP",
}
