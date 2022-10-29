#!/usr/bin/env bash

# This script is meant to be used by go/template's maintainers
# to trigger the creation of a new release branch and PR according
# to the release flow described in the docs (../docs/release.md).

# The script sets the defined new version in the version.txt file
# and automatically opens a new PR on GitHub.

set -e

SVU="go run github.com/caarlos0/svu@latest"

echo "Current version is $($SVU current)"
read -p "Enter your intended release version: " VERSION

echo ""
echo "Preparing release $VERSION"
TRIMMED_VERSION="${VERSION#"v"}"
if [[ "$TRIMMED_VERSION" = "$VERSION" ]]; then
  echo "Version tags should start with a 'v'"
  exit 1
fi

BRANCH_NAME=release-"$TRIMMED_VERSION"
git checkout -b "$BRANCH_NAME"

echo "$TRIMMED_VERSION" >config/version.txt
git add config/version.txt
git commit -m "chore: prepare release $TRIMMED_VERSION"
git push --set-upstream origin "$BRANCH_NAME"

if command -v gh >/dev/null; then
  gh pr create --fill
else
  echo "GitHub CLI is not installed."
  echo "Pls open a PR for $BRANCH_NAME into main"
fi

git checkout main
