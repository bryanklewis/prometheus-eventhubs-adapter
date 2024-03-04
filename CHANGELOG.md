# Changelog
All notable changes to this project will be documented in this file. Format is based on [Keep a Changelog]( https://keepachangelog.com/en/1.0.0/ ).
This project adheres to [Semantic Versioning]( https://semver.org/ ).

Template:
```
## master - \[Unreleased\]
### Added
### Changed
### Deprecated
### Removed
### Fixed
```

## v0.5.4 - 04 March 2024
### Changed
- Bump crypto lib, thanks to @matthewhudsonedb

## v0.5.3 - 20 November 2023
### Changed
- Bump grpc lib, thanks to @matthewhudsonedb
- Bump all project libs

## v0.5.2 - 25 August 2023
### Changed
- Bump promhttp lib to fix resource overrun security issue, thanks to @matthewhudsonedb

## v0.5.1 - 11 July 2023
### Changed
- Bump go, gin-gonic and alpine to fix security vulnerabilities, thanks to @gbulloch-edb

## v0.5.0 - 31 January 2022
### Added
- partition key send support, thanks to @jotka

## v0.4.4 - 6 December 2021
### Changed
- Recreate hub on send error
- Bump golang and dependency versions

## v0.4.3 - 25 October 2021
### Changed
- Remove hardcoded listen_address in container

## v0.4.2 - 14 December 2020
### Changed
- build with golang 1.15.6
- use azure build agent{-latest}
- bump project dependencies
### Removed
- viper aliases in config

## v0.4.1 - 18 December 2019
### Changed
- build with golang 1.13.5
- use azure event hubs sdk v3
- bump project dependencies, except viper
### Removed
- dead code getCounterValue, part of throughput calc
- search for config in home directory
### Fixed
- linting simplify sample.Metric[]

## v0.3.3 - 14 October 2019
### Changed
- build with golang 1.13.1
- use azure event hubs sdk v2.0.3
### Removed
- throughput calc function
### Fixed
- debug message in hub package not accurate

## v0.3.2 - 20 August 2019
### Changed
- Use Alpine instead of scratch based Docker image
- Push to Docker Hub after build
- Prefix non-http metrics with 'adapter'
### Fixed
- Build ldflags not populating on linux

## v0.3.1 - 19 August 2019
### Added
- Docker image
### Fixed
- Build ldflags not populating on windows
- Documentation fixes

## v0.3.0 - 08 August 2019
### Added
- Public release
- Send batch events
- Azure Pipelines config

## v0.2.0 - 26 July 2019
### Added
- Avro JSON
### Changed
- Use viper library for config

## v0.1.0 - 12 July 2019
### Added
- Adapter created
