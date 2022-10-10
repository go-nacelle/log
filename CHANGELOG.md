# Changelog

## [Unreleased]

## [v2.0.1] - 2022-10-10

### Added

- Added `WithLogger` and `FromContext`. [#8](https://github.com/go-nacelle/log/pull/8)

## [v2.0.0] - 2021-05-31

### Added

- Exposed the interface `MinimalLogger` and its constructor `FromMinimalLogger`. [#5](https://github.com/go-nacelle/log/pull/5)

### Removed

- Removed mocks package. [#6](https://github.com/go-nacelle/log/pull/6)

### Changed

- Renamed `ReplayAdapter` and `RollupAdapter` and to `ReplayLogger` and `RollupLogger`, respectively. [#5](https://github.com/go-nacelle/log/pull/5)

## [v1.1.2] - 2020-09-30

### Removed

- Removed dependency on [aphistic/sweet](https://github.com/aphistic/sweet) by rewriting tests to use [testify](https://github.com/stretchr/testify). [#3](https://github.com/go-nacelle/log/pull/3)

## [v1.1.1] - 2019-11-19

### Fixed

- Fixed bad console output. [db6e246](https://github.com/go-nacelle/log/commit/db6e24657334615a099e39bae0359179778016e4), [45875f1](https://github.com/go-nacelle/log/commit/45875f173a0db48fc3f615d96a4f83e015cdf130)

## [v1.1.0] - 2019-10-07

### Added

- Added `WithIndirectCaller` to control the number of stack frames to omit. [#2](https://github.com/go-nacelle/log/pull/2)

### Removed

- Removed dependency on [aphistic/gomol](https://github.com/aphistic/gomol) by rewriting base logger internally. [4e537aa](https://github.com/go-nacelle/log/commit/4e537aa0e5a08638bfb45f5153e8deccf6e1d00d)

### Changed

- Changed log field blacklist from a comma-separated list to a json-encoded array. [96b9d53](https://github.com/go-nacelle/log/commit/96b9d53baff25f7c0436799f520c3d4a5970941e)

## [v1.0.1] - 2019-06-20

### Added

- Added mocks package. [d24aad2](https://github.com/go-nacelle/log/commit/d24aad20df4c5b24dbdff3860c348af82abed169)

## [v1.0.0] - 2019-06-17

### Changed

- Migrated from [efritz/nacelle](https://github.com/efritz/nacelle).

[Unreleased]: https://github.com/go-nacelle/log/compare/v2.0.1...HEAD
[v1.0.0]: https://github.com/go-nacelle/log/releases/tag/v1.0.0
[v1.0.1]: https://github.com/go-nacelle/log/compare/v1.0.0...v1.0.1
[v1.1.0]: https://github.com/go-nacelle/log/compare/v1.0.1...v1.1.0
[v1.1.1]: https://github.com/go-nacelle/log/compare/v1.1.0...v1.1.1
[v1.1.2]: https://github.com/go-nacelle/log/compare/v1.1.1...v1.1.2
[v2.0.0]: https://github.com/go-nacelle/log/compare/v1.1.2...v2.0.0
[v2.0.1]: https://github.com/go-nacelle/log/compare/v2.0.0...v2.0.1
