package nix

import (
	"github.com/google/jsonschema-go/jsonschema"
)

var _ Type = (*Number)(nil)

// Number represents a floating point number option.
type Number struct{}

// UnmarshalSchema extracts option properties from a JSON Schema. It errors if
// the given schema and this property types mismatch.
func (*Number) UnmarshalSchema(schema *jsonschema.Schema) error {
	if schema.Type != jsonSchemaNumber {
		return ErrTypeMismatch
	}

	return nil
}

// String marshals this property as a nix type definition.
func (*Number) String() string {
	return "number"
}
