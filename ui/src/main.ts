import App from './App.svelte';

declare global {
  interface Window {
    app: any;
  }
}

const app = new App({
  target: document.body,
});
window.app = app;

export default app;
