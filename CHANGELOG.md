# Changelog

## [0.0.11](https://github.com/baptistegh/go-lakekeeper/compare/v0.0.10...v0.0.11) (2025-07-21)


### ⚠ BREAKING CHANGES

* add explicit context argument to all API methods ([#92](https://github.com/baptistegh/go-lakekeeper/issues/92))

### Features

* add explicit context argument to all API methods ([#92](https://github.com/baptistegh/go-lakekeeper/issues/92)) ([7eb0818](https://github.com/baptistegh/go-lakekeeper/commit/7eb0818a1b6cfe90a766be3ad842ff8b1d5827a1))
* add integration with go-iceberg for catalog endpoints ([#89](https://github.com/baptistegh/go-lakekeeper/issues/89)) ([553afcb](https://github.com/baptistegh/go-lakekeeper/commit/553afcbfc4b30966ee0f4a5b1dd3be53e96d0ef2))
* **warehouse:** add deprecation notice for GetProtection ([#96](https://github.com/baptistegh/go-lakekeeper/issues/96)) ([df774ba](https://github.com/baptistegh/go-lakekeeper/commit/df774baaac5af01e8514d529523daddb00cd4835))
* **warehouse:** add few missing methods ([#94](https://github.com/baptistegh/go-lakekeeper/issues/94)) ([20e080b](https://github.com/baptistegh/go-lakekeeper/commit/20e080b70cd32600c4744711ce472f89447888c8))
* **warehouse:** add get statistics ([#95](https://github.com/baptistegh/go-lakekeeper/issues/95)) ([cc8ecff](https://github.com/baptistegh/go-lakekeeper/commit/cc8ecffc5a3ba428e8c81a91b1a1678c1aa80be2))
* **warehouse:** add GetNamespaceProtection ([#94](https://github.com/baptistegh/go-lakekeeper/issues/94)) ([20e080b](https://github.com/baptistegh/go-lakekeeper/commit/20e080b70cd32600c4744711ce472f89447888c8))
* **warehouse:** add GetTableProtection method ([#96](https://github.com/baptistegh/go-lakekeeper/issues/96)) ([df774ba](https://github.com/baptistegh/go-lakekeeper/commit/df774baaac5af01e8514d529523daddb00cd4835))
* **warehouse:** add GetViewProtection method ([#96](https://github.com/baptistegh/go-lakekeeper/issues/96)) ([df774ba](https://github.com/baptistegh/go-lakekeeper/commit/df774baaac5af01e8514d529523daddb00cd4835))
* **warehouse:** add ListSoftDeletedTabular ([#94](https://github.com/baptistegh/go-lakekeeper/issues/94)) ([20e080b](https://github.com/baptistegh/go-lakekeeper/commit/20e080b70cd32600c4744711ce472f89447888c8))
* **warehouse:** add SetNamespaceProtection ([#94](https://github.com/baptistegh/go-lakekeeper/issues/94)) ([20e080b](https://github.com/baptistegh/go-lakekeeper/commit/20e080b70cd32600c4744711ce472f89447888c8))
* **warehouse:** add SetTableProtection method ([#96](https://github.com/baptistegh/go-lakekeeper/issues/96)) ([df774ba](https://github.com/baptistegh/go-lakekeeper/commit/df774baaac5af01e8514d529523daddb00cd4835))
* **warehouse:** add SetViewProtection method ([#96](https://github.com/baptistegh/go-lakekeeper/issues/96)) ([df774ba](https://github.com/baptistegh/go-lakekeeper/commit/df774baaac5af01e8514d529523daddb00cd4835))
* **warehouse:** add table and view protection methods ([#96](https://github.com/baptistegh/go-lakekeeper/issues/96)) ([df774ba](https://github.com/baptistegh/go-lakekeeper/commit/df774baaac5af01e8514d529523daddb00cd4835))
* **warehouse:** add UndropTabular ([#94](https://github.com/baptistegh/go-lakekeeper/issues/94)) ([20e080b](https://github.com/baptistegh/go-lakekeeper/commit/20e080b70cd32600c4744711ce472f89447888c8))


### Miscellaneous Chores

* prepare release 0.0.11 ([afa161a](https://github.com/baptistegh/go-lakekeeper/commit/afa161a43e419f61143ef8c5e92c46035ae5d437))

## 0.0.10 (2025-07-19)

<!-- Release notes generated using configuration in .github/release.yml at main -->

## What's Changed
### 🎉 Features
* feat(permission): remove project scope on warehouse by @baptistegh in https://github.com/baptistegh/go-lakekeeper/pull/87


**Full Changelog**: https://github.com/baptistegh/go-lakekeeper/compare/v0.0.9...v0.0.10

## 0.0.9 (2025-07-18)

<!-- Release notes generated using configuration in .github/release.yml at main -->

## What's Changed
### 🎉 Features
* feat: add control on bootstrap user role by @baptistegh in https://github.com/baptistegh/go-lakekeeper/pull/82
* feat(permission): add warehouse interfaces by @baptistegh in https://github.com/baptistegh/go-lakekeeper/pull/85
* feat(permission): add missing GetAccess on role by @baptistegh in https://github.com/baptistegh/go-lakekeeper/pull/86
### Other Changes
* chore(ci): add v0.9.3 support by @baptistegh in https://github.com/baptistegh/go-lakekeeper/pull/80


**Full Changelog**: https://github.com/baptistegh/go-lakekeeper/compare/v0.0.8...v0.0.9

## 0.0.8 (2025-07-17)

<!-- Release notes generated using configuration in .github/release.yml at main -->

## What's Changed
### 🎉 Features
* feat(permission): add role interfaces by @baptistegh in https://github.com/baptistegh/go-lakekeeper/pull/78


**Full Changelog**: https://github.com/baptistegh/go-lakekeeper/compare/v0.0.7...v0.0.8

## 0.0.7 (2025-07-16)

<!-- Release notes generated using configuration in .github/release.yml at main -->

## What's Changed
### 🎉 Features
* feat(permission): implement server permissions interfaces by @baptistegh in https://github.com/baptistegh/go-lakekeeper/pull/52
* feat(permissions): add filtering support to server get access endpoint by @baptistegh in https://github.com/baptistegh/go-lakekeeper/pull/69
* feat(permission): add project interface support by @baptistegh in https://github.com/baptistegh/go-lakekeeper/pull/75
* feat(project): add get api statistics endpoint support by @baptistegh in https://github.com/baptistegh/go-lakekeeper/pull/70
### ✅ Bug Fixes
* fix(permission): rename all project related objects in server by @baptistegh in https://github.com/baptistegh/go-lakekeeper/pull/74
### 📚 Documentation
* chore: clean CHANGELOG.md by @baptistegh in https://github.com/baptistegh/go-lakekeeper/pull/50
### Other Changes
* chore: DRY in integration tests by @baptistegh in https://github.com/baptistegh/go-lakekeeper/pull/76


**Full Changelog**: https://github.com/baptistegh/go-lakekeeper/compare/v0.0.6...v0.0.7

## 0.0.6 (2025-07-15)

<!-- Release notes generated using configuration in .github/release.yml at main -->

## What's Changed
### Other Changes
* chore(release-please): fix previous tag by @baptistegh in https://github.com/baptistegh/go-lakekeeper/pull/46
* chore(release-please): rework v0.0.0 by @baptistegh in https://github.com/baptistegh/go-lakekeeper/pull/48


**Full Changelog**: https://github.com/baptistegh/go-lakekeeper/commits/v0.0.6
