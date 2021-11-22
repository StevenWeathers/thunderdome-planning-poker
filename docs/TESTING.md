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
