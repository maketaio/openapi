package model

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/maketaio/openapi/util/ptr"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

type TypeKind int

const (
	TypeUnknown TypeKind = iota
	TypeInt32
	TypeInt64
	TypeFloat64
	TypeBool
	TypeString
	TypeObject
	TypeArray
	TypeRef
)

type Type struct {
	Kind   TypeKind
	Enum   []EnumConst // For simple kinds where enums can be defined
	Fields []Field     // For struct kind
	Elem   *Type       // For slice, map and struct kind
	Ref    string      // For reference kind, the ID of the declaration being referenced

	// Validation
	Min, Max         *int64   // For int32, int64, string, map, slice
	ExclMin, ExclMax bool     // For int32, int64, float64
	MultipleOf       *int64   // For int32, int64
	Len              *int64   // For string, map, slice
	Pattern          string   // For string
	Format           string   // For string
	MinF, MaxF       *float64 // For float64
}

type EnumConst struct {
	// Name is the enum name extracted as is from x-enum-varnames.
	// An empty string means the enum name is not specified.
	Name    string
	Int32   *int32
	Int64   *int64
	Float64 *float64
	Str     *string
	Bool    *bool
	Doc     []string
}

type Field struct {
	Name       string
	Type       *Type
	Required   bool
	Deprecated bool
	Doc        []string
}

// Declaration represents a top-level type or class to be generated. An OpenAPI schema becomes
// a declaration when it is defined at the top level of components.schemas or when it is an enum or a
// nested object schema.
type Declaration struct {
	ID         string
	Type       *Type
	Doc        []string
	Loc        Location
	Deprecated bool
}

// Registry collects and stores declarations
type Registry struct {
	decls map[string]*Declaration
	ids   []string
}

func NewRegistry() *Registry {
	return &Registry{
		decls: map[string]*Declaration{},
	}
}

