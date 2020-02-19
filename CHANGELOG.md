# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [unreleased]

## [1.0.7] - 2020-02-19
### Changed
- enhance refanging by supporting `<DOT>` refanging

## [1.0.6] - 2019-09-23
### Fixed
- fixed an issue with defanging TLDs containing periods

### Changed
- enhance refanging by supporting additional schemes

## [1.0.5] - 2019-08-26
### Fixed
- fixed an issue with joeguo/tldextract

## [1.0.4] - 2019-07-11
### Fixed
- add support for `<.>` refanging
- add support for `<dot>` refanging
- fix issue where url was mangled when there was no scheme provided

## [1.0.3] - 2018-05-30
### Fixed
- add support for `[dot]` refanging

## [1.0.2] - 2018-05-08
### Fixed
- fixed defanging for IPv4 addresses. 

## [1.0.1] - 2018-05-03
### Changed
- remove trailing newline for clipboard output. 
- improve accuracy of IOC extraction. 

## [1.0.0] - 2018-05-02
- Initial Release


[unreleased]: https://github.com/jakewarren/defang/compare/v1.0.7...HEAD
[1.0.7]: https://github.com/jakewarren/defang/compare/v1.0.6...v1.0.7
[1.0.6]: https://github.com/jakewarren/defang/compare/v1.0.5...v1.0.6
[1.0.5]: https://github.com/jakewarren/defang/compare/v1.0.4...v1.0.5
[1.0.4]: https://github.com/jakewarren/defang/compare/v1.0.3...v1.0.4
[1.0.3]: https://github.com/jakewarren/defang/compare/v1.0.2...v1.0.3
[1.0.2]: https://github.com/jakewarren/defang/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/jakewarren/defang/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/jakewarren/defang/releases/tag/v1.0.0
