package codec

import "github.com/maketaio/openapi/runtime/fields"

type Code int

const (
	CodeTypeMismatch Code = iota
	CodeMissingField
	CodeUnknownField
)

type ValueKind int

const (
	KindString ValueKind = iota
	KindNumber
	KindInteger
	KindBoolean
	KindObject
	KindArray
	KindNull
	KindAny
)

type Issue struct {
	Path     fields.Path
	Code     Code
	Expected ValueKind // for type mismatches
	Actual   ValueKind // for type mismatches
	Allowed  []string  // for unknown field
	Message  string
}

func (i *Issue) Error() string {
	return i.Message
}
