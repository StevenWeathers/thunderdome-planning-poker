# UI AGENTS.md

Instructions for work under `ui/`. The [root AGENTS.md](../AGENTS.md) still applies; this file adds UI detail.

## Stack

Svelte 5, TypeScript, Vite, Tailwind CSS **3**, `typesafe-i18n`, `@lucide/svelte` icons, Node 26. This is a client-side
single page app served by the Go binary: no SvelteKit, SSR or file based routing.

## Commands

Run from `ui/`.

```bash
npm ci
npm run build                  # production build into dist/
npm run autobuild              # rebuild dist/ on change
npm run format                 # oxfmt write
npm run prettier               # oxfmt check (run by CI)
npm test                       # all tests with coverage
npx vitest run src/components/team/UpdateUser.test.ts
npx vitest run src/components/team/UpdateUser.test.ts -t "modal dialog"
npm run locales                # regenerate i18n types
npx playwright install --with-deps   # once, for the test browser
```

- Don't rely on `npm run dev`: `index.html` is a Go template that the backend fills in. Run the app through the backend
  (`task dev` from the root).
- With `serve --live` the backend reads `ui/dist` from disk, so pair it with `npm run autobuild` for UI iteration.
- `npm run check` currently reports many existing errors and is not a CI gate.

## Where things go

- `src/App.svelte`: router, locale setup, and the `xfetch` and `notifications` instances.
- `src/config.ts`: `appRoutes`, `AppConfig`, `PathPrefix`, supported locales.
- `src/pages/<domain>/`: route components. `src/components/<domain>/`: feature components.
- `src/components/{global,forms,table}/`: shared building blocks.
- `src/types/`: shared types. `src/stores.ts`: `user`, `activeAlerts`, `dismissedAlerts`, `dir`.
- `src/validationUtils.ts`, `src/dateUtils.ts`, `src/websocketUtil.ts`: check these before writing a helper.
- `public/`: static assets.

## Components

- Use `<script lang="ts">`, a `Props` interface with `$props()`, `$state`, `$derived`, `$effect`, `$bindable` and
  snippets (`{@render children?.()}`). Pass callbacks as props.
- Don't copy legacy `export let`, `on:` directives or `createEventDispatcher` patterns into new code.
- Check a shared component's `Props` before using it. Names vary: `SolidButton` takes `onClick` and `testid`, while
  native elements take `onclick` and `data-testid`.

| Need               | Use                                                                               |
|--------------------|-----------------------------------------------------------------------------------|
| Dialog             | `global/Modal` (`closeModal`, `ariaLabel`); it handles Escape and focus trapping  |
| Delete confirmation | `global/DeleteConfirmation`                                                      |
| Buttons            | `global/SolidButton`, `global/HollowButton`                                       |
| Inputs             | `forms/TextInput` (`bind:value`, exposes `focus()`), `SelectInput`, `Checkbox`, `Toggle` |
| Tables and CRUD    | `table/TableContainer`, `TableNav`, `Table`, `TableRow`, `RowCol`, `CrudActions`  |

Reference pair: [UserTitlesList](src/components/team/UserTitlesList.svelte) and
[UserTitleForm](src/components/team/UserTitleForm.svelte).

## Routing

- Add the path to `appRoutes` in `src/config.ts`, then register it with `router.on` in `App.svelte`, inside the matching
  feature flag block (`FeaturePoker`, `FeatureRetro`, ...).
- Pages receive `xfetch`, `router`, `notifications` and their route params as props. Pass these down; don't create new
  instances.
- Navigate with `router.route(...)`.
- Guard pages as neighbors do: check `$user.id`, then `validateUserIsRegistered` or `validateUserIsAdmin` as needed.
  These checks are for the interface only; the API must enforce authorization.
- Load pages with a literal dynamic `import()` passed to `showPage`. A static page import puts the page in the
  first load. Import a large dependency's CSS in the component that needs it.

## API calls

```ts
xfetch(`${apiPrefix}/user-titles`, { method: 'POST', body: { name } })
  .then(res => res.json())
  .then(result => {
    /* payload is in result.data */
  })
  .catch(error => {
    const status = error?.[1]?.status;
    notifications.danger($LL.userTitleSaveError());
  });
```

- The method defaults to `GET`, or `POST` when `body` is set. Pass the body as an object.
- Don't add `PathPrefix` to `xfetch` endpoints; it is added for you. A 401 logs the user out.
- Build team URLs for all three scopes (standalone, organization, department), as `pages/team/Team.svelte` does.

## Websockets

- Open connections with `Sockette` from `src/lib/sockette` and `getWebsocketAddress()`; close them in `onDestroy`.
- Messages are `{ type, value }`. Follow the existing `sendSocketEvent` and `onSocketMessage` handlers, including how
  each event encodes `value` as a string.
- Keep the existing handling of close codes `4001` to `4004`.
- Event names must match the Go handlers in `internal/http/<game>/`.

## Styling and accessibility

- Use Tailwind utilities with `dark:` variants. Dark mode is class based.
- Support right-to-left Persian: prefer logical utilities (`ms-`, `me-`, `ps-`, `pe-`, `start-`, `end-`).
- Label inputs, give dialogs and icon-only buttons an accessible name, and keep keyboard operation working.

## Localization

- Put all user facing text, including `aria-label`s and placeholders in new code, in i18n: `$LL.key()`.
- Add the key to all eight files, `src/i18n/{de,en,es,fa,fr,it,pt,ru}/index.ts`, with a real translation.
- Run `npm run locales`. Don't edit generated `src/i18n/i18n-*.ts` files.

## Security

- Use `{@html}` only for content the server sanitizes or for trusted markup. If you add a new rich text field,
  sanitize it in the Go data service.
- Scripts need the CSP nonce. Don't add inline scripts or external origins without a matching backend policy change.
- `vite.config.ts` and `index.html` contain Go template expressions (`{{.Nonce}}`, `{{.UIConfig...}}`); preserve them.

## Testing

- Name tests `Component.test.ts`, beside the component or in `__tests__/`. Tests run in real Chromium, not jsdom.
- Use `render` from `vitest-browser-svelte`, and `page` and `userEvent` from `vitest/browser`.
- Prefer label, role and test id locators; assert absence with `.query()` returning `null`.
- Pass `xfetch`, `notifications` and `router` as `vi.fn()` stubs.
- For socket driven pages, mock `lib/sockette` with `vi.mock` and call `flushSync()` after each message. See
  [Retro.test.ts](src/pages/retro/Retro.test.ts).
- English is preloaded by `vitest.setup.ts`.

Before finishing, run `npm test` and `npm run prettier`, and say what you could not verify. Visual changes need
screenshots in the pull request.

## Don't edit

`dist/`, `coverage/`, generated i18n files, and `package-lock.json` (except through npm). `src/lib/` is vendored; change it
only if the task requires it.
