# Testing Thunderdome

In effort to maintain a stable application Thunderdome is supported by automated testing

## End to End Testing

[Cypress](https://www.cypress.io/) is used for End to End testing of Thunderdome, keep tests atomic!

```
cd e2e
npm install
npm run cypress:open
```

### Conventions for Writing E2E Tests

- Follow [cypress best practices](https://docs.cypress.io/guides/references/best-practices)
- Favor adding `data-testid"` attribute to dom elements for selection, never use css classes
    - custom cypress command `cy.getByTestId()` exists to select dom elements simply by passing the `data-testid`
      attribute value e.g. `cy.getByTestId('user-delete')` would look for `data-testid="user-delete"`

## Frontend Unit Testing

Frontend unit tests are run with [Jest](https://jestjs.io/) and utilize [@testing-library/svelte](https://testing-library.com/docs/svelte-testing-library/intro).

To run the tests `npm test` for a single run or `npm test:watch` to actively watch for test changes.

### Conventions for Frontend Unit Tests

- Check out [Unit testing svelte component](https://sveltesociety.dev/recipes/testing-and-debugging/unit-testing-svelte-component/) guide
- Test files are named `{componentName}.test.js`
- Test folders are named `__tests__` and live alongside the source code to which the test files are for
  - Example `frontend/src/pages/__tests__/*.test.js` corresponds to `frontend/src/pages/*.svelte`
- Snapshot testing should only be used for individual components not page components

## Go Unit Testing

Run `make testgo` to run go tests