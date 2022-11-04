# Releasing a new version of go-template

Since we want to support that even users that install `go-template` with `go install` get the right version output of `go version`
we currently embed a version file in this repo at [`config/version.txt`](../config/version.txt).
This file has to be updated on every release.

To keep the process as easy as possible some helper scripts and workflows have been created to ensure the correct release process.

To create a new release please do the following:

- Checkout the main branch
- `git pull` to update it
- Execute `make release` to create a new release PR (if you don't have the GitHub CLI installed you need to create the PR yourself from the branch that was created by the script)
- Get someone else to approve your PR. This also ensures that all maintainers agree that a new version should be released
- After the PR is merged, the workflow `tag-release` should be executed automatically
  - This workflow listens to all changes on the `config/version.txt` file
  - It creates a new tag on the main branch using the version defined in the file
- The newly created tag will then trigger the `release` workflow automatically which creates a new release
