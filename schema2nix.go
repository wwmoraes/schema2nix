// Package schema2nix provides utilities to convert type definitions in JSON
// Schema to its equivalent in Nix/Nixpkgs options.
package schema2nix

import (
	"bytes"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"

	"github.com/wwmoraes/schema2nix/nix"
)

// Convert parses a JSON schema into a Nix option and marshals it to its textual
// representation. The result is valid Nix module that defines the option with
// either Nixpkgs' lib.types.submodule for objects or lib.options.mkOption for
// all other types.
func Convert(schema *jsonschema.Schema) ([]byte, error) {
	tmpl, tmplLength := nix.ModuleTemplate()

	option, err := nix.NewMkOptionFromSchema(schema)
	if err != nil {
		return nil, fmt.Errorf("parsing root schema: %w", err)
	}

	var data string
	if schema.Type == "object" {
		data = option.AsType().String()
	} else {
		data = option.String()
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(data)+tmplLength))

	err = tmpl.Execute(buf, data)
	if err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return buf.Bytes(), nil
}
