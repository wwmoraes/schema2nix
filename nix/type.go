package nix

import (
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
)

// Type is a value that can retrieve its constraints from a JSON Schema
// definition and represent itself as a nix type.
type Type interface {
	SchemaUnmarshaler
	fmt.Stringer
}

// NewTypeFromSchema returns a nix [Type] that represents the schema type.
func NewTypeFromSchema(schema *jsonschema.Schema) (Type, error) {
	var result Type

	switch schema.Type {
	case jsonSchemaString:
		result = &Str{}
	case jsonSchemaBoolean:
		result = &Bool{}
	case jsonSchemaObject:
		result = &Attrs{}
	case jsonSchemaInteger:
		result = &Int{}
	case jsonSchemaNumber:
		result = &Number{}
	case jsonSchemaArray:
		result = &List{}
	default:
		result = &Anything{}
	}

	err := result.UnmarshalSchema(schema)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling schema into nix type: %w", err)
	}

	return result, nil
}
