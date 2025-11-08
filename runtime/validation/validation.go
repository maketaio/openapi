package validation

import (
	"fmt"

	"github.com/maketaio/openapi/runtime/fields"
)

type Code int

const (
	CodeIntMax Code = iota
	CodeIntMin
	CodeIntExclMax
	CodeIntExclMin
	CodeIntMultipleOf
	CodeNumMax
	CodeNumMin
	CodeNumExclMax
	CodeNumExclMin
	CodeStrMinLen
	CodeStrMaxLen
	CodeStrLen
	CodeStrPattern
	CodeArrMinItems
	CodeArrMaxItems
	CodeArrLen
	CodeObjMinProps
	CodeObjMaxProps
	CodeObjLen
)

type Issue struct {
	Path    fields.Path
	Code    Code
	Params  Params
	Message string
}

func (i *Issue) Error() string {
	return i.Message
}

type Params struct {
	IntMax        *int64
	IntMin        *int64
	IntExclMax    *int64
	IntExclMin    *int64
	IntMultipleOf *int64
	NumMax        *float64
	NumMin        *float64
	NumExclMax    *float64
	NumExclMin    *float64
	StrMinLen     *int64
	StrMaxLen     *int64
	StrLen        *int64
	StrPattern    string
	ArrMinItems   *int64
	ArrMaxItems   *int64
	ArrLen        *int64
	ObjMinProps   *int64
	ObjMaxProps   *int64
	ObjLen        *int64
}

func NewIntMaxIssue(path fields.Path, max int64) *Issue {
	return &Issue{
		Path: path,
		Code: CodeIntMax,
		Params: Params{
			IntMax: &max,
		},
		Message: fmt.Sprintf("%s must be less than or equal to %d", path, max),
	}
}

func NewIntMinIssue(path fields.Path, min int64) *Issue {
	return &Issue{
		Path: path,
		Code: CodeIntMin,
		Params: Params{
			IntMin: &min,
		},
		Message: fmt.Sprintf("%s must be greater than or equal to %d", path, min),
	}
}

func NewIntExclMaxIssue(path fields.Path, max int64) *Issue {
	return &Issue{
		Path: path,
		Code: CodeIntExclMax,
		Params: Params{
			IntExclMax: &max,
		},
		Message: fmt.Sprintf("%s must be less than %d", path, max),
	}
}

func NewIntExclMinIssue(path fields.Path, min int64) *Issue {
	return &Issue{
		Path: path,
		Code: CodeIntExclMin,
		Params: Params{
			IntExclMin: &min,
		},
		Message: fmt.Sprintf("%s must be greater than %d", path, min),
	}
}

func NewNumMaxIssue(path fields.Path, max float64) *Issue {
	return &Issue{
		Path: path,
		Code: CodeNumMax,
		Params: Params{
			NumMax: &max,
		},
		Message: fmt.Sprintf("%s must be less than or equal to %f", path, max),
	}
}

func NewNumMinIssue(path fields.Path, min float64) *Issue {
	return &Issue{
		Path: path,
		Code: CodeNumMin,
		Params: Params{
			NumMin: &min,
		},
		Message: fmt.Sprintf("%s must be greater than or equal to %f", path, min),
	}
}

func NewNumExclMaxIssue(path fields.Path, max float64) *Issue {
	return &Issue{
		Path: path,
		Code: CodeNumExclMax,
		Params: Params{
			NumExclMax: &max,
		},
		Message: fmt.Sprintf("%s must be less than %f", path, max),
	}
}

func NewNumExclMinIssue(path fields.Path, min float64) *Issue {
	return &Issue{
		Path: path,
		Code: CodeNumExclMin,
		Params: Params{
			NumExclMin: &min,
		},
		Message: fmt.Sprintf("%s must be greater than %f", path, min),
	}
}

func NewStrMinLenIssue(path fields.Path, min int64) *Issue {
	return &Issue{
		Path: path,
		Code: CodeStrMinLen,
		Params: Params{
			StrMinLen: &min,
		},
		Message: fmt.Sprintf("%s must be at least %d characters", path, min),
	}
}

func NewStrMaxLenIssue(path fields.Path, max int64) *Issue {
	return &Issue{
		Path: path,
		Code: CodeStrMaxLen,
		Params: Params{
			StrMaxLen: &max,
		},
		Message: fmt.Sprintf("%s must be at most %d characters", path, max),
	}
}

func NewStrLenIssue(path fields.Path, len int64) *Issue {
	return &Issue{
		Path: path,
		Code: CodeStrLen,
		Params: Params{
			StrLen: &len,
		},
		Message: fmt.Sprintf("%s must have %d characters", path, len),
	}
}

func NewStrPatternIssue(path fields.Path, pattern string) *Issue {
	return &Issue{
		Path: path,
		Code: CodeStrPattern,
		Params: Params{
			StrPattern: pattern,
		},
		Message: fmt.Sprintf("%s must match pattern %s", path, pattern),
	}
}

func NewArrMinItemsIssue(path fields.Path, min int64) *Issue {
	return &Issue{
		Path: path,
		Code: CodeArrMinItems,
		Params: Params{
			ArrMinItems: &min,
		},
		Message: fmt.Sprintf("%s must have at least %d items", path, min),
	}
}

func NewArrMaxItemsIssue(path fields.Path, max int64) *Issue {
	return &Issue{
		Path: path,
		Code: CodeArrMaxItems,
		Params: Params{
			ArrMaxItems: &max,
		},
		Message: fmt.Sprintf("%s must have at most %d items", path, max),
	}
}

func NewArrLenIssue(path fields.Path, len int64) *Issue {
	return &Issue{
		Path: path,
		Code: CodeArrLen,
		Params: Params{
			ArrLen: &len,
		},
		Message: fmt.Sprintf("%s must have %d items", path, len),
	}
}

func NewObjMinPropsIssue(path fields.Path, min int64) *Issue {
	return &Issue{
		Path: path,
		Code: CodeObjMinProps,
		Params: Params{
			ObjMinProps: &min,
		},
		Message: fmt.Sprintf("%s must have at least %d properties", path, min),
	}
}

func NewObjMaxPropsIssue(path fields.Path, max int64) *Issue {
	return &Issue{
		Path: path,
		Code: CodeObjMaxProps,
		Params: Params{
			ObjMaxProps: &max,
		},
		Message: fmt.Sprintf("%s must have at most %d properties", path, max),
	}
}

func NewObjLenIssue(path fields.Path, len int64) *Issue {
	return &Issue{
		Path: path,
		Code: CodeObjLen,
		Params: Params{
			ObjLen: &len,
		},
		Message: fmt.Sprintf("%s must have %d properties", path, len),
	}
}
