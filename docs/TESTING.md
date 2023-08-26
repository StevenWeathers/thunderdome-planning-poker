# Testing Thunderdome

In effort to maintain a stable application Thunderdome is supported by automated testing

## End to End Testing

[Playwright](https://playwright.dev/) is used for End to End testing of Thunderdome, keep tests atomic!

```
cd e2e
npm install
npm test
```

## Frontend Unit Testing

Frontend unit tests are run with [Jest](https://jestjs.io/) and
utilize [@testing-library/svelte](https://testing-library.com/docs/svelte-testing-library/intro).

To run the tests `npm test` for a single run or `npm test:watch` to actively watch for test changes.

### Conventions for Frontend Unit Tests

- Check
  out [Unit testing svelte component](https://sveltesociety.dev/recipes/testing-and-debugging/unit-testing-svelte-component/)
  guide
- Test files are named `{componentName}.test.js`
- Test folders are named `__tests__` and live alongside the source code to which the test files are for
    - Example `ui/src/pages/__tests__/*.test.js` corresponds to `ui/src/pages/*.svelte`
- Snapshot testing should only be used for individual components not page components

## Go Unit Testing

Run `make testgo` to run go tests