package generator

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/maketaio/api/internal/util/set"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type kind int

const (
	kindUnknown kind = iota
	kindInt32
	kindInt64
	kindFloat64
	kindBool
	kindString
	kindTime
	kindBytes
	kindStruct
	kindSlice
	kindMap
	kindJSONRaw
	kindRef // named alias/reference to a top-level model
)

type goType struct {
	kind   kind
	enum   []enumConst // For simple kinds where enums can be defined
	fields []goField   // For struct kind
	elem   *goType     // For slice, map and struct kind
	ref    string      // For reference kind

	// Validation
	min, max, exclMin, exclMax     *int     // For int32, int64, string, map, slice
	multipleOf                     int      // For int32, int64
	len                            *int     // For string, map, slice
	minF, maxF, exclMinF, exclMaxF *float64 // For float64
}

type enumConst struct {
	name string
	lit  string
	doc  []string
}

type goField struct {
	name       string
	jsonName   string
	typ        *goType
	required   bool
	deprecated bool
	doc        []string
}

// decl represents a top level type declaration to be generated
// an OpenAPI schema is considered top level declaration when it is:
//   - defined at the top level of components.schemas or
//   - an enum or
//   - a nested object schema
type decl struct {
	name       string
	typ        *goType
	doc        []string
	deprecated bool
	path       []string
}

// collector is used to collect top level type declarations
type collector struct {
	decls   []*decl
	names   set.Set[string]
	counter int
}

func newCollector() *collector {
	return &collector{
		names: set.NewSet[string](),
	}
}

func (c *collector) addDecl(name string, typ *goType, schema *base.Schema, path []string) string {
	name = toTitle(name)

	for {
		if !c.names.Has(name) {
			break
		}
		c.counter++
		name = name + strconv.Itoa(c.counter)
	}

	c.names.Add(name)
	c.decls = append(c.decls, &decl{
		name:       name,
		typ:        typ,
		doc:        toDocLines(schema.Description),
		deprecated: deref(schema.Deprecated, false),
		path:       path,
	})

	return name
}

