# Project Instructions

## Project Overview

Go SDK for the CentralNic Reseller (fka RRPProxy) backend API, successor to the deprecated
Hexonet Go SDK. Node is present only for the release tooling, never for product code.

- **Language / runtime:** Go (`go 1.26.0` in `go.mod`)
- **Module path:** `github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5`
- **Default branch:** `master`

## Architecture

- **The module path carries the major version (`/v5`).** semantic-release rewrites it in
  `go.mod` and every `*.go` import on a major bump — never change it by hand.
- **`GetVersion()` in `apiclient/apiclient.go` is written by the release.** Its version
  string is replaced from `.releaserc.json`; do not edit it manually.
- **`go.mod` versions are chosen for dnscontrol compatibility** (see the comments there) —
  check dnscontrol before raising the `go` directive or a requirement.

## Coding Standards

- **Static analysis:** golangci-lint v2, config `.github/linters/.golangci.yml` — `standard`
  plus `gosec`, `unconvert`, `gocyclo` (min-complexity 20), `goconst`, `gocritic`, `govet`
  (with `shadow`), `revive`; `errcheck` with `check-blank: true`.
- **Formatting:** `goimports` (enabled as the golangci-lint formatter).
- Lint issues are reported only for changes since `origin/master` (`new-from-rev`) — do not
  read a clean run as a clean codebase.

## Testing

- **Framework:** Go `testing`; `github.com/stretchr/testify/assert` only where a test file
  already uses it (`apiclient`, `idntranslator`). Tests sit beside the source as
  `<package>/<package>_test.go`.
- **No real API calls in unit tests.** Responses come from `rtm.AddTemplate(...)` response
  templates and an `httptest` server (see `apiclient/apiclient_test.go`).
- The demo (`demo/demo.go`) does call the live API, with `CNR_TEST_USER` /
  `CNR_TEST_PASSWORD` — it is not a test.

### Running tests

```sh
./scripts/test-go.sh    # go test -v -race -coverprofile=coverage.out ./... + HTML report
pnpm run test-demo      # go mod tidy && go run demo/demo.go — live API, needs credentials
```

golangci-lint runs only in the shared CI workflow, against `.github/linters/.golangci.yml`;
the devcontainer does not install it.

## Build, CI & Policies

- **CI** delegates to
  [rtldev-middleware-shareable-workflows](https://github.com/centralnicgroup-opensource/rtldev-middleware-shareable-workflows)
  (`go-sdk-test.yml` on pull requests, `go-sdk-release.yml` on push to `master`). Do not
  reimplement them in this repository's workflow stubs.
- **Releases:** semantic-release (`.releaserc.json`) writes the changelog to `HISTORY.md`
  and commits it back to `master`.
- **Lockfiles:** `pnpm-lock.yaml` and `go.sum` are committed deliberately; `package-lock.json`
  is git-ignored — pnpm only.
- **Supply chain:** `minimumReleaseAge: 1440` in `pnpm-workspace.yaml` is enforced
  workspace-wide — wait out a red build rather than lowering it.
- **Devcontainer:** the frame is in `.devcontainer/`; the shared behaviour comes from the
  `devbase` devcontainer Feature (`devbase:2`) and Go from the `go` Feature. Never fork
  devbase's scripts here — repository-specific setup goes in `.devcontainer/post-create.sh`.
- **Runtime version policy:** Go from the `go` directive in `go.mod`; Node `engines` and
  `devEngines.packageManager` (pnpm `^11.0.0`) in `package.json` cover the release tooling.

## Git Conventions

- **Commit messages:** Conventional Commits with **mandatory scope** —
  `<type>(<scope>): <summary>`. Never append a `Co-Authored-By:` trailer.
- **Commit type selection:** `fix`/`feat` are reserved for source changes — they trigger a
  release. Everything else uses a non-releasing type: `ci` (workflows, devcontainer),
  `build` (build tooling), `chore`, `docs`, `test`, `refactor`.
- **Breaking changes:** add `BREAKING CHANGE: <summary>` to the commit body after a blank
  line, and document the consumer upgrade path in the same change — it bumps the `/vN`
  module path.
- **Branch creation:** before creating a branch, and again after any PR merges,
  run `git checkout master && git pull --ff-only` first — it brings in the
  merged commits and their release tags. Never branch from a stale local
  `master`, from `origin/master` (it becomes the new branch's upstream) or from
  another feature branch.
- **Branch naming:** prefix with the Jira issue ID — `RSRMID-1234-short-description`.
- **Pull requests:** include the Jira issue link in the description; after opening, add the
  PR URL as a comment on the Jira issue.
- **Merging:** rebase-merge (`gh pr merge --rebase`). Squash merges are disabled at the
  repository level.

## Atlassian / JIRA

Work is tracked in **Jira Cloud**, project `RSRMID` — not GitHub Issues. Load the
`jira-conventions` skill before creating, editing or transitioning an issue: it carries the
required fields, this repository's title prefix and component, and the worklog rule for Done and Cancelled.

## Do NOT

- Add dependencies without explicit request
- Add `Co-Authored-By:` trailers to commit messages
- Fork the `devbase` devcontainer Feature's scripts into this repository
- Hand-edit `HISTORY.md`, the `GetVersion()` string or the `/vN` module path — the release
  owns them
- Commit `.env` or display its contents — it is git-ignored for local environment
  variables such as credentials
