package nix_test

import (
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/stretchr/testify/require"

	"github.com/wwmoraes/schema2nix/nix"
)

func TestInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		schema      *jsonschema.Schema
		assertError require.ErrorAssertionFunc
		want        string
	}{
		{
			name: "without limits",
			schema: &jsonschema.Schema{
				Type: "integer",
			},
			want:        `mkOption {default = null;type = nullOr int;}`,
			assertError: require.NoError,
		},
		{
			name: "with unsigned (>=0) limit",
			schema: &jsonschema.Schema{
				Type:    "integer",
				Minimum: jsonschema.Ptr(float64(0)),
			},
			want:        `mkOption {default = null;type = nullOr ints.unsigned;}`,
			assertError: require.NoError,
		},
		{
			name: "with positive (>0) limit",
			schema: &jsonschema.Schema{
				Type:    "integer",
				Minimum: jsonschema.Ptr(float64(1)),
			},
			want:        `mkOption {default = null;type = nullOr ints.positive;}`,
			assertError: require.NoError,
		},
		{
			name: "with arbitrary minimum",
			schema: &jsonschema.Schema{
				Type:    "integer",
				Minimum: jsonschema.Ptr(float64(42)),
			},
			want:        `mkOption {default = null;type = nullOr (addCheck int (x: x >= 42));}`,
			assertError: require.NoError,
		},
		{
			name: "with arbitrary maximum",
			schema: &jsonschema.Schema{
				Type:    "integer",
				Maximum: jsonschema.Ptr(float64(42)),
			},
			want:        `mkOption {default = null;type = nullOr (addCheck int (x: x <= 42));}`,
			assertError: require.NoError,
		},
		{
			name: "between values",
			schema: &jsonschema.Schema{
				Type:    "integer",
				Minimum: jsonschema.Ptr(float64(1)),
				Maximum: jsonschema.Ptr(float64(10)),
			},
			want:        `mkOption {default = null;type = nullOr (ints.between 1 10);}`,
			assertError: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			option, err := nix.NewMkOptionFromSchema(tt.schema)
			tt.assertError(t, err)
			require.IsType(t, &nix.Int{}, option.Type)

			got := option.String()
			require.Equal(t, tt.want, got)
		})
	}
}
