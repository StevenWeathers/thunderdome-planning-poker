# Testing Thunderdome

In effort to maintain a stable application Thunderdome is supported by automated testing

## End to End Testing

[Playwright](https://playwright.dev/) is used for End to End testing of Thunderdome, keep tests atomic!

```bash
cd e2e
npm ci
npx playwright install --with-deps
npm test
```

## Frontend Unit Testing

Frontend unit tests are run with [Vitest](https://vitest.dev/) and use
[Playwright](https://playwright.dev/) browser automation for Svelte component testing via
[vitest-browser-svelte](https://github.com/sanderdl/vitest-browser-svelte).

If not already installed, install UI deps and Playwright

```bash
cd ui
npm ci
npx playwright install --with-deps
```

Then you can run tests

```bash
task test-ui
```

### Conventions for Frontend Unit Tests

- Refer to [vitest-browser-svelte documentation](https://github.com/sanderdl/vitest-browser-svelte) for Svelte component testing patterns
- Test files are named `{componentName}.test.ts`
- Test files can be placed either:
    - Alongside source files as `{componentName}.test.ts`
    - In `__tests__` folders: `ui/src/pages/__tests__/*.test.ts` corresponds to `ui/src/pages/*.svelte`
- Tests run in a Playwright browser environment (Chromium by default)

## Go Unit Testing

Run `task test-go` to run go tests