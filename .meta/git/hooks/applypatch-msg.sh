#!/usr/bin/env sh

set -eu

export GIT_REFLOG_ACTION=applypatch-msg

exec git hook run --ignore-missing --to-stdin=/dev/stdin commit-msg -- "$@"
