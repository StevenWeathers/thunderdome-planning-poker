<script>
    import { AppConfig } from '../config.js'

    let count = 0
    let defaultTimeout = AppConfig.ToastTimeout
    let toasts = []
    let themes = {
        danger: '#bb2124',
        success: '#22bb33',
        warning: '#f0ad4e',
        info: '#5bc0de',
        default: '#aaaaaa',
    }

    function animateOut(node, { delay = 0, duration = 300 }) {
        function vhTOpx(value) {
            var w = window,
                d = document,
                e = d.documentElement,
                g = d.getElementsByTagName('body')[0],
                x = w.innerWidth || e.clientWidth || g.clientWidth,
                y = w.innerHeight || e.clientHeight || g.clientHeight
            return (y * value) / 100
        }

        return {
            delay,
            duration,
            css: t =>
                `opacity: ${
                    (t - 0.5) * 1
                }; transform-origin: top right; transform: scaleX(${
                    (t - 0.5) * 1
                });`,
        }
    }

    function createToast(msg, theme, timeout) {
        const background = themes[theme] || themes['default']

        timeout = timeout || defaultTimeout

        toasts.unshift({
            id: count,
            msg,
            background,
            timeout,
            width: '100%',
        })
        toasts = toasts
        count = count + 1
    }

    export function removeToast(id) {
        toasts = toasts.filter(t => t.id != id)
    }

    export function show(msg, timeout, theme = 'default') {
        createToast(msg, theme, timeout)
    }

    export function danger(msg, timeout) {
        show(msg, timeout, 'danger')
    }

    export function warning(msg, timeout) {
        show(msg, timeout, 'warning')
    }

    export function info(msg, timeout) {
        show(msg, timeout, 'info')
    }

    export function success(msg, timeout) {
        show(msg, timeout, 'success')
    }
</script>

<style>
    .toasts {
        list-style: none;
        position: fixed;
        top: 0;
        right: 0;
        padding: 0;
        margin: 0;
        z-index: 9999;
    }

    .toasts > .toast {
        position: relative;
        margin: 10px;
        min-width: 40vw;
        position: relative;
        animation: animate-in 350ms forwards;
        color: #fff;
    }

    .toasts > .toast > .content {
        padding: 10px;
        display: block;
        font-weight: 500;
    }

    .toasts > .toast > .progress {
        position: absolute;
        bottom: 0;
        background-color: rgb(0, 0, 0, 0.3);
        height: 6px;
        width: 100%;
        animation-name: shrink;
        animation-timing-function: linear;
        animation-fill-mode: forwards;
    }

    .toasts > .toast:before,
    .toasts > .toast:after {
        content: '';
        position: absolute;
        z-index: -1;
        top: 50%;
        bottom: 0;
        left: 10px;
        right: 10px;
        border-radius: 100px / 10px;
    }

    .toasts > .toast:after {
        right: 10px;
        left: auto;
        transform: skew(8deg) rotate(3deg);
    }

    @keyframes animate-in {
        0% {
            width: 0;
            opacity: 0;
            transform: scale(1.15) translateY(20px);
        }
        100% {
            width: 40vw;
            opacity: 1;
            transform: scale(1) translateY(0);
        }
    }

    @keyframes shrink {
        0% {
            width: 40vw;
        }
        100% {
            width: 0;
        }
    }
</style>

<ul class="toasts">
    {#each toasts as toast (toast.id)}
        <li
            class="toast"
            style="background: {toast.background};"
            out:animateOut
        >
            <div class="content" data-testid="notification-msg">
                {toast.msg}
            </div>
            <div
                class="progress"
                style="animation-duration: {toast.timeout}ms;"
                on:animationend="{() => removeToast(toast.id)}"
            ></div>
        </li>
    {/each}
</ul>
