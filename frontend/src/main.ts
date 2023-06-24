import App from './App.svelte'

declare global {
    let app: any;
}

window.app = new App({
    target: document.body,
})

export default app
