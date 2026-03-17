GOCOVERDIR ?= build/coverage/integration
GOFLAGS += -covermode=atomic -race -shuffle=on -mod=readonly -trimpath

export GOCOVERDIR GOFLAGS

GOMODULE != go list -m

define GO_TEST_IGNORE_PATTERNS
$(strip
cmd/schema2nix/main.go:
.gen.go:
)
endef

# files used by go generate
define GO_GENERATE_SOURCES
$(strip
)
endef

# files produces by go generate
define GO_GENERATE_TARGETS
$(strip
)
endef

define TEST_PACKAGES
$(strip
)
endef

GO_SOURCES = $(shell git ls-files '*.go') $(strip ${GO_GENERATE_TARGETS}) nix/module.gotmpl

all: bin/schema2nix
all: gomod2nix.toml

check:: GOFLAGS=
check::
	@golangci-lint run

dist:: GOFLAGS=
dist::
	goreleaser release --clean --snapshot --skip archive,before,nfpm,sign --release-notes CHANGELOG.md

.PHONY: coverage
coverage: build/coverage/all.txt | build/coverage/
	go tool cover -func=$< | sed 's|${GOMODULE}/||g' | grep -v '100.0%' | column -t

test::
	go test -v -covermode=atomic -race -shuffle=on -mod=readonly -trimpath ./...

## make magic, not war ;)

build/coverage/all.txt: build/coverage/unit.part.txt build/coverage/integration.part.txt
	go tool gocovmerge $^ \
	| grep $(if ${GO_TEST_IGNORE_PATTERNS},-v '$(subst $(space),\|,${GO_TEST_IGNORE_PATTERNS})',.) \
	> $@

build/coverage/unit.part.txt: ${GO_SOURCES} go.sum | build/coverage/
	go test --coverprofile=$@ $(addprefix ./,$(addsuffix /...,${TEST_PACKAGES}))
	sed -i'' '#$(subst .,\.,$(subst $(space),\|,${GO_TEST_IGNORE_PATTERNS}))#d' $@

build/coverage/integration.part.txt: ${GO_SOURCES} go.sum | ${GOCOVERDIR}/
	-@rm -rf "${GOCOVERDIR}/*" 2>/dev/null || true
	go run -cover ./cmd/schema2nix/... .tmp/crush.schema.json > /dev/null 2>&1
	go tool covdata textfmt -i=${GOCOVERDIR} -o=$@ $(if ${TEST_PACKAGES},-pkg="$(addprefix ${GOMODULE}/,${TEST_PACKAGES})")
	sed -i'' '#$(subst .,\.,$(subst $(space),\|,${GO_TEST_IGNORE_PATTERNS}))#d' $@

bin/%: ${GO_SOURCES} go.sum
	go build -o ./$@ ./cmd/$(patsubst bin/%,%,$@)/...

go.sum: GOFLAGS-=-mod-readonly
go.sum: ${GO_SOURCES} go.mod
	@go mod tidy -v -x
	@touch $@

gomod2nix.toml: go.sum
	gomod2nix generate

build/coverage/%.html: build/coverage/%.txt | build/coverage/
	go tool cover -html=$< -o $@

${GO_GENERATE_TARGETS} &: ${GO_GENERATE_SOURCES}
	go generate ./...

hooks/pre-bump::
	goreleaser check
	goreleaser healthcheck

hooks/post-bump::
	env -u GOFLAGS op run -- goreleaser release --clean --release-notes CHANGELOG.md
