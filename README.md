## Semgrep Report to Github Comment script

This script is used to convert a semgrep report to a github markdown suitable for commenting on PRs.

__N.B This script is far from perfect and is a work in progress.__

## Installation

### Using Pre-compiled Binaries

You can download the correct version for your platform with this one-liner:

```bash
curl -L "https://github.com/rileywilliams/semgrep2github/releases/download/v0.1.2/semgrep2github_$(uname -s | tr '[:upper:]' '[:lower:]')_$([ "$(uname -m)" = "x86_64" ] && echo "amd64" || echo "arm64").tar.gz" | tar -xz
```

Or download pre-compiled binaries for your platform from the [releases page](https://github.com/rileywilliams/semgrep2github/releases/latest).

### Usage

```bash
semgrep2github -report report.json
```


### References

Report schema comes from https://github.com/semgrep/semgrep-interfaces/blob/main/semgrep_output_v1.jsonschema.

### Release Process

This project uses [GoReleaser](https://goreleaser.com/) to build and publish releases.

To create a new release:
```bash
git tag -a v0.x.x
git push origin v0.x.x
goreleaser release
```