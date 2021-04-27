#!/bin/bash

COLOR_RED='\033[0;31m'
COLOR_NONE='\033[0m'
CURRENT_GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

if [ $CURRENT_GIT_BRANCH != 'main' ]; then
  printf "\n"
  printf "${COLOR_RED} Error: The release.sh script must be run while on the main branch. \n ${COLOR_NONE}"
  printf "\n"

  exit 1
fi

if [ $# -ne 1 ]; then
  printf "\n"
  printf "${COLOR_RED} Error: Release version argument required. \n\n ${COLOR_NONE}"
  printf " Example: \n\n    ./tools/release.sh 0.9.0 \n\n"
  printf "  Example (make): \n\n    make release version=0.9.0 \n"
  printf "\n"

  exit 1
fi

# Ensure there is no leading 'v'
RELEASE_VERSION=$(echo $1 | sed -e '/^v/s/^v\(.*\)$/\1/g')
GIT_USER=$(git config user.email)

echo "Generating release for v${RELEASE_VERSION} using system user git user ${GIT_USER}"

git checkout -b release/v${RELEASE_VERSION}

# Auto-generate CHANGELOG updates
git-chglog --next-tag v${RELEASE_VERSION} -o CHANGELOG.md

# Update version in version.go file
echo -e "package version\n\n// Version of this library\nconst Version string = \"${RELEASE_VERSION}\"" > internal/version/version.go

# Commit CHANGELOG updates
git add CHANGELOG.md internal/version/version.go
git commit --no-verify -m "chore(changelog): Update CHANGELOG for v${RELEASE_VERSION}"
git push --no-verify origin release/v${RELEASE_VERSION}

