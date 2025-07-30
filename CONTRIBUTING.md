# Contributing to go-lakekeeper

Thank you for your interest in contributing to **go-lakekeeper**! 🎉  
We welcome issues, bug fixes, new features, and documentation improvements.

Please take a moment to review this guide before submitting a pull request.

## Table of Contents

- [Contributing to go-lakekeeper](#contributing-to-go-lakekeeper)
  - [Table of Contents](#table-of-contents)
  - [Code of Conduct](#code-of-conduct)
  - [Getting Started](#getting-started)
  - [How to Contribute](#how-to-contribute)
    - [Reporting Bugs](#reporting-bugs)
    - [Suggesting Features](#suggesting-features)
    - [Submitting Changes](#submitting-changes)
  - [Coding Guidelines](#coding-guidelines)
  - [Testing](#testing)
  - [Pull Request Process](#pull-request-process)
  - [License](#license)

---

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md) to help create a welcoming and respectful community.

## Getting Started

1. Fork the repository and clone it locally.
2. Install [Go 1.24+](https://golang.org/dl/).
3. Run `go mod tidy` to install dependencies.
4. Run `make build` or use your preferred method to build the CLI/client.

## How to Contribute

### Reporting Bugs

- Use [GitHub Issues](https://github.com/baptistegh/go-lakekeeper/issues).
- Include a clear title and description.
- Include steps to reproduce, expected behavior, and actual behavior.

### Suggesting Features

- Open an issue
- Provide context and possible use cases.

### Submitting Changes

- Fork the repo and create a new branch: `git checkout -b feature/my-feature`
- Make your changes
- Format your code: `make fmt`
- Run tests: `make test`
- Commit and push your branch
- Open a pull request and follow the PR template

## Coding Guidelines

- Use `gofmt`, `golangci-lint`, and `go vet` before submitting.
- Keep your changes focused: small, atomic PRs are easier to review.
- Use clear and descriptive commit messages.
- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines.

## Testing

Run all tests before opening a PR:

```sh
# Unit tests
make test

# Integration tests
make test-integration
```

If adding a feature or fixing a bug, include relevant unit and/or integration tests.

## Pull Request Process

1. Ensure all CI checks pass (unit tests, linters, etc.)
2. Make sure your branch is up to date with `main`
3. Link related issues in the PR description
4. A maintainer will review your PR and may request changes
5. Once approved, your PR will be merged into `main`

## License

By contributing, you agree that your contributions will be licensed under the [Apache 2.0 License](LICENSE).

---

Thanks again for your support and contributions!
Let's build something great together.