package nix

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
)

var _ Type = (*Attrs)(nil)

// Attrs represents an attribute set option.
type Attrs struct {
	Properties   map[string]Option
	FreeformType string
	IsFreeform   bool
}

// UnmarshalSchema extracts a JSON schema object type definition into its Nix
// attribute set equivalent.
//
//nolint:gocognit // TODO refactor
func (property *Attrs) UnmarshalSchema(schema *jsonschema.Schema) error {
	if schema.Type != jsonSchemaObject {
		return ErrTypeMismatch
	}

	if isFreeformType(schema) {
		property.FreeformType = "lazyAttrsOf anything"
	}

	property.Properties = make(map[string]Option, len(schema.Properties))

	var (
		prop Option
		err  error
	)

	for name, schema := range schema.Properties {
		option, err := NewMkOptionFromSchema(schema)
		if err != nil {
			return fmt.Errorf("unmarshalling schema %q into nix %q: %w", schema.Type, reflect.TypeOf(option.AsType()), err)
		}

		property.Properties[name] = option
	}

	if schema.AdditionalProperties != nil {
		for name, schema := range schema.AdditionalProperties.Properties {
			prop, err = NewMkOptionFromSchema(schema)
			if err != nil {
				fmt.Fprintf(os.Stderr, "/* TODO %q property (type %q) */", name, schema.Type)

				continue
			}

			property.Properties[name] = prop
		}
	}

	return nil
}

// String marshals this property as a nix type definition.
func (property *Attrs) String() string {
	if len(property.Properties) == 0 {
		return "attrs"
	}

	var buf strings.Builder

	buf.WriteString("submodule {")

	if property.FreeformType != "" {
		fmt.Fprintf(&buf, "freeformType = %s;", property.FreeformType)
	}

	buf.WriteString("options = {")

	for name, prop := range property.Properties {
		fmt.Fprintf(&buf, "%s = %s;", SafeIdentifier(name), prop.String())
	}

	buf.WriteString("};}")

	return buf.String()
}

func isFreeformType(schema *jsonschema.Schema) bool {
	if schema == nil {
		return true
	}

	if schema.AdditionalProperties == nil {
		return true
	}

	if schema.AdditionalProperties.Not != nil && reflect.ValueOf(*schema.AdditionalProperties.Not).IsZero() {
		return false
	}

	return true
}
