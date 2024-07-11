#!/bin/bash


DEFAULT_BRANCH='main'
CURRENT_GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

SRCDIR=${SRCDIR:-"."}
GOBIN=$(go env GOPATH)/bin
VER_PACKAGE="internal/version"
VER_CMD=${GOBIN}/svu
VER_BUMP=${GOBIN}/gobump
CHANGELOG_CMD=${GOBIN}/git-chglog
CHANGELOG_FILE=CHANGELOG.md
REL_CMD=${GOBIN}/goreleaser
RELEASE_NOTES_FILE=${SRCDIR}/tmp/relnotes.md
SPELL_CMD=${GOBIN}/misspell

if [ $CURRENT_GIT_BRANCH != ${DEFAULT_BRANCH} ]; then
  echo "Not on ${DEFAULT_BRANCH}, skipping"
  exit 0
fi

# Compare versions
VER_CURR=$(${VER_CMD} current --strip-prefix)
VER_NEXT=$(${VER_CMD} next --strip-prefix)

echo " "
echo "Comparing tag versions..."
echo "Current version: ${VER_CURR}"
echo "Next version:    ${VER_NEXT}"
echo " "

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

echo "Generating release for v${VER_NEXT} with git user ${GIT_USER}"

# Update package version in version.go file using svu
${VER_BUMP} set ${VER_NEXT} -r -w ${VER_PACKAGE}

# Auto-generate CHANGELOG updates
${CHANGELOG_CMD} --next-tag v${VER_NEXT} -o ${CHANGELOG_FILE}
${SPELL_CMD}  -source text -w ${CHANGELOG_FILE}

# Commit CHANGELOG updates
git add CHANGELOG.md ${VER_PACKAGE}

git commit --no-verify -m "chore(release): release v${VER_NEXT}"
git push --no-verify origin HEAD:${DEFAULT_BRANCH}

if [ $? -ne 0 ]; then
  echo "Failed to push branch updates, exiting"
  exit 1
fi

# Tag and push
git tag v${VER_NEXT}
git push --no-verify origin HEAD:${DEFAULT_BRANCH} --tags

if [ $? -ne 0 ]; then
  echo "Failed to push tag, exiting."
  exit 1
fi

# Generate release notes for GoReleaser to add to the GitHub release description
${CHANGELOG_CMD} -o ${RELEASE_NOTES_FILE} ${VER_NEXT} --sort semver

# Correct spelling mistakes in release notes
${SPELL_CMD} -source text -w ${RELEASE_NOTES_FILE}

# Publish the release
${REL_CMD} release --release-notes=${RELEASE_NOTES_FILE}
