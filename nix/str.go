package nix

import (
	"fmt"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
)

// Str represents a string option.
type Str struct {
	Enum []string
}

// UnmarshalSchema extracts option properties from a JSON Schema. It errors if
// the given schema and this property types mismatch.
func (property *Str) UnmarshalSchema(schema *jsonschema.Schema) error {
	if schema.Type != jsonSchemaString {
		return ErrTypeMismatch
	}

	if len(schema.Enum) > 0 {
		property.Enum = make([]string, 0, len(schema.Enum))
		for _, entry := range schema.Enum {
			value, ok := entry.(string)
			if !ok {
				return ErrEnumValueTypeMismatch
			}

			property.Enum = append(property.Enum, value)
		}
	}

	return nil
}

// String marshals this property as a nix type definition.
func (property *Str) String() string {
	if len(property.Enum) > 0 {
		//nolint:gocritic // %q escapes the separators
		return fmt.Sprintf(`enum ["%s"]`, strings.Join(property.Enum, `" "`))
	}

	return "str"
}