// walk recursively walks an OpenAPI schema for top level type declarations,
// and returns a goType that is a direct type or references the declaration.
func (c *collector) walk(name string, path []string, schema *base.Schema) (*goType, error) {
	topLevel := len(path) == 1

	if schema.Type == nil || len(schema.Type) == 0 {
		return nil, fmt.Errorf("schema %s has no type", strings.Join(path, "/"))
	}

	switch schema.Type[0] {
	case "string":
		if len(schema.Enum) > 0 {
			enum, err := makeConsts(name, schema, strconv.Quote, func(val string) string { return val })
			if err != nil {
				return nil, err
			}

			name := c.addDecl(name, &goType{kind: kindString, enum: enum}, schema, path)
			return &goType{
				kind: kindRef,
				ref:  name,
			}, nil
		}

		switch schema.Format {
		case "date", "date-time":
			if topLevel {
				name := c.addDecl(name, &goType{kind: kindTime}, schema, path)
				return &goType{
					kind: kindRef,
					ref:  name,
				}, nil
			}

			return &goType{
				kind: kindTime,
			}, nil

		case "binary", "byte":
			if topLevel {
				name := c.addDecl(name, &goType{kind: kindBytes}, schema, path)
				return &goType{
					kind: kindRef,
					ref:  name,
				}, nil
			}

			return &goType{
				kind: kindBytes,
			}, nil

		default:
			if topLevel {
				name := c.addDecl(name, &goType{kind: kindString}, schema, path)
				return &goType{
					kind: kindRef,
					ref:  name,
				}, nil
			}

			return &goType{
				kind: kindString,
			}, nil
		}

	case "integer":
		switch schema.Format {
		case "int64":
			if len(schema.Enum) > 0 {
				enum, err := makeConsts(name, schema, func(v int64) string {
					return strconv.FormatInt(v, 10)
				}, func(v int64) string {
					return strconv.FormatInt(v, 10)
				})
				if err != nil {
					return nil, err
				}

				name := c.addDecl(name, &goType{kind: kindInt64, enum: enum}, schema, path)
				return &goType{
					kind: kindRef,
					ref:  name,
				}, nil
			}

			if topLevel {
				name := c.addDecl(name, &goType{kind: kindInt64}, schema, path)
				return &goType{
					kind: kindRef,
					ref:  name,
				}, nil
			}

			return &goType{
				kind: kindInt64,
			}, nil
		}

		if len(schema.Enum) > 0 {
			enum, err := makeConsts(name, schema, func(v int32) string {
				return strconv.FormatInt(int64(v), 10)
			}, func(v int32) string {
				return strconv.FormatInt(int64(v), 10)
			})
			if err != nil {
				return nil, err
			}

			name := c.addDecl(name, &goType{kind: kindInt32, enum: enum}, schema, path)
			return &goType{
				kind: kindRef,
				ref:  name,
			}, nil
		}

		if topLevel {
			name := c.addDecl(name, &goType{kind: kindInt32}, schema, path)
			return &goType{
				kind: kindRef,
				ref:  name,
			}, nil
		}

		return &goType{
			kind: kindInt32,
		}, nil

	case "number":
		if len(schema.Enum) > 0 {
			enum, err := makeConsts(name, schema, func(v float64) string {
				return strconv.FormatFloat(v, 'f', -1, 64)
			}, func(v float64) string {
				return strings.Replace(strconv.FormatFloat(v, 'f', -1, 64), ".", "_", -1)
			})
			if err != nil {
				return nil, err
			}

			name := c.addDecl(name, &goType{kind: kindFloat64, enum: enum}, schema, path)
			return &goType{
				kind: kindRef,
				ref:  name,
			}, nil
		}

		if topLevel {
			name := c.addDecl(name, &goType{kind: kindFloat64}, schema, path)
			return &goType{
				kind: kindRef,
				ref:  name,
			}, nil
		}

		return &goType{
			kind: kindFloat64,
		}, nil

	case "boolean":
		if len(schema.Enum) > 0 {
			enum, err := makeConsts(name, schema, strconv.FormatBool, strconv.FormatBool)
			if err != nil {
				return nil, err
			}

			name := c.addDecl(name, &goType{kind: kindBool, enum: enum}, schema, path)
			return &goType{
				kind: kindRef,
				ref:  name,
			}, nil
		}

		if topLevel {
			name := c.addDecl(name, &goType{kind: kindBool}, schema, path)
			return &goType{
				kind: kindRef,
				ref:  name,
			}, nil
		}

		return &goType{
			kind: kindBool,
		}, nil

	case "array":
		if schema.Items == nil || schema.Items.IsB() {
			typ := &goType{
				kind: kindSlice,
				elem: &goType{
					kind: kindJSONRaw,
				},
			}

			if schema.Items != nil && schema.Items.IsB() && !schema.Items.B {
				typ.elem = &goType{
					kind: kindStruct,
				}

				zero := 0
				typ.len = &zero
			}

			if topLevel {
				name := c.addDecl(name, typ, schema, path)
				return &goType{
					kind: kindRef,
					ref:  name,
				}, nil
			}

			return typ, nil
		}

		typ, err := c.walk(toTitle(name)+"Item", append(path, "*"), schema.Items.A.Schema())
		if err != nil {
			return nil, err
		}

		if topLevel {
			name := c.addDecl(name, &goType{kind: kindSlice, elem: typ}, schema, path)
			return &goType{
				kind: kindRef,
				ref:  name,
			}, nil
		}

		return &goType{
			kind: kindSlice,
			elem: typ,
		}, nil

	case "object":
		additionalPropMap, err := c.walkAdditionalProps(name, append(path, "*"), schema)
		if err != nil {
			return nil, fmt.Errorf("failed to get additional properties type for %s: %w", strings.Join(path, "/"), err)
		}

		if orderedmap.Len(schema.Properties) == 0 {
			if additionalPropMap == nil {
				if topLevel {
					name := c.addDecl(name, &goType{kind: kindStruct}, schema, path)
					return &goType{
						kind: kindRef,
						ref:  name,
					}, nil
				}

				return &goType{
					kind: kindStruct,
				}, nil
			}

			if topLevel {
				name := c.addDecl(name, additionalPropMap, schema, path)

				return &goType{
					kind: kindRef,
					ref:  name,
				}, nil
			}

			return additionalPropMap, nil
		}

		if orderedmap.Len(schema.Properties) == 0 {
			if topLevel {
				name := c.addDecl(name, &goType{kind: kindStruct}, schema, path)
				return &goType{
					kind: kindRef,
					ref:  name,
				}, nil
			}

			return &goType{
				kind: kindStruct,
			}, nil
		}

		fields := make([]goField, 0, orderedmap.Len(schema.Properties))

		for prop := schema.Properties.First(); prop != nil; prop = prop.Next() {
			propName := prop.Key()
			goName := toTitle(propName)
			propSchema := prop.Value().Schema()

			typ, err := c.walk(toTitle(name)+goName, append(path, propName), propSchema)
			if err != nil {
				return nil, err
			}

			fields = append(fields, goField{
				name:       goName,
				jsonName:   propName,
				typ:        typ,
				required:   slices.Contains(schema.Required, propName),
				deprecated: deref(propSchema.Deprecated, false),
				doc:        toDocLines(propSchema.Description),
			})
		}

		typ := &goType{
			kind:   kindStruct,
			fields: fields,
		}

		if additionalPropMap != nil {
			typ.elem = additionalPropMap.elem
		}

		name := c.addDecl(name, typ, schema, path)

		return &goType{
			kind: kindRef,
			ref:  name,
		}, nil

	default:
		return nil, fmt.Errorf("unhandled type %s for %s", schema.Type[0], strings.Join(path, "/"))
	}
}

