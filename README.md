# go-repo-template

A simple template for Go projects.

## Utilities

### [CI](.github/workflows/ci.yml)

Builds, tests, lints and checks dependencies for known vulnerabilities on every push and pull request to `main`.

- **Tests** run with `-race` and `-shuffle=on`; coverage is uploaded as a workflow artifact.
- [golangci-lint](https://github.com/golangci/golangci-lint): lints using [.golangci.yml](.golangci.yml), pinned to the version in [.golangci-lint-version](.golangci-lint-version).
- [govulncheck](https://go.dev/doc/security/vuln/): reports known CVEs in dependencies that your code actually calls.
- [Codecov](https://codecov.io) (optional): uploads and visualizes test coverage. See [Optional utilities](#optional-utilities).

### [CodeQL](.github/workflows/codeql.yml)

Static analysis that finds vulnerable patterns in your own code (injection, unsafe path handling, etc.). Runs on pushes, pull requests and weekly.

### [Dependabot](.github/dependabot.yml)

Weekly version updates for Go modules and GitHub Actions (including SHA-pinned actions), with minor/patch updates grouped and a cooldown before new releases are proposed.

### [CodeRabbit](.coderabbit.yaml) (optional)

AI pull request review configured by [.coderabbit.yaml](.coderabbit.yaml).

### Other

- [.gitattributes](.gitattributes): forces LF line endings so formatting checks behave the same on every OS.

## Getting started

1. **Create a repository from this template**

   Use the `Use this template` button on GitHub, or the [GitHub CLI](https://cli.github.com):

   ```shell
   gh repo create <repo> --template ItakawaM/go-repo-template --public --clone
   ```

2. **Set your module path**

   ```shell
   go mod edit -module github.com/<user>/<repo>
   ```

3. **Update the [LICENSE](LICENSE)** with your name and the current year (or replace it).

4. **Configure repository settings**

   Settings are not copied from templates. Under `Settings -> Advanced Security`:

    - Enable `Dependabot alerts`
    - Enable `Dependabot security updates` (and optionally `Grouped security updates`)
    - Enable `Secret scanning` and `Push protection`
    - Leave `CodeQL default setup` off: this repository uses the advanced setup in [codeql.yml](.github/workflows/codeql.yml), and the two conflict

   Dependabot version updates need no setting; they are enabled by [dependabot.yml](.github/dependabot.yml).

## Optional utilities

### Codecov

1. Sign in to [Codecov](https://codecov.io) and add the repository.
2. Add the upload token as the `CODECOV_TOKEN` repository secret.
3. Uncomment the Codecov step in [ci.yml](.github/workflows/ci.yml) and adjust [codecov.yml](codecov.yml) if needed.

See the [Codecov quick start](https://docs.codecov.com/docs/quick-start).

### CodeRabbit

1. Install the CodeRabbit app for the repository by following the [quickstart](https://docs.coderabbit.ai/getting-started/quickstart).
2. Fill in the placeholder instructions in [.coderabbit.yaml](.coderabbit.yaml).