# Changelog

## [0.1.0](https://github.com/baptistegh/go-lakekeeper/compare/github.com/baptistegh/go-lakekeeper-v0.0.5...github.com/baptistegh/go-lakekeeper-v0.1.0) (2025-07-15)


### ⚠ BREAKING CHANGES

* rename project Default method to GetDefault ([#21](https://github.com/baptistegh/go-lakekeeper/issues/21))
* create management module for related apis ([#20](https://github.com/baptistegh/go-lakekeeper/issues/20))
* **storage:** not sending errors back on storage creds/profile options func ([#16](https://github.com/baptistegh/go-lakekeeper/issues/16))
* init client structure ([#5](https://github.com/baptistegh/go-lakekeeper/issues/5))

### Features

* **auth:** add k8s service account token authentication ([#27](https://github.com/baptistegh/go-lakekeeper/issues/27)) ([ebe81b6](https://github.com/baptistegh/go-lakekeeper/commit/ebe81b628b92aea85ff53d1b35b71f22497d9b9d))
* **core:** add Ptr helper method ([#13](https://github.com/baptistegh/go-lakekeeper/issues/13)) ([afcc394](https://github.com/baptistegh/go-lakekeeper/commit/afcc39440a14a8f0d3aafac87ed31a260017c6ac))
* create management module for related apis ([#20](https://github.com/baptistegh/go-lakekeeper/issues/20)) ([996ddaf](https://github.com/baptistegh/go-lakekeeper/commit/996ddaf63405969c2f8394f987dd519e57fdae7e))
* init client structure ([#5](https://github.com/baptistegh/go-lakekeeper/issues/5)) ([370c4e3](https://github.com/baptistegh/go-lakekeeper/commit/370c4e3c9fdc123f1974dc9fd88dc58a013e6916))
* **project:** add missing methods DeleteDefault/RenameDefault ([#22](https://github.com/baptistegh/go-lakekeeper/issues/22)) ([9cd9be6](https://github.com/baptistegh/go-lakekeeper/commit/9cd9be6a4b74d6ae100a08110ce0b98fc68583a1))
* rename project Default method to GetDefault ([#21](https://github.com/baptistegh/go-lakekeeper/issues/21)) ([55fc6ec](https://github.com/baptistegh/go-lakekeeper/commit/55fc6ec57eb12cfc4a00d80d3266381c7cb6e50a))
* **role:** add search method ([#23](https://github.com/baptistegh/go-lakekeeper/issues/23)) ([fbc7dc4](https://github.com/baptistegh/go-lakekeeper/commit/fbc7dc48d050a643e0253c7a30d538d8d8cbe2cc))
* **storage:** not sending errors back on storage creds/profile options func ([#16](https://github.com/baptistegh/go-lakekeeper/issues/16)) ([ad0319c](https://github.com/baptistegh/go-lakekeeper/commit/ad0319cfa289824a23cbc9ee4b4ca0a208047884))
* **test:** add missing tests for warehouse service ([#15](https://github.com/baptistegh/go-lakekeeper/issues/15)) ([3dd130c](https://github.com/baptistegh/go-lakekeeper/commit/3dd130cfa28de17ac6fc761bf2739993a9b9e3ef))
* **test:** add unit tests ([#6](https://github.com/baptistegh/go-lakekeeper/issues/6)) ([e4be90a](https://github.com/baptistegh/go-lakekeeper/commit/e4be90a6bee68b4234558678db77428158c9b5e7))
* **test:** proposal unit tests structure ([#10](https://github.com/baptistegh/go-lakekeeper/issues/10)) ([127e624](https://github.com/baptistegh/go-lakekeeper/commit/127e6248e8f5e5f197ef21ada04035da3845ee6d))
* **user:** add search and list methods ([#25](https://github.com/baptistegh/go-lakekeeper/issues/25)) ([346bf56](https://github.com/baptistegh/go-lakekeeper/commit/346bf565bdca71e0385a8b5799ea0021b2ab8cf7))


### Bug Fixes

* **ci:** delete old release please config ([#34](https://github.com/baptistegh/go-lakekeeper/issues/34)) ([044a2de](https://github.com/baptistegh/go-lakekeeper/commit/044a2dee48a376f5f2714ac7615e8960e197aea3))
* ensure services are initialized with client ([#4](https://github.com/baptistegh/go-lakekeeper/issues/4)) ([5c9e0b0](https://github.com/baptistegh/go-lakekeeper/commit/5c9e0b02612d051fd8bf9501fd12f1366988c7db))
* temporary workaround for endpoints that does not use x-project-id ([#14](https://github.com/baptistegh/go-lakekeeper/issues/14)) ([4e2023b](https://github.com/baptistegh/go-lakekeeper/commit/4e2023b741c9be7ed015b5256cf64f3e20c4b5ed))
