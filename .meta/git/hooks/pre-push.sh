#!/usr/bin/env sh

set -u

export GIT_REFLOG_ACTION=pre-push

git stash push --keep-index --include-untracked | grep -vqFx "No local changes to save"
STASHED=$?

remake all test && cog check --from-latest-tag > /dev/null
STATUS=$?

if [ "${STASHED}" -eq "0" ]; then
  git stash pop --quiet || true
fi

exit ${STATUS}
