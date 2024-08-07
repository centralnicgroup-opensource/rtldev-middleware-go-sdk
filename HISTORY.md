## [4.0.7](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v4.0.6...v4.0.7) (2024-08-07)


### Bug Fixes

* **deps:** bump golang.org/x/net from 0.27.0 to 0.28.0 ([81dca7d](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/81dca7dfb65225f46695e092b0bf9d08b7e3fd4c))

## [4.0.6](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v4.0.5...v4.0.6) (2024-07-25)


### Performance Improvements

* **idntranslator.go:** extend transitional processing TLDs array with '.art' TLD ([27ad3d2](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/27ad3d2db2cdb4a398e29379344fe13f3adde813))

## [4.0.5](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v4.0.4...v4.0.5) (2024-07-08)


### Bug Fixes

* **deps:** bump golang.org/x/net from 0.26.0 to 0.27.0 ([70a54f1](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/70a54f1b245dc8928348a006ce898279b248f5e4))

## [4.0.4](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v4.0.3...v4.0.4) (2024-06-05)


### Bug Fixes

* **deps:** bump golang.org/x/net from 0.25.0 to 0.26.0 ([a730f62](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/a730f62bdfc8d50750c111e0cde8c985e637f03b))
* **deps:** bump golang.org/x/text from 0.15.0 to 0.16.0 ([573cd71](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/573cd71d8aa06962327b605dfea100dc1f1b3e80))

## [4.0.3](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v4.0.2...v4.0.3) (2024-05-22)


### Bug Fixes

* **deps:** Bumped deps version to make go-sdk compatible with dns-control ([cc323a0](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/cc323a0c231da347b5fa527f8521e611c898f551))

## [4.0.3](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v4.0.2...v4.0.3) (2024-05-22)


### Bug Fixes

* **deps:** Bumped deps version to make go-sdk compatible with dns-control ([cc323a0](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/cc323a0c231da347b5fa527f8521e611c898f551))

## [4.0.2](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v4.0.1...v4.0.2) (2024-05-22)


### Bug Fixes

* **go.mod:** bumped go version 1.22.3 to 1.18 as dnscontrol supports 1.18 ([ed377c5](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/ed377c5375f64408f5914dddd62ae42534ddcd53))

## [4.0.1](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v4.0.0...v4.0.1) (2024-05-21)


### Bug Fixes

* **package.json:** reviewed package name from rtldev-middleware-go-sdk to go-sdk for go packagist ([eb75e17](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/eb75e17ae7137bf0313f9bbf2772b560e0a77937))

# [4.0.0](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v3.5.6...v4.0.0) (2024-05-21)


### Bug Fixes

* **idn translator:** replaced ([9fc05af](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/9fc05afb82b09d543e0d9eb5c8adfe55f0710683))
* **lib structure:** avoid import cycle and reported linter issues ([8832b05](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/8832b05b477beb3adbaa6a85c10076cfc3dbe10c))
* **response class:** merge with responsetemplate, add responsetranslator ([fef0a4b](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/fef0a4b5edc8181657e91e8c8c84401344a67b92))
* **response/-templatemanager:** review & patch failing tests ([fad8f86](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/fad8f8689b7d21f4b7abc3f84c23a02f75e33984))


### Performance Improvements

* **idntranslator & ResponseTranslator class:** deprecated API IDN Conversion, integrated GOLang IDN library with tests, and ResponseTranslator tests ([d6b42ef](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/d6b42ef3e2f7a363911412708d8eeade667a7b41))


### BREAKING CHANGES

* **response class:** Brought our library to the next golang level and applied a huge restructuring to be future safe

# [4.0.0](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v3.5.6...v4.0.0) (2024-05-21)


### Bug Fixes

* **idn translator:** replaced ([9fc05af](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/9fc05afb82b09d543e0d9eb5c8adfe55f0710683))
* **lib structure:** avoid import cycle and reported linter issues ([8832b05](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/8832b05b477beb3adbaa6a85c10076cfc3dbe10c))
* **response class:** merge with responsetemplate, add responsetranslator ([fef0a4b](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/fef0a4b5edc8181657e91e8c8c84401344a67b92))
* **response/-templatemanager:** review & patch failing tests ([fad8f86](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/fad8f8689b7d21f4b7abc3f84c23a02f75e33984))


### Performance Improvements

* **idntranslator & ResponseTranslator class:** deprecated API IDN Conversion, integrated GOLang IDN library with tests, and ResponseTranslator tests ([d6b42ef](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/d6b42ef3e2f7a363911412708d8eeade667a7b41))


### BREAKING CHANGES

* **response class:** Brought our library to the next golang level and applied a huge restructuring to be future safe

## [3.5.6](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v3.5.5...v3.5.6) (2023-11-30)


### Bug Fixes

* **apiclient.go:** patched an issue where autoconvertidn was overriding commands ([3c7ec16](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/3c7ec16306a59842a48c999d57b53ad596d1d1c8))

## [3.5.5](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v3.5.4...v3.5.5) (2023-01-17)


