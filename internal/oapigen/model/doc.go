// Package model contains logic that processes OpenAPI schemas into generic data structures
// suitable for code generation. It deliberately does not emit any code; this is left to generators
// which may use the constructs provided by this package for their own implementation.
//
// # Glossary
//
//   - Declaration: A named, top-level unit of code generation. A Declaration is created for
//     any schema defined at components.schemas, and for certain nested schemas that are
//     "hoisted" (e.g., nested object schemas) and for enums on simple types. Declarations
//     are addressable by an ID (see Location) and contain a Type that describes their shape.
//     Think “what will become a type/alias/const block in the target language”.
//
//   - Type: A structural description used to model shapes. Types can be primitive
//     (int32, string, …), composite (object, array, map), or a TypeRef pointing to a
//     Declaration by ID. Types are anonymous by themselves; they become named/embeddable
//     only when wrapped by a Declaration.
//
//   - Location: Identifies where a Declaration came from. For top-level schemas it’s just the
//     schema name; for hoisted nested schemas it is the root schema name plus a path of
//     segments (e.g., /properties/address, /items, /additionalProperties). The Location’s
//     string form is used as the Declaration ID.
//
// # Example
//
// Consider the following OpenAPI schema:
//
//	components:
//		schemas:
//			User:
//				type: object
//				properties:
//					name:
//						type: string
//					address:
//						type: object
//						properties:
//							street:
//								type: string
//
// In this case, 2 declarations are generated:
//
//	Declaration{
//		ID: "User"
//		Type: Type{
//			Kind: model.TypeObject
//			Fields: []model.Field{
//				// Field models
//			}
//		}
//		Loc: Location{
//			Root: "User"
//			Path: nil
//		}
//	}
//
//	Declaration{
//		ID: "User/properties/address"
//		Type: Type{
//			Kind: model.TypeObject
//			Fields: []model.Field{
//				// Field models
//			}
//		}
//		Loc: Location{
//			Root: "User"
//			Path: []Segment{
//				{Kind: SegmentProperty, Name: "address"},
//			}
//		}
//	}
package model
