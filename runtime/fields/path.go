package fields

import "strings"

type Path []string

func (p Path) String() string {
	return strings.Join(p, ".")
}

func (p Path) Field(s string) Path {
	return append(p, s)
}