func toDocLines(doc string) []string {
	if doc == "" {
		return nil
	}

	normalized := strings.ReplaceAll(strings.ReplaceAll(doc, "\r\n", "\n"), "\r", "\n")
	return strings.Split(normalized, "\n")
}

func makeConsts[T any](name string, schema *base.Schema, litFunc func(T) string, suffixFunc func(T) string) ([]enumConst, error) {
	if len(schema.Enum) == 0 {
		return nil, nil
	}

	var varnames []string
	var descriptions []string

	if orderedmap.Len(schema.Extensions) > 0 {
		varnameExt, found := schema.Extensions.Get("x-enum-varnames")
		if found {
			if err := varnameExt.Decode(&varnames); err != nil {
				return nil, fmt.Errorf("failed to unmarshal x-enum-varnames: %w", err)
			}
		}

		descriptionExt, found := schema.Extensions.Get("x-enum-descriptions")
		if found {
			if err := descriptionExt.Decode(&descriptions); err != nil {
				return nil, fmt.Errorf("failed to unmarshal x-enum-descriptions: %w", err)
			}
		}
	}

	var result []enumConst

	for i, n := range schema.Enum {
		var v T
		if err := n.Decode(&v); err != nil {
			return nil, fmt.Errorf("failed to unmarshal enum value %d: %w", i, err)
		}

		suffix := ""
		if len(varnames) > 0 {
			suffix = varnames[i]
		}

		if len(suffix) == 0 {
			suffix = suffixFunc(v)
		}

		desc := ""
		if len(descriptions) > 0 {
			desc = descriptions[i]
		}

		c := enumConst{
			name: name + toTitle(suffix),
			lit:  litFunc(v),
			doc:  toDocLines(desc),
		}

		result = append(result, c)
	}

	return result, nil
}

func (c *collector) walkAdditionalProps(name string, path []string, schema *base.Schema) (*goType, error) {
	if schema.AdditionalProperties == nil {
		return &goType{
			kind: kindMap,
			elem: &goType{
				kind: kindJSONRaw,
			},
		}, nil
	}

	if schema.AdditionalProperties.IsB() {
		if !schema.AdditionalProperties.B {
			return nil, nil
		}

		return &goType{
			kind: kindMap,
			elem: &goType{
				kind: kindJSONRaw,
			},
		}, nil
	}

	typ, err := c.walk(toTitle(name)+"Entry", append(path, "*"), schema.AdditionalProperties.A.Schema())
	if err != nil {
		return nil, err
	}

	return &goType{
		kind: kindMap,
		elem: typ,
	}, nil
}

func toTitle(s string) string {
	if s == "" {
		return ""
	}

	return cases.Title(language.English, cases.NoLower).String(s)
}

func deref[T any](v *T, def T) T {
	if v == nil {
		return def
	}

	return *v
}
