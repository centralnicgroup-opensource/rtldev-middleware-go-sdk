# [3.0.0](https://github.com/hexonet/go-sdk/compare/v2.3.0...v3.0.0) (2020-05-08)


### Features

* **apiclient:** automatic IDN conversion of API command parameters to punycode ([407c105](https://github.com/hexonet/go-sdk/commit/407c105d9d9f13a77fe68a9c1793596933edbd58))


### BREAKING CHANGES

* **apiclient:** Even though thought and build for internal purposes, we launch a major version for
this change. type of cmd parameter changes from map[string]inteface{} to map[string]string.

# [2.3.0](https://github.com/hexonet/go-sdk/compare/v2.2.3...v2.3.0) (2020-03-13)


### Features

* **apiclient:** support bulk parameter in api commands using slices ([c11db41](https://github.com/hexonet/go-sdk/commit/c11db411d22860929a12a4639f0b6422a95e1351))

## [2.2.3](https://github.com/hexonet/go-sdk/compare/v2.2.2...v2.2.3) (2019-10-04)


### Bug Fixes

* **responsetemplate/mgr:** improve description of `423 Empty API response` ([ce11490](https://github.com/hexonet/go-sdk/commit/ce11490))

## [2.2.2](https://github.com/hexonet/go-sdk/compare/v2.2.1...v2.2.2) (2019-09-19)


### Bug Fixes

* **release process:** migrate configuration ([a717401](https://github.com/hexonet/go-sdk/commit/a717401))

## [2.2.1](https://github.com/hexonet/go-sdk/compare/v2.2.0...v2.2.1) (2019-08-19)


### Bug Fixes

* **APIClient:** change default SDK url ([64d89d5](https://github.com/hexonet/go-sdk/commit/64d89d5))

# [2.2.0](https://github.com/hexonet/go-sdk/compare/v2.1.0...v2.2.0) (2019-04-17)


### Features

* **responsetemplate:** add IsPending method ([faa9c4d](https://github.com/hexonet/go-sdk/commit/faa9c4d))

# [2.1.0](https://github.com/hexonet/go-sdk/compare/v2.0.1...v2.1.0) (2019-04-03)


### Features

* **apiclient:** review user-agent header usage ([ed719e5](https://github.com/hexonet/go-sdk/commit/ed719e5))

## [2.0.1](https://github.com/hexonet/go-sdk/compare/v2.0.0...v2.0.1) (2018-11-12)


### Bug Fixes

* **pkg:** readd missing root-folder go file ([b4ffd6a](https://github.com/hexonet/go-sdk/commit/b4ffd6a))

# [2.0.0](https://github.com/hexonet/go-sdk/compare/v1.2.1...v2.0.0) (2018-11-12)


### Code Refactoring

* **pkg:** migration to generic cross-sdk structure; add CI/CD ([31778a1](https://github.com/hexonet/go-sdk/commit/31778a1))


### BREAKING CHANGES

* **pkg:** Downward incompatible, reviewed from scratch
