# Changelog

## [1.2.0](https://github.com/anowarislam/ado/compare/v1.1.4...v1.2.0) (2025-11-26)


### Features

* **config:** add validate command ([#39](https://github.com/anowarislam/ado/issues/39)) ([8516448](https://github.com/anowarislam/ado/commit/8516448ba4a14c1031a6ab72983b76273b58ec4b))
* **logging:** implement structured logging infrastructure ([#42](https://github.com/anowarislam/ado/issues/42)) ([dc87e08](https://github.com/anowarislam/ado/commit/dc87e08b5fde6287f4ce9bf4893e3caccb318e7c))


### Documentation

* **adr:** 0001 - establish ADR + spec development workflow ([#31](https://github.com/anowarislam/ado/issues/31)) ([465af49](https://github.com/anowarislam/ado/commit/465af494a6b1e75a13329e08425f7352c69ddc7b))
* **adr:** 0002 - structured logging ([#40](https://github.com/anowarislam/ado/issues/40)) ([4c5d377](https://github.com/anowarislam/ado/commit/4c5d377c9434ae289f248828a1eaac04465f4e2d))
* integrate three-phase workflow across documentation ([#34](https://github.com/anowarislam/ado/issues/34)) ([75c396e](https://github.com/anowarislam/ado/commit/75c396e75acec3c4653375d9b71798ad2f4268b8))
* **spec:** command config validate ([#38](https://github.com/anowarislam/ado/issues/38)) ([968163c](https://github.com/anowarislam/ado/commit/968163cb2fe68bd38abc09b5426b40b0a0d71798))
* **spec:** establish specification framework for features and commands ([#33](https://github.com/anowarislam/ado/issues/33)) ([acc6723](https://github.com/anowarislam/ado/commit/acc6723afab3449f252856a7fd507c5299cf3aa1))
* **spec:** feature structured-logging ([#41](https://github.com/anowarislam/ado/issues/41)) ([334910f](https://github.com/anowarislam/ado/commit/334910f04adc3fa587914fb2b28620ec2c2c1511))

## [1.1.4](https://github.com/anowarislam/ado/compare/v1.1.3...v1.1.4) (2025-11-26)


### Bug Fixes

* **ci:** strip v prefix from tag for cosign signing ([#29](https://github.com/anowarislam/ado/issues/29)) ([1a18c1d](https://github.com/anowarislam/ado/commit/1a18c1d292bd801924cfbc3f8e223b63703d703d))

## [1.1.3](https://github.com/anowarislam/ado/compare/v1.1.2...v1.1.3) (2025-11-26)


### Bug Fixes

* **ci:** use GITHUB_TOKEN for Docker registry login ([#26](https://github.com/anowarislam/ado/issues/26)) ([619be65](https://github.com/anowarislam/ado/commit/619be6505106cb668518bafcad1ccce87046dcce))
* **ci:** use input tag for cosign signing in workflow_dispatch ([#28](https://github.com/anowarislam/ado/issues/28)) ([6357385](https://github.com/anowarislam/ado/commit/63573857c6b8cb991ef4f6f47921fbca3e6b1925))

## [1.1.2](https://github.com/anowarislam/ado/compare/v1.1.1...v1.1.2) (2025-11-26)


### Bug Fixes

* **docker:** remove zoneinfo copy from scratch image ([#23](https://github.com/anowarislam/ado/issues/23)) ([4773d43](https://github.com/anowarislam/ado/commit/4773d430e432e2701ae9d1dcedb9afed723f0519))

## [1.1.1](https://github.com/anowarislam/ado/compare/v1.1.0...v1.1.1) (2025-11-26)


### Bug Fixes

* **docker:** add extra_files to goreleaser docker config ([#21](https://github.com/anowarislam/ado/issues/21)) ([09f05f5](https://github.com/anowarislam/ado/commit/09f05f5d73967f91d567e1f8c4dd4d414b6e6c6d))

## [1.1.0](https://github.com/anowarislam/ado/compare/v1.0.2...v1.1.0) (2025-11-26)


### Features

* **ci:** add container pipeline for multi-arch Docker images ([#19](https://github.com/anowarislam/ado/issues/19)) ([a359a18](https://github.com/anowarislam/ado/commit/a359a18ab00808942c083d9c99471f4332a4c2f6))


### CI/CD

* add Claude Code GitHub workflows ([#17](https://github.com/anowarislam/ado/issues/17)) ([b5e445f](https://github.com/anowarislam/ado/commit/b5e445f907b3953243f565ed0d0fcabdc277f069))


### Documentation

* add comprehensive release automation guide ([00812f2](https://github.com/anowarislam/ado/commit/00812f24035b939dcfabeb3b3a1908c54c828b94))
* update RELEASE.md with PR workflow and branch protection ([#18](https://github.com/anowarislam/ado/issues/18)) ([59ec3a6](https://github.com/anowarislam/ado/commit/59ec3a6ea55b68bfc57efc1e8a4cb088c7e72434))

## [1.0.2](https://github.com/anowarislam/ado/compare/v1.0.1...v1.0.2) (2025-11-25)


### Bug Fixes

* **ci:** use GitHub App token for release creation ([2bde573](https://github.com/anowarislam/ado/commit/2bde573e4f3dd98b48d5d54b6725e2f45fdbbf5b))

## [1.0.1](https://github.com/anowarislam/ado/compare/v1.0.0...v1.0.1) (2025-11-25)


### Bug Fixes

* **config:** add release permissions ([#11](https://github.com/anowarislam/ado/issues/11)) ([cda5ebb](https://github.com/anowarislam/ado/commit/cda5ebb4991909b62c21104d877f13f55f2e7962))


### Dependencies

* **deps:** add dependabot grouping and auto-merge workflow ([710fea2](https://github.com/anowarislam/ado/commit/710fea26f9a671c47a640749e5104a4bd9ee5b4d))
* **deps:** bump the github-actions group with 4 updates ([#13](https://github.com/anowarislam/ado/issues/13)) ([05f4ae3](https://github.com/anowarislam/ado/commit/05f4ae34528b05cc7cdb3019167dcc66e1a257a6))
* **deps:** bump the python-dependencies group ([#14](https://github.com/anowarislam/ado/issues/14)) ([9329cc6](https://github.com/anowarislam/ado/commit/9329cc611e3e13021c3bd91113ced7b113952c8a))
* **release:** add changelog sections for build, ci, refactor, docs ([3c02686](https://github.com/anowarislam/ado/commit/3c026868297487546329267694e3c2ae863d5e25))


### Code Refactoring

* **ci:** separate goreleaser into own workflow ([c2af8ff](https://github.com/anowarislam/ado/commit/c2af8fff09adb0dc00c751f5efe22cb422ccc0c6))

## 1.0.0 (2025-11-25)


### Features

* initial project skeleton with Go CLI and Python lab ([4ebd438](https://github.com/anowarislam/ado/commit/4ebd438985db26886e8abcff1d70b3bb10c4a6dd))

## Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

This file is automatically updated by [release-please](https://github.com/googleapis/release-please).
