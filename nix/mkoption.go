package nix

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
)

var _ Option = (*MkOption)(nil)

// MkOption represents a nix option, or more precisely a nixpkgs'
// lib.options.mkOption attributes. It provides a way to configure an option
// and its underlying type, and ultimately generate its nix source code
// representation.
type MkOption struct {
	Type

	// applyFn      string
	defaultValue string
	// defaultText  string
	description string
	// example      string
	// internal     bool
	// nullable bool
	// readOnly     bool
	// visible      bool
}

// NewMkOptionFromSchema unmarshals a JSON schema into its equivalent Nix
// option. Types call it recursively where applicable (objects, arrays,
// oneOf/anyOf, etc).
func NewMkOptionFromSchema(schema *jsonschema.Schema) (*MkOption, error) {
	var option MkOption

	err := option.UnmarshalSchema(schema)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling schema %q into nix %q: %w", schema.Type, reflect.TypeOf(option.AsType()), err)
	}

	return &option, nil
}

// UnmarshalSchema extracts option properties from a JSON Schema. It errors if
// the given schema and this property types mismatch.
func (opt *MkOption) UnmarshalSchema(schema *jsonschema.Schema) error {
	var err error

	opt.Type, err = NewTypeFromSchema(schema)
	if err != nil {
		return fmt.Errorf("unmarshalling schema into type: %w", err)
	}

	opt.defaultValue = string(schema.Default)
	opt.description = schema.Description

	return nil
}

func (opt *MkOption) String() string {
	var buf bytes.Buffer

	buf.WriteString("mkOption {")

	optType := opt.Type.String()

	if opt.defaultValue != "" {
		fmt.Fprintf(&buf, "default = %s;", opt.defaultValue)
	} else {
		if strings.HasPrefix(optType, "listOf") {
			buf.WriteString("default = [];")
		} else {
			buf.WriteString("default = null;")

			optType = "nullOr " + SafeExpression(optType)
		}
	}

	fmt.Fprintf(&buf, "type = %s;", optType)

	if opt.description != "" {
		fmt.Fprintf(&buf, "description = %q;", opt.description)
	}

	buf.WriteString("}")

	return buf.String()
}

// AsType returns the bare option type.
func (opt *MkOption) AsType() Type {
	return opt.Type
}
