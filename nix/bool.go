package nix

import (
	"github.com/google/jsonschema-go/jsonschema"
)

var _ Type = (*Bool)(nil)

// Bool represents a boolean option.
type Bool struct{}

// UnmarshalSchema extracts option properties from a JSON Schema. It errors if
// the given schema and this property types mismatch.
func (*Bool) UnmarshalSchema(schema *jsonschema.Schema) error {
	if schema.Type != jsonSchemaBoolean {
		return ErrTypeMismatch
	}

	return nil
}

// String marshals this property as a nix type definition.
func (*Bool) String() string {
	return "bool"
}
