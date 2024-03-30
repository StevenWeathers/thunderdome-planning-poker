# Contributing to Thunderdome

We want to make contributing to this project as easy and transparent as possible, whether it's:

Reporting a bug
Discussing the current state of the code
Submitting a fix
Proposing new features
Becoming a maintainer

## Code of Conduct

Please read the [Code of Conduct](CODE_OF_CONDUCT.md) document.

## We Develop with GitHub

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

### Report bugs using Github's [issues](https://github.com/StevenWeathers/thunderdome-planning-poker/issues)

We use GitHub issues to track public bugs. Report a bug by [opening a new issue](); it's that easy!

### We Use [Github Flow](https://docs.github.com/en/get-started/quickstart/github-flow), So All Code Changes Happen Through Pull Requests

Pull requests are the best way to propose changes to the codebase (we
use [GitHub Flow](https://docs.github.com/en/get-started/quickstart/github-flow). We actively welcome your pull
requests:

1. Fork the repo and create your branch from `master`.
2. If you've added code that should be tested, add tests.
3. If you've changed/added APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue that pull request!

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see
the [tags on this repository](https://github.com/StevenWeathers/thunderdome-planning-poker/tags).

## Coding Conventions

- Thunderdome's priorities in architectural design are SaaS first, self-hosted second, always open-source.
    - Avoid vendor lock-in, use open standards such as open telemetry.
    - Keep the infrastructure requirements minimal, e.g. only requires Postgres, no cloud specific features.
    - Everything should be able to be bundled in the compiled Go binary including UI assets
        - This includes the UI assets and content

### Go code conventions

- Follow [Effective Go](https://go.dev/doc/effective_go)
  and [Code Review Comments Guide](https://go.dev/wiki/CodeReviewComments) from the Go project as much as
  possible within reason.
- Go is not an Object Oriented Programming language, we favor simplicity.
- Use standard library packages as much as possible, new dependencies should come with a valid reason for adding
  another dependency.
- All Go code is linted with `golangci-lint` on every commit.

### UI code conventions

- The Javascript framework is [Svelte](https://svelte.dev/) transpiled
  with [Typescript](https://www.typescriptlang.org/)
- CSS framework is [Tailwind](https://tailwindcss.com/)
- All UI code is linted with `npm run prettier` on every commit.

### End-to-End tests

End-to-End testing utilizes [Playwright](https://playwright.dev/) and aid in validating that Thunderdome's primary
features continue to work with every code change.

- End-to-End tests run on every pull-request commit.
- Write end-to-end tests that can be run in parallel.

## Authors

* **[Steven Weathers](https://github.com/StevenWeathers)** - *Creator and maintainer*

See also the list of [contributors](https://github.com/StevenWeathers/thunderdome-planning-poker/contributors) who
participated in
this project.

## Any contributions you make will be under the Apache 2.0 Software License

In short, when you submit code changes, your submissions are understood to be under the
same [Apache 2.0 License](http://www.apache.org/licenses/LICENSE-2.0) that covers the project. Feel free to contact the
maintainers if that's a concern.

## License

By contributing, you agree that your contributions will be licensed under Apache 2.0.
