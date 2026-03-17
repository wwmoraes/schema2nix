package nix

import (
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
)

var _ Type = (*Raw)(nil)

// Raw represents an unknown/undefined type that accepts any arbitrary value.
type Raw struct {
	schemaType string
}

// UnmarshalSchema extracts option properties from a JSON Schema. It errors if
// the given schema and this property types mismatch.
func (property *Raw) UnmarshalSchema(schema *jsonschema.Schema) error {
	property.schemaType = schema.Type

	return nil
}

// String marshals this property as a nix type definition.
func (property *Raw) String() string {
	return fmt.Sprintf("raw /* TODO add JSON Schema %q support */", property.schemaType)
}
