package nix

import (
	"github.com/google/jsonschema-go/jsonschema"
)

var (
	_ Type = (*Int)(nil)

	minimumLimitBit byte = 1 << 0
	maximumLimitBit byte = 1 << 1
)

// Int represents an integer option.
type Int struct {
	intLimit
}

// UnmarshalSchema extracts option properties from a JSON Schema. It errors if
// the given schema and this property types mismatch.
func (property *Int) UnmarshalSchema(schema *jsonschema.Schema) error {
	if schema.Type != jsonSchemaInteger {
		return ErrTypeMismatch
	}

	var limits byte

	if schema.Minimum != nil {
		limits |= minimumLimitBit
	}

	if schema.Maximum != nil {
		limits |= maximumLimitBit
	}

	switch limits {
	case minimumLimitBit | maximumLimitBit:
		property.intLimit = intBetween{Minimum: int64(*schema.Minimum), Maximum: int64(*schema.Maximum)}
	case maximumLimitBit:
		property.intLimit = intMaximum{int64(*schema.Maximum)}
	case minimumLimitBit:
		property.intLimit = intMinimum{int64(*schema.Minimum)}
	default:
		property.intLimit = intUnlimited{}
	}

	return nil
}
