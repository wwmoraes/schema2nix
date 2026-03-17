// Package nix provides utilities to work with Nix source code.
package nix

import (
	_ "embed"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"text/template"
	"unicode"

	"github.com/google/jsonschema-go/jsonschema"
)

// ReservedIdentifierSymbols contains a regular expression to match reserved
// symbols that cannot go on a bare identifier in Nix i.e. identifiers with such
// characters must be quoted.
//
// TODO test reserved symbols compiles
const ReservedIdentifierSymbols = `[.$]`

var (
	// ErrEnumValueTypeMismatch occurs when unmarshaling an enum value that does not match its parent type
	ErrEnumValueTypeMismatch = errors.New("enum value type mismatch")
	// ErrTypeMismatch occurs when unmarshaling an incompatible type
	ErrTypeMismatch = errors.New("type mismatch")

	//go:embed module.gotmpl
	moduleTemplate []byte

	nixReservedWords = [...]string{
		"assert",
		"else",
		"if",
		"in",
		"let",
		"then",
		"with",
	}

	reservedIdentifierSymbolsRegexp = regexp.MustCompile(ReservedIdentifierSymbols)
)

// SchemaUnmarshaler represents any type-related values that can extract
// relevant information from a JSON schema.
type SchemaUnmarshaler interface {
	UnmarshalSchema(schema *jsonschema.Schema) error
}

// ModuleTemplate parses/loads the module template to use.
func ModuleTemplate() (*template.Template, int) {
	return template.Must(template.New("").Parse(string(moduleTemplate))), len(moduleTemplate)
}

// SafeExpression checks if an expression is safe to use as-is i.e. it doesn't
// have associativity issues. In practice it goes overly-simplistic and checks
// if there's any unicode spacing runes in the expression, wrapping it in
// parenthesis if so.
func SafeExpression(str string) string {
	if strings.ContainsFunc(str, unicode.IsSpace) {
		return fmt.Sprintf("(%s)", str)
	}

	return str
}

// SafeIdentifier checks if a given string is safe to use as an identifier,
// quoting it if necessary.
func SafeIdentifier(str string) string {
	index := slices.Index(nixReservedWords[:], str)
	if index != -1 && nixReservedWords[index] == str {
		return fmt.Sprintf("%q", str)
	}

	matched := reservedIdentifierSymbolsRegexp.MatchString(str)

	if matched {
		return fmt.Sprintf("%q", str)
	}

	return str
}
