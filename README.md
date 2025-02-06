## Semgrep Report to Github Comment script

This script is used to convert a semgrep report to a github markdown suitable for commenting on PRs.

__N.B This script is far from perfect and is a work in progress.__

### Usage

```bash
semgrep2github -report report.json
```


### References

Report schema comes from https://github.com/semgrep/semgrep-interfaces/blob/main/semgrep_output_v1.jsonschema.