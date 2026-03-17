package nix

import (
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
)

var _ Type = (*Anything)(nil)

// Anything represents an unknown/undefined type that accepts any arbitrary value.
type Anything struct {
	schemaType string
}

// UnmarshalSchema extracts option properties from a JSON Schema. It errors if
// the given schema and this property types mismatch.
func (property *Anything) UnmarshalSchema(schema *jsonschema.Schema) error {
	property.schemaType = schema.Type

	return nil
}

// String marshals this property as a nix type definition.
func (property *Anything) String() string {
	return fmt.Sprintf("anything /* TODO add JSON Schema %q support */", property.schemaType)
}
