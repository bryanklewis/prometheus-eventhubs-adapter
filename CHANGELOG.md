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

## v0.4.0 - 18 December 2019
### Changed
- build with golang 1.13.5
- use azure event hubs sdk v3
- bump project dependencies
### Removed
- dead code getCounterValue, part of throughput calc
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
