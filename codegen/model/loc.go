package model

type SegmentKind int

const (
	// SegmentProperty is used for object property segments.
	SegmentProperty SegmentKind = iota
	// SegmentAdditionalProperties is used for additional properties segments.
	SegmentAdditionalProperties
	// SegmentItems is used for array items segments.
	SegmentItems
)

// Segment represents a segment of a path to a model.
type Segment struct {
	// Kind shows what type of segment this is.
	Kind SegmentKind
	// Name distinguishes segments of the same kind, if applicable.
	Name string
}

// Location represents a location of a model. If a model is a top level schema, then its loc is simply
// the schema name. If a model is "hoisted" because it was a nested object schema, then its loc
// contains the parent schema name plus the path to the nested schema.
type Location struct {
	// Root is the name of the top level schema.
	Root string
	// Path contains the path to the nested model, if applicable.
	Path []Segment
}

func (l Location) String() string {
	loc := l.Root
	for _, seg := range l.Path {
		switch seg.Kind {
		case SegmentProperty:
			loc += "/properties/" + seg.Name
		case SegmentAdditionalProperties:
			loc += "/additionalProperties"
		case SegmentItems:
			loc += "/items"
		}
	}
	return loc
}

func (l Location) IsTopLevel() bool {
	return len(l.Path) == 0
}

func (l Location) WithProperty(name string) Location {
	return Location{
		Root: l.Root,
		Path: append(l.Path, Segment{
			Kind: SegmentProperty,
			Name: name,
		}),
	}
}

func (l Location) WithAdditionalProperties() Location {
	return Location{
		Root: l.Root,
		Path: append(l.Path, Segment{
			Kind: SegmentAdditionalProperties,
		}),
	}
}

func (l Location) WithItems() Location {
	return Location{
		Root: l.Root,
		Path: append(l.Path, Segment{
			Kind: SegmentItems,
		}),
	}
}
