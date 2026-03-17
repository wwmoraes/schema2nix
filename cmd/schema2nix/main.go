// Binary schema2nix converts a JSON Schema definition into a Nix options module.
package main

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"

	"github.com/goccy/go-json"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/kr/pretty"
	"github.com/spf13/cobra"

	"github.com/wwmoraes/schema2nix"
)

const flagNameDebug = "debug"

var (
	exitCode atomic.Int32
	version  = "0.0.0"
)

func main() {
	defer func() {
		os.Exit(int(exitCode.Load()))
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	rootCmd := cobra.Command{
		Use:     "schema2nix [SCHEMA-FILE]",
		Args:    cobra.ExactArgs(1),
		RunE:    run,
		Version: version,
	}

	rootCmd.PersistentFlags().Bool(flagNameDebug, false, "prints debugging information on the standard error pipe")

	err := rootCmd.ExecuteContext(ctx)
	if err == nil {
		return
	}

	if errors.Is(err, context.Canceled) {
		exitCode.CompareAndSwap(0, 2)

		return
	}

	exitCode.CompareAndSwap(0, 1)
}

func run(cmd *cobra.Command, args []string) error {
	data, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	schema, err := parseSchema(data)
	if err != nil {
		return fmt.Errorf("parsing schema: %w", err)
	}

	debug, err := cmd.Flags().GetBool(flagNameDebug)
	if err != nil {
		return fmt.Errorf("getting %q flag value: %w", flagNameDebug, err)
	}

	if debug {
		fmt.Fprintln(os.Stderr, pretty.Sprint(schema))
	}

	nixSrc, err := schema2nix.Convert(schema)
	if err != nil {
		return fmt.Errorf("converting schema: %w", err)
	}

	fmt.Println(string(nixSrc))

	return nil
}

func parseSchema(data []byte) (*jsonschema.Schema, error) {
	var schema jsonschema.Schema

	err := json.Unmarshal(data, &schema)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling JSON: %w", err)
	}

	resolvedSchema, err := schema.Resolve(nil)
	if err != nil {
		return nil, fmt.Errorf("resolving schema: %w", err)
	}

	err = resolveDefs(&schema, &schema)
	if err != nil {
		return nil, fmt.Errorf("resolving defs: %w", err)
	}

	return resolvedSchema.Schema(), nil
}

//nolint:gocognit,gocyclo // TODO refactor
func resolveDefs(root, schema *jsonschema.Schema) error {
	if root == nil || schema == nil {
		return nil
	}

	if root.Defs == nil {
		root.Defs = make(map[string]*jsonschema.Schema)
	}

	maps.Copy(root.Defs, schema.Defs)

	for _, subSchema := range schema.Defs {
		err := resolveDefs(root, subSchema)
		if err != nil {
			return err
		}
	}

	if schema.Ref != "" {
		key, found := strings.CutPrefix(schema.Ref, "#/$defs/")
		if found {
			def, found := root.Defs[key]
			if found {
				*schema = *def
			}
		}
	}

	for _, subSchema := range schema.Properties {
		err := resolveDefs(root, subSchema)
		if err != nil {
			return err
		}
	}

	err := resolveDefs(root, schema.AdditionalItems)
	if err != nil {
		return err
	}

	err = resolveDefs(root, schema.AdditionalProperties)
	if err != nil {
		return err
	}

	for _, subSchema := range schema.AllOf {
		err := resolveDefs(root, subSchema)
		if err != nil {
			return err
		}
	}

	for _, subSchema := range schema.AnyOf {
		err := resolveDefs(root, subSchema)
		if err != nil {
			return err
		}
	}

	return nil
}
