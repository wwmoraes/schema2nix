package nix

import (
	"fmt"
)

// Option is a Nix/Nixpgs option definition that can unmarshal a JSON schema
// into itself and generate its own Nix code representation.
type Option interface {
	SchemaUnmarshaler
	fmt.Stringer

	// AsType returns the unwrapped type definition. It is useful for aggregate
	// types such as lists and for attribute sets at the root level where you need
	// the submodule definition without the mkOption wrapping.
	AsType() Type
}
