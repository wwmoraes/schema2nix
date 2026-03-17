#!/usr/bin/env sh

set -eu

export GIT_REFLOG_ACTION=pre-applypatch

exec git hook run --ignore-missing --to-stdin=/dev/stdin pre-commit -- "$@"
