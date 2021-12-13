<style>
    .animate-ping {
        -webkit-animation: td-ping 1.5s cubic-bezier(0, 0, 0.2, 1) infinite;
        animation: td-ping 1.5s cubic-bezier(0, 0, 0.2, 1) infinite;
        transform-box: fill-box;
        transform-origin: center center;
    }

    @-webkit-keyframes td-ping {
        75%,
        to {
            opacity: 0;
            transform: scale(2);
        }
    }

    @keyframes td-ping {
        75%,
        to {
            opacity: 0;
            transform: scale(2);
        }
    }
</style>

<div
    class="absolute bottom-0 right-0 w-1/4 h-1/4 rounded bg-red-200 animate-ping"
></div>
<div class="absolute bottom-0 right-0 w-1/4 h-1/4 rounded shadow-md">
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024">
        <rect
            class="text-white fill-current"
            width="1024"
            height="1024"
            rx="192"
            ry="192"></rect>
        <rect
            class="fill-current text-red-500"
            fill="red"
            x="128"
            y="128"
            width="768"
            height="768"
            rx="96"
            ry="96"></rect>
    </svg>
</div>
