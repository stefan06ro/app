# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

- Add values service extracted from app-operator.

### Added

- Add annotation and key packages extracted from app-operator.

### Changed

- Updated apiextensions to v3.4.0.
- Prepare module v3.

## [2.0.0] - 2020-08-11

### Changed

- Updated Kubernetes dependencies to v1.18.5.

## [0.2.3] - 2020-06-23

### Changed

- Update apiextensions to avoid displaying empty strings in app CRs.

## [0.2.2] - 2020-06-01

### Changed

- Set version label value to 0.0.0 so control plane app CRs are reconciled by
  app-operator-unique.

## [0.2.1] - 2020-04-24

- Fix module path (was accidentaly declared as gitlab.com/...).

## [0.2.0] - 2020-04-24

### Changed

- migrate from dep to go modules (build-only changes)

## [0.1.0] - 2020-04-24

### Added

- First release

[Unreleased]: https://github.com/giantswarm/app/compare/v2.0.0...HEAD
[2.0.0]: https://github.com/giantswarm/app/compare/v0.2.3...v2.0.0
[0.2.3]: https://github.com/giantswarm/app/compare/v0.2.2...v0.2.3
[0.2.2]: https://github.com/giantswarm/app/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/giantswarm/app/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/giantswarm/app/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/giantswarm/app/releases/tag/v0.1.0