func (r *Registry) Collect(dm *libopenapi.DocumentModel[v3.Document]) error {
	for pair := dm.Model.Components.Schemas.First(); pair != nil; pair = pair.Next() {
		l := Location{Root: pair.Key()}
		_, err := r.visit(l, pair.Value())
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Registry) Get(id string) (*Declaration, bool) {
	m, ok := r.decls[id]
	return m, ok
}

// Range iterates declarations in insertion order (schema document order and nested traversal order).
// To stop early, return false.
func (r *Registry) Range(fn func(id string, decl *Declaration) bool) {
	for _, id := range r.ids {
		if !fn(id, r.decls[id]) {
			return
		}
	}
}

// addDecl adds a declaration and returns its ID built from path
func (r *Registry) addDecl(l Location, typ *Type, schema *base.Schema) string {
	m := &Declaration{
		ID:   l.String(),
		Type: typ,
		Loc:  l,
	}

	if schema != nil {
		m.Doc = toDocLines(schema.Description)
		m.Deprecated = ptr.Deref(schema.Deprecated, false)
	}

	r.decls[m.ID] = m
	r.ids = append(r.ids, m.ID)
	return m.ID
}

// visit traverses a schema proxy and converts it into an internal Type representation.
// For top-level schemas, certain nested schemas or enum schemas, visit will register a new
// declaration and return a Type referencing it. For all other cases, it returns the constructed
// Type directly.
func (r *Registry) visit(l Location, sp *base.SchemaProxy) (*Type, error) {
	if sp.IsReference() {
		loc := sp.GetReference()
		parts := strings.Split(loc, "/")
		typ := &Type{
			Kind: TypeRef,
			Ref:  parts[len(parts)-1],
		}

		if l.IsTopLevel() {
			return &Type{
				Kind: TypeRef,
				Ref:  r.addDecl(l, typ, nil),
			}, nil
		}

		return typ, nil
	}

	schema := sp.Schema()

	if schema.Type == nil || len(schema.Type) == 0 {
		return nil, fmt.Errorf("schema %s has no type", l)
	}

	if len(schema.Type) > 1 {
		return nil, fmt.Errorf("schema %s has multiple types, which is not supported at the moment", l)
	}

	switch schema.Type[0] {
	case "string":
		return r.visitStr(l, schema)
	case "integer":
		return r.visitInt(l, schema)
	case "number":
		return r.visitNum(l, schema)
	case "boolean":
		return r.visitBool(l, schema)
	case "array":
		return r.visitArr(l, schema)
	case "object":
		return r.visitObj(l, schema)
	default:
		return nil, fmt.Errorf("unhandled type %s for %s", schema.Type[0], l)
	}
}

func (r *Registry) visitStr(l Location, schema *base.Schema) (*Type, error) {
	typ := &Type{
		Kind:    TypeString,
		Pattern: schema.Pattern,
		Format:  schema.Format,
	}

	if schema.MaxLength != nil && schema.MinLength != nil && *schema.MaxLength == *schema.MinLength {
		typ.Len = schema.MaxLength
	} else {
		typ.Min = schema.MinLength
		typ.Max = schema.MaxLength
	}

	if len(schema.Enum) > 0 {
		var err error
		typ.Enum, err = makeConsts(schema, func(c *EnumConst, v string) { c.Str = &v })
		if err != nil {
			return nil, err
		}
	}

	if l.IsTopLevel() || len(typ.Enum) > 0 {
		return &Type{
			Kind: TypeRef,
			Ref:  r.addDecl(l, typ, schema),
		}, nil
	}

	return typ, nil
}

func (r *Registry) visitInt(l Location, schema *base.Schema) (*Type, error) {
	var typ *Type

	if schema.Format == "int32" {
		typ = &Type{Kind: TypeInt32}

		if len(schema.Enum) > 0 {
			var err error
			typ.Enum, err = makeConsts(schema, func(c *EnumConst, v int32) { c.Int32 = &v })
			if err != nil {
				return nil, err
			}
		}
	} else {
		typ = &Type{Kind: TypeInt64}

		if len(schema.Enum) > 0 {
			var err error
			typ.Enum, err = makeConsts(schema, func(c *EnumConst, v int64) { c.Int64 = &v })
			if err != nil {
				return nil, err
			}
		}
	}

	b := getNumericBounds(schema)

	if b.max != nil {
		if math.Trunc(*b.max) != *b.max {
			return nil, fmt.Errorf("schema %s has a maximum or exclusiveMaximum that is not an integer", l)
		}

		typ.Max = ptr.To(int64(*b.max))
		typ.ExclMax = b.exclMax
	}

	if b.min != nil {
		if math.Trunc(*b.min) != *b.min {
			return nil, fmt.Errorf("schema %s has a minimum or exclusiveMinimum that is not an integer", l)
		}

		typ.Min = ptr.To(int64(*b.min))
		typ.ExclMin = b.exclMin
	}

	if schema.MultipleOf != nil {
		if math.Trunc(*schema.MultipleOf) != *schema.MultipleOf {
			return nil, fmt.Errorf("schema %s has a multipleOf that is not an integer", l)
		}

		typ.MultipleOf = ptr.To(int64(*schema.MultipleOf))
	}

	if l.IsTopLevel() || len(typ.Enum) > 0 {
		return &Type{
			Kind: TypeRef,
			Ref:  r.addDecl(l, typ, schema),
		}, nil
	}

	return typ, nil
}

func (r *Registry) visitNum(l Location, schema *base.Schema) (*Type, error) {
	typ := &Type{Kind: TypeFloat64}

	if len(schema.Enum) > 0 {
		var err error
		typ.Enum, err = makeConsts(schema, func(c *EnumConst, v float64) { c.Float64 = &v })
		if err != nil {
			return nil, err
		}
	}

	b := getNumericBounds(schema)
	typ.MaxF, typ.ExclMax = b.max, b.exclMax
	typ.MinF, typ.ExclMin = b.min, b.exclMin

	if l.IsTopLevel() || len(typ.Enum) > 0 {
		return &Type{
			Kind: TypeRef,
			Ref:  r.addDecl(l, typ, schema),
		}, nil
	}

	return typ, nil
}

func (r *Registry) visitBool(l Location, schema *base.Schema) (*Type, error) {
	typ := &Type{Kind: TypeBool}

	if len(schema.Enum) > 0 {
		var err error
		typ.Enum, err = makeConsts(schema, func(c *EnumConst, v bool) { c.Bool = &v })
		if err != nil {
			return nil, err
		}
	}

	if l.IsTopLevel() || len(typ.Enum) > 0 {
		return &Type{
			Kind: TypeRef,
			Ref:  r.addDecl(l, typ, schema),
		}, nil
	}

	return typ, nil
}

func (r *Registry) visitArr(l Location, schema *base.Schema) (*Type, error) {
	typ := &Type{
		Kind: TypeArray,
	}

	if schema.MaxItems != nil && schema.MinItems != nil && *schema.MaxItems == *schema.MinItems {
		typ.Len = schema.MaxItems
	} else {
		typ.Min = schema.MinItems
		typ.Max = schema.MaxItems
	}

	if schema.Items == nil {
		typ.Elem = &Type{
			Kind: TypeUnknown,
		}
	} else if schema.Items.IsB() {
		typ.Elem = &Type{
			Kind: TypeUnknown,
		}

		if !schema.Items.B {
			typ.Len = ptr.To(int64(0))
			typ.Min = nil
			typ.Max = nil
		}
	} else {
		var err error
		typ.Elem, err = r.visit(l.WithItems(), schema.Items.A)
		if err != nil {
			return nil, err
		}
	}

	if l.IsTopLevel() {
		return &Type{
			Kind: TypeRef,
			Ref:  r.addDecl(l, typ, schema),
		}, nil
	}

	return typ, nil
}

func (r *Registry) visitObj(l Location, schema *base.Schema) (*Type, error) {
	typ := &Type{
		Kind: TypeObject,
	}

	if schema.MaxProperties != nil && schema.MinProperties != nil && *schema.MaxProperties == *schema.MinProperties {
		typ.Len = schema.MaxProperties
	} else {
		typ.Min = schema.MinProperties
		typ.Max = schema.MaxProperties
	}

	var err error
	typ.Elem, err = r.visitAdditionalProps(l, schema)
	if err != nil {
		return nil, err
	}

	if orderedmap.Len(schema.Properties) == 0 {
		if l.IsTopLevel() {
			return &Type{
				Kind: TypeRef,
				Ref:  r.addDecl(l, typ, schema),
			}, nil
		}

		return typ, nil
	}

	typ.Fields = make([]Field, 0, orderedmap.Len(schema.Properties))

	for prop := schema.Properties.First(); prop != nil; prop = prop.Next() {
		name := prop.Key()

		ft, err := r.visit(l.WithProperty(name), prop.Value())
		if err != nil {
			return nil, err
		}

		propSchema := prop.Value().Schema()

		typ.Fields = append(typ.Fields, Field{
			Name:       name,
			Type:       ft,
			Required:   slices.Contains(schema.Required, name),
			Deprecated: ptr.Deref(propSchema.Deprecated, false),
			Doc:        toDocLines(propSchema.Description),
		})
	}

	if l.IsTopLevel() {
		return &Type{
			Kind: TypeRef,
			Ref:  r.addDecl(l, typ, schema),
		}, nil
	}

	return typ, nil
}

func (r *Registry) visitAdditionalProps(l Location, schema *base.Schema) (*Type, error) {
	if schema.AdditionalProperties == nil {
		return &Type{
			Kind: TypeUnknown,
		}, nil
	}

	if schema.AdditionalProperties.IsB() {
		if !schema.AdditionalProperties.B {
			return nil, nil
		}

		return &Type{
			Kind: TypeUnknown,
		}, nil
	}

	return r.visit(l.WithAdditionalProperties(), schema.AdditionalProperties.A)
}

func toDocLines(doc string) []string {
	if doc == "" {
		return nil
	}

	normalized := strings.ReplaceAll(strings.ReplaceAll(doc, "\r\n", "\n"), "\r", "\n")
	return strings.Split(normalized, "\n")
}

func makeConsts[T any](schema *base.Schema, build func(*EnumConst, T)) ([]EnumConst, error) {
	if len(schema.Enum) == 0 {
		return nil, nil
	}

	var varnames []string
	var descriptions []string

	if orderedmap.Len(schema.Extensions) > 0 {
		xvarname, found := schema.Extensions.Get("x-enum-varnames")
		if found {
			if err := xvarname.Decode(&varnames); err != nil {
				return nil, fmt.Errorf("failed to unmarshal x-enum-varnames: %w", err)
			}
		}

		xdescription, found := schema.Extensions.Get("x-enum-descriptions")
		if found {
			if err := xdescription.Decode(&descriptions); err != nil {
				return nil, fmt.Errorf("failed to unmarshal x-enum-descriptions: %w", err)
			}
		}
	}

	var result []EnumConst

	for i, n := range schema.Enum {
		var v T
		if err := n.Decode(&v); err != nil {
			return nil, fmt.Errorf("failed to unmarshal enum value %d: %w", i, err)
		}

		name := ""
		if i < len(varnames) {
			name = varnames[i]
		}

		desc := ""
		if i < len(descriptions) {
			desc = descriptions[i]
		}

		c := EnumConst{
			Name: name,
			Doc:  toDocLines(desc),
		}

		build(&c, v)
		result = append(result, c)
	}

	return result, nil
}

type numericBounds struct {
	min, max         *float64
	exclMin, exclMax bool
}

func getNumericBounds(schema *base.Schema) numericBounds {
	b := numericBounds{}
	if schema.Maximum != nil {
		b.max = schema.Maximum
	}

	if schema.Minimum != nil {
		b.min = schema.Minimum
	}

	if schema.ExclusiveMaximum != nil && (schema.ExclusiveMaximum.IsB() || schema.ExclusiveMaximum.A) {
		if schema.ExclusiveMaximum.IsB() {
			if b.max == nil {
				b.max = &schema.ExclusiveMaximum.B
				b.exclMax = true
			} else {
				exclMax := schema.ExclusiveMaximum.B
				if exclMax <= *b.max {
					b.max = &exclMax
					b.exclMax = true
				}
			}
		} else if b.max != nil {
			b.exclMax = true
		}
	}

	if schema.ExclusiveMinimum != nil && (schema.ExclusiveMinimum.IsB() || schema.ExclusiveMinimum.A) {
		if schema.ExclusiveMinimum.IsB() {
			if b.min == nil {
				b.min = &schema.ExclusiveMinimum.B
				b.exclMin = true
			} else {
				exclMin := schema.ExclusiveMinimum.B
				if exclMin >= *b.min {
					b.min = &exclMin
					b.exclMin = true
				}
			}
		} else if b.min != nil {
			b.exclMin = true
		}
	}

	return b
}
