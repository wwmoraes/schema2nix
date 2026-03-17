MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables
.SUFFIXES:

-include .env
-include .env.local
-include .env.secrets

export

-include $(shell git ls-files '**/*.mk')

.PHONY: all
all:
	@nix build --no-link --print-out-paths

.PHONY: check
check::
	@nix flake check --print-build-logs --verbose

.PHONY: clean
clean::
	-rm -rf bin build

.PHONY: dist
dist:: ;

.PHONY: release
release:
	cog bump --auto

.PHONY: test
test:: ;

## make magic, not war ;)

.tmp/crush.nix: .tmp/crush.schema.json bin/schema2nix
	./bin/schema2nix $< 2> .tmp/crush.log | tee .tmp/crush.raw.nix | nixfmt | tee /dev/tty > $@

.tmp/crush.schema.json:
	curl -fsSLo $@ https://charm.land/crush.json

%/:
	@test -d $@ || mkdir $@