### Bug Fixes

* **new release:** for new namespace/module name in go.mod ([345318c](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/345318c3b88269512ed9eda32eed7cdcf2d5831a))

## [3.5.4](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v3.5.3...v3.5.4) (2022-06-24)

### Bug Fixes

- **apiclient:** add CONNECTION_URL to debug output in debug mode ([3caac20](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/3caac2095677ff22d3c543f5b2d5a3577a7b99eb))

## [3.5.3](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v3.5.2...v3.5.3) (2022-03-23)

### Bug Fixes

- **ot&e:** url updated for OT&E environment ([4aed2be](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/4aed2be94b4a81341d2bb18ca6939dd6b01dae84))

## [3.5.2](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v3.5.1...v3.5.2) (2021-04-09)

### Performance Improvements

- **fixed version:** included version number in go.mod ([903ca2f](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/903ca2f9b3065730cb19af4c7ac06e440b8655cb))

## [3.5.1](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v3.5.0...v3.5.1) (2021-01-21)

### Bug Fixes

- **ci:** migration from Travis CI to github actions ([3461c59](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/3461c59779134ef614e5a1599d2c13ccc1203343))

# [3.5.0](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v3.4.0...v3.5.0) (2020-05-11)

### Features

- **logger:** possibility to override debug mode's default logging mechanism. See README.md ([dc71ed9](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/dc71ed9417e838aae7c4e09834cd31e8f33764ef))

# [3.4.0](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v3.3.2...v3.4.0) (2020-05-08)

### Features

- **response:** possibility of placeholder vars in standard responses to improve error details ([87df76b](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/87df76b39b0e267f4acf12dcc695ba599e233bc4))

## [3.3.2](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v3.3.1...v3.3.2) (2020-05-08)

### Bug Fixes

- **security:** replace passwords whereever they could be used for output ([d698ab7](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/d698ab79af58216e5ae5bb8561b0c3b4bb1a796d))

## [3.3.1](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v3.3.0...v3.3.1) (2020-05-08)

### Bug Fixes

- **messaging:** return a specific error template in case code or description are missing ([faf78c4](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/faf78c413217c2b4c26632e08b497280c2a8c351))

# [3.3.0](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v3.2.0...v3.3.0) (2020-05-08)

### Features

- **apiclient:** allow to specify additional libraries vai SetUserAgent ([a440863](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/a44086372f9a0a1ad4d32671e98c1beab9dceb3b))

# [3.2.0](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v3.1.0...v3.2.0) (2020-05-08)

### Features

- **response:** added GetCommandPlain (getting used command in plain text) ([1e00417](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/1e00417222a37a2fc25d6e53e3224a3fdda4c950))

# [3.1.0](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v3.0.0...v3.1.0) (2020-05-08)

### Features

- **apiclient:** support the `High Performance Proxy Connection Setup`. see README.md ([3487c88](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/3487c8800001d9b790c0c398dbdcc3d78efc2863))

# [3.0.0](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v2.3.0...v3.0.0) (2020-05-08)

### Features

- **apiclient:** automatic IDN conversion of API command parameters to punycode ([407c105](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/407c105d9d9f13a77fe68a9c1793596933edbd58))

### BREAKING CHANGES

- **apiclient:** Even though thought and build for internal purposes, we launch a major version for
  this change. type of cmd parameter changes from map[string]inteface{} to map[string]string.

# [2.3.0](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v2.2.3...v2.3.0) (2020-03-13)

### Features

- **apiclient:** support bulk parameter in api commands using slices ([c11db41](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/c11db411d22860929a12a4639f0b6422a95e1351))

## [2.2.3](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v2.2.2...v2.2.3) (2019-10-04)

### Bug Fixes

- **responsetemplate/mgr:** improve description of `423 Empty API response` ([ce11490](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/ce11490))

## [2.2.2](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v2.2.1...v2.2.2) (2019-09-19)

### Bug Fixes

- **release process:** migrate configuration ([a717401](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/a717401))

## [2.2.1](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v2.2.0...v2.2.1) (2019-08-19)

### Bug Fixes

- **APIClient:** change default SDK url ([64d89d5](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/64d89d5))

# [2.2.0](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v2.1.0...v2.2.0) (2019-04-17)

### Features

- **responsetemplate:** add IsPending method ([faa9c4d](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/faa9c4d))

# [2.1.0](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v2.0.1...v2.1.0) (2019-04-03)

### Features

- **apiclient:** review user-agent header usage ([ed719e5](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/ed719e5))

## [2.0.1](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v2.0.0...v2.0.1) (2018-11-12)

### Bug Fixes

- **pkg:** readd missing root-folder go file ([b4ffd6a](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/b4ffd6a))

# [2.0.0](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/compare/v1.2.1...v2.0.0) (2018-11-12)

### Code Refactoring

- **pkg:** migration to generic cross-sdk structure; add CI/CD ([31778a1](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/commit/31778a1))

### BREAKING CHANGES

- **pkg:** Downward incompatible, reviewed from scratch
