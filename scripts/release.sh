#!/bin/bash

DEFAULT_BRANCH='main'
CURRENT_GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

SRCDIR=${SRCDIR:-"."}
GOBIN=$(go env GOPATH)/bin
VER_PACKAGE="internal/version"
VER_CMD=${GOBIN}/svu
VER_BUMP=${GOBIN}/gobump
CHANGELOG_CMD=${GOBIN}/git-chglog
REL_CMD=${GOBIN}/goreleaser
RELEASE_NOTES_FILE=${SRCDIR}/tmp/relnotes.md

# Compare versions
VER_CURR=$(${VER_CMD} current --strip-prefix)
VER_NEXT=$(${VER_CMD} next --strip-prefix)

if [ $CURRENT_GIT_BRANCH != ${DEFAULT_BRANCH} ]; then
  echo "Not on ${DEFAULT_BRANCH}, skipping"
  exit 0
fi

if [ "${VER_CURR}" = "${VER_NEXT}" ]; then
  echo "No new version recommended, skipping"
  exit 0
fi

GIT_USER=$(git config user.name)
GIT_EMAIL=$(git config user.email)

if [ -z "${GIT_USER}" ]; then
  echo "git user.name not set"
  exit 1
fi
if [ -z "${GIT_EMAIL}" ]; then
  echo "git user.email not set"
  exit 1
fi

echo "Generating release for v${VER_NEXT}"

# Bump package version
${VER_BUMP} set ${VER_NEXT} -r -w ${VER_PACKAGE}

# Auto-generate CHANGELOG updates
${CHANGELOG_CMD} --next-tag v${VER_NEXT} -o CHANGELOG.md

# Commit CHANGELOG updates
git add CHANGELOG.md ${VER_PACKAGE}/

git commit --no-verify -m "chore(release): Releasing v${VER_NEXT}"
git tag v${VER_NEXT}
git push --no-verify origin HEAD:${DEFAULT_BRANCH} --tags

# Make release notes
${CHANGELOG_CMD} --silent -o ${RELEASE_NOTES_FILE} v${VER_NEXT}

# Publish the release
${REL_CMD} release --release-notes=${RELEASE_NOTES_FILE}
