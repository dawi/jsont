package jsont

type TokenType int

const (
	Empty = iota
	Whitespace
	ObjectStart
	ObjectEnd
	ArrayStart
	ArrayEnd
	Colon
	Comma
	FieldName
	True
	False
	Null
	Integer
	Float
	String
	Unknown
)

type Token struct {
	Type  TokenType
	Value string
}
