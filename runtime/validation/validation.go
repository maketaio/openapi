package validation

import "strings"

type Code string

const (
	CodeIntGT         Code = "int.gt"
	CodeIntGTE        Code = "int.gte"
	CodeIntLT         Code = "int.lt"
	CodeIntLTE        Code = "int.lte"
	CodeIntMultipleOf Code = "int.multipleOf"
	CodeNumGT         Code = "num.gt"
	CodeNumGTE        Code = "num.gte"
	CodeNumLT         Code = "num.lt"
	CodeNumLTE        Code = "num.lte"
)

type Path []string

func (p Path) String() string {
	return strings.Join(p, ".")
}

func (p Path) Field(s string) Path {
	return append(p, s)
}

type Error struct {
	Path    Path   `json:"path"`
	Detail  Detail `json:"detail"`
	Message string `json:"message"`
}

type Detail struct {
	Code          Code     `json:"code"`
	IntGT         *int64   `json:"int.gt,omitzero"`
	IntGTE        *int64   `json:"int.gte,omitzero"`
	IntLT         *int64   `json:"int.lt,omitzero"`
	IntLTE        *int64   `json:"int.lte,omitzero"`
	IntMultipleOf *int64   `json:"int.multipleOf,omitzero"`
	NumGT         *float64 `json:"num.gt,omitzero"`
	NumGTE        *float64 `json:"num.gte,omitzero"`
	NumLT         *float64 `json:"num.lt,omitzero"`
	NumLTE        *float64 `json:"num.lte,omitzero"`
}

func (e *Error) Error() string {
	return e.Message
}
