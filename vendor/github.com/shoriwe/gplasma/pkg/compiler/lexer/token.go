package lexer

/*
	Token Kinds
*/

const (
	Unknown uint8 = iota
	PendingEscape
	Comment
	Whitespace
	Literal
	Tab
	IdentifierKind
	JunkKind
	Separator
	Punctuation
	Assignment
	Comparator
	Operator
	SingleQuoteString
	DoubleQuoteString
	Integer
	HexadecimalInteger
	BinaryInteger
	OctalInteger
	Float
	ScientificFloat
	CommandOutput
	ByteString
	Keyboard
	Boolean
	NoneType
	EOF

	Comma
	Colon
	SemiColon
	NewLine

	Pass
	Super
	End
	If
	Unless
	As
	Raise
	Else
	Elif
	While
	Do
	For
	Until
	Switch
	Case
	Default
	Yield
	Return
	Continue
	Break
	Redo
	Module
	Def
	Lambda
	Interface
	Class
	Try
	Except
	Finally
	BEGIN
	END
	Context

	Assign
	NegateBitsAssign
	BitwiseOrAssign
	BitwiseXorAssign
	BitWiseAndAssign
	BitwiseLeftAssign
	BitwiseRightAssign
	AddAssign
	SubAssign
	StarAssign
	DivAssign
	FloorDivAssign
	ModulusAssign
	PowerOfAssign

	Not
	SignNot
	NegateBits

	And
	Or
	Xor
	In

	Equals
	NotEqual
	GreaterThan
	GreaterOrEqualThan
	LessThan
	LessOrEqualThan

	BitwiseOr
	BitwiseXor
	BitWiseAnd
	BitwiseLeft
	BitwiseRight

	Add // This is also an unary operator
	Sub // This is also an unary operator
	Star
	Div
	FloorDiv
	Modulus
	PowerOf

	True
	False
	None

	OpenParentheses
	CloseParentheses
	OpenSquareBracket
	CloseSquareBracket
	OpenBrace
	CloseBrace
	DollarSign
	Dot
)

/*
	Regular Expressions
*/

type Token struct {
	String      string
	DirectValue uint8
	Kind        uint8
	Line        int
	Column      int
	Index       int
}

/*
	Keyboards
*/

var (
	CommaChar              = ','
	ColonChar              = ':'
	SemiColonChar          = ';'
	NewLineChar            = '\n'
	OpenParenthesesChar    = '('
	CloseParenthesesChar   = ')'
	OpenSquareBracketChar  = '['
	CloseSquareBracketChar = ']'
	OpenBraceChar          = '{'
	CloseBraceChar         = '}'
	DollarSignChar         = '$'
	DotChar                = '.'
	BitwiseOrChar          = '|'
	BitwiseXorChar         = '^'
	BitWiseAndChar         = '&'
	AddChar                = '+'
	SubChar                = '-'
	StarChar               = '*'
	DivChar                = '/'
	ModulusChar            = '%'
	LessThanChar           = '<'
	GreatThanChar          = '>'
	NegateBitsChar         = '~'
	SignNotChar            = '!'
	EqualsChar             = '='
	WhiteSpaceChar         = ' '
	TabChar                = '\t'
	CommentChar            = '#'
	BackSlashChar          = '\\'

	PassString      = "pass"
	SuperString     = "super"
	EndString       = "end"
	IfString        = "if"
	UnlessString    = "unless"
	ElseString      = "else"
	ElifString      = "elif"
	WhileString     = "while"
	DoString        = "do"
	ForString       = "for"
	UntilString     = "until"
	SwitchString    = "switch"
	CaseString      = "case"
	DefaultString   = "default"
	YieldString     = "yield"
	ReturnString    = "return"
	ContinueString  = "continue"
	BreakString     = "break"
	RedoString      = "redo"
	ModuleString    = "module"
	DefString       = "def"
	LambdaString    = "lambda"
	InterfaceString = "interface"
	ClassString     = "class"
	TryString       = "try"
	ExceptString    = "except"
	FinallyString   = "finally"
	AndString       = "and"
	OrString        = "or"
	XorString       = "xor"
	InString        = "in"
	BEGINString     = "BEGIN"
	ENDString       = "END"
	NotString       = "not"
	TrueString      = "True"
	FalseString     = "False"
	NoneString      = "None"
	ContextString   = "context"
	RaiseString     = "raise"
	AsString        = "as"
)
