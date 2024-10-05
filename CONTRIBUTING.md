# Contributing

Volt is open-source, and **all contributions are welcome**.
Contributions can be related to:
- performance optimization
- new features
- bug fixes
- unit tests coverage
- documentation

## Pull Requests
The default branch: _master_

The PR are required to:
- contain a clear and detailed description
- add documentation (if necessary)
- not break the current performances
- pass the go linter
- match the unit tests coverage (or improve it) throuh codecov

The PR will be squashed merged once the CI is passed and code review is done.

### CI
Some [Github Actions](https://github.com/akmonengine/volt/tree/master/.github/workflows) are required:
- **Golangci-lint** should pass
- **Unit Tests** through codecov should pass


Other actions are not required but informative:
- **Goreportcard**

## Issues
Feel free to open [issues](https://github.com/akmonengine/volt/issues) for bugs or improvements.
Please attach a code sample to reproduce the bug. The resolution proposals (through PR) are very welcome.

