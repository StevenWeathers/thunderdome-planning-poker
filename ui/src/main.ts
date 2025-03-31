import App from './App.svelte';
import { mount } from 'svelte';

declare global {
  interface Window {
    app: any;
  }
}

const app = mount(App, {
  target: document.body,
});
window.app = app;

export default app;
