package nix

import (
	"github.com/google/jsonschema-go/jsonschema"
)

var _ Type = (*List)(nil)

// List represents a list option.
type List struct {
	subType Option
}

// UnmarshalSchema extracts option properties from a JSON Schema. It errors if
// the given schema and this property types mismatch.
func (property *List) UnmarshalSchema(schema *jsonschema.Schema) error {
	if schema.Type != jsonSchemaArray {
		return ErrTypeMismatch
	}

	subType, err := NewMkOptionFromSchema(schema.Items)
	if err != nil {
		return err
	}

	property.subType = subType

	return nil
}

// String marshals this property as a nix type definition.
func (property *List) String() string {
	return "listOf " + property.subType.String()
}
