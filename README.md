# schema2nix

> converts JSON schemas to nix options

[![GitHub Issues](https://img.shields.io/github/issues/wwmoraes/schema2nix.svg)](https://github.com/wwmoraes/schema2nix/issues)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/wwmoraes/schema2nix.svg)](https://github.com/wwmoraes/schema2nix/pulls)

![Codecov](https://img.shields.io/codecov/c/github/wwmoraes/schema2nix)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

---

## 📝 Table of Contents

- [About](#about)
- [Usage](#usage)

## About

`schema2nix` parses JSON Schemas up to draft-07 and generates its equivalent
static options using Nixpkgs' `lib.types`. The output is a full module that you
can save to a file and import in your own definitions.

Its main use-case is to quickly add Nix typed settings for programs that offer a
JSON Schema specification, in particular for the Nixpkgs/nix-darwin/home-manager
projects.

### Inspiration

There's a few prior art attempts to bridge JSON schema and Nix options:

- [friedow's fromJsonSchema](https://github.com/friedow/fromJsonSchema), a pure Nix option that converts the schema at runtime using `builtins.fromJSON`
- [Lehmanator's json-schema-nix](https://github.com/Lehmanator/json-schema-nix), an abandoned idea to convert semantic/generic schemas
- [anna328p's nix-json-schema](https://github.com/anna328p/nix-json-schema), another converter abandoned idea

`fromJsonSchema` is a good option if you need only basic type conversion (i.e.
no enum/either) and a runtime conversion instead of an extra source file.


## Usage

`schema2nix` is meant to run as a pre-processor, generating an options module
that does not require the tool at runtime, while also allowing customization
of the final types. This also has the advantage that you won't slow down
evaluation as importing a file from your repository doesn't trigger an
[IFD](https://nix.dev/manual/nix/2.33/language/import-from-derivation).

That said, the flake in this repository provides the `mkOptionFromSchema`
function to run it at runtime by generating the schema in the Nix store, with
the drawback that it'll constitute an
[IFD](https://nix.dev/manual/nix/2.33/language/import-from-derivation).

> [!NOTE]
> The output is terse as the generator is purely functional (pun-intended),
> skipping any embelishments such as indentation or new lines. You can pipe
> `schema2nix` to `nixfmt` or another formatter of your preference to "prettify"
> it.

> [!TIP]
> You can convert schemas in YAML/TOML to JSON with tools such as
> [remarshal](https://github.com/remarshal-project/remarshal)

### pre-processor

Run the binary to generate the module:

```shell
schema2nix path/to/schema.json > options.nix
# OR through Nix
nix run github:wwmoraes/schema2nix -- path/to/schema.json > options.nix
```

It is provided in the overlay as a package called `schema2nix` so you can add
it to your devShell for instance.

You can then import the generated file to use as your options:

```nix
{
  pkgs,
  ...
}:
{
  options = {
    programs.foo = {
      settings = import ./options.nix { inherit (pkgs) lib; };
      };
    };
  };

  config = { ... };
}
```

### import from derivation

Add this repository as an overlay to your nixpkgs import, then you can use
the `mkOptionFromSchema` function, which outputs a nix plain-text file with an
importable module that declares the options. You can integrate it on your own
options definitions:

```nix
{
  pkgs,
  ...
}:
{
  options = {
    programs.foo = {
      settings = pkgs.mkOptionFromSchema {
        src = ./schema.json;
      };
    };
  };

  config = { ... };
}
```

The difference between `mkOptionFromSchema` and friedow's `fromJsonSchema` is
that the former is a derivation, that is, it runs the conversion only once for
each source version. It is also cacheable as any other nix derivation.

## Feature support

All JSON schema types are supported:

| JSON Schema | Nix               | Status |
|-------------|-------------------|--------|
| `array`     | `listOf anything` | ✅     |
| `boolean`   | `bool`            | ✅     |
| `integer`   | `int`             | ✅     |
| `null`      | `null`            | 🚧     |
| `number`    | `number`          | ✅     |
| `object`    | `attrs`           | ✅     |
| `string`    | `str`             | ✅     |

As a minimum `schema2nix` detects those types and sets the generated type to its
equivalent open-ended version in nix as seen above.

JSON Schema is powerful enough to specify constraints on top of plain types. In
that sense, `schema2nix` tries to identify constraints and generate constrained
types where possible:

| JSON Schema                                                | Nix                                                 | Status |
|------------------------------------------------------------|-----------------------------------------------------|--------|
| types with `enum`                                          | `enum [ ... ]`                                      | ✅     |
| `object` with declared properties                          | `submodule ...`                                     | ✅     |
| `object` with `additionalProperties` true/open-ended       | `submodule { freeformType = lazyAttrsOf raw; ... }` | ✅     |
| integers with minimum -128 and maximum 127                 | `ints.s8`                                           | 🚧     |
| integers with minimum -32768 and maximum 32767             | `ints.s16`                                          | 🚧     |
| integers with minimum -2147483648 and maximum 2147483647   | `ints.s32`                                          | 🚧     |
| integers with minimum 0                                    | `ints.unsigned`                                     | ✅     |
| integers with minimum 0 and maximum 255                    | `ints.u8`                                           | 🚧     |
| integers with minimum 0 and maximum 65535                  | `ints.u16`                                          | 🚧     |
| integers with minimum 0 and maximum 4294967295             | `ints.u32`                                          | 🚧     |
| integers with minimum 1                                    | `ints.positive`                                     | ✅     |
| integers with arbitrary minimum                            | `addCheck int (x: x >= ...)`                        | ✅     |
| integers with arbitrary minimum and maximum                | `ints.between`                                      | 🚧     |
| integers with arbitrary maximum                            | `addCheck int (x: x <= ...)`                        | 🚧     |
| numbers with arbitrary minimum and maximum                 | `numbers.between`                                   | 🚧     |
| numbers with minimum 0                                     | `numbers.nonnegative`                               | 🚧     |
| numbers with minimum 0 (exclusive)                         | `numbers.positive`                                  | 🚧     |

There's a few Nix types that have no clear way to apply yet:

| JSON Schema | Nix                          | Status |
|-------------|------------------------------|--------|
| ???         | `commas`                     | ⏳     |
| ???         | `envVar`                     | ⏳     |
| ???         | `ints.between`               | ⏳     |
| ???         | `lines`                      | ⏳     |
| ???         | `package`                    | ⏳     |
| ???         | `path`                       | ⏳     |
| ???         | `port` (alias of `ints.u16`) | ⏳     |
| ???         | `separatedString`            | ⏳     |
| ???         | `strMatching`                | ⏳     |
| ???         | `anything`                   | ⏳     |
| ???         | `unspecified`                | ⏳     |
| ???         | `boolByOr`                   | ⏳     |
| ???         | `float`                      | ⏳     |
| ???         | `nonEmptyStr`                | ⏳     |
| ???         | `singleLineStr`              | ⏳     |
| ???         | `passwdEntry`                | ⏳     |
| ???         | `fileset`                    | ⏳     |
| ???         | `shellPackage`               | ⏳     |
| ???         | `pkgs`                       | ⏳     |
| ???         | `pathInStore`                | ⏳     |
| ???         | `externalPath`               | ⏳     |
| ???         | `nonEmptyListOf`             | ⏳     |
| ???         | `attrsOf`                    | ⏳     |
| ???         | `lazyAttrsOf`                | ⏳     |
| ???         | `uniq`                       | ⏳     |
| ???         | `nullOf`                     | ⏳     |
| ???         | `deferredModule`             | ⏳     |
| ???         | `either`                     | ⏳     |
| ???         | `oneOf`                      | ⏳     |

Note: Nix types as per Nixpkgs' `lib.types`.

Legend:
- ✅: Fully supported
- 🚧: In construction
- ⏳: Pending
- ❌: Not supported (and won't be)
