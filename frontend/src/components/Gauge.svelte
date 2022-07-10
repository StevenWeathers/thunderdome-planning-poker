<script>
    import { afterUpdate } from 'svelte'
    import Snap from 'snapsvg-cjs'

    export let percentage = 100
    export let text = ''
    export let color = 'blue'
    export let stat = percentage
    export let count = ''

    let svgElem

    let polar_to_cartesian, svg_circle_arc_path, animate_arc

    polar_to_cartesian = function (cx, cy, radius, angle) {
        let radians
        radians = ((angle - 90) * Math.PI) / 180.0
        return [
            Math.round((cx + radius * Math.cos(radians)) * 100) / 100,
            Math.round((cy + radius * Math.sin(radians)) * 100) / 100,
        ]
    }

    svg_circle_arc_path = function (x, y, radius, start_angle, end_angle) {
        let end_xy, start_xy
        start_xy = polar_to_cartesian(x, y, radius, end_angle)
        end_xy = polar_to_cartesian(x, y, radius, start_angle)
        return (
            'M ' +
            start_xy[0] +
            ' ' +
            start_xy[1] +
            ' A ' +
            radius +
            ' ' +
            radius +
            ' 0 0 0 ' +
            end_xy[0] +
            ' ' +
            end_xy[1]
        )
    }

    animate_arc = function (ratio, svg) {
        const arc = svg.select('.gaugeNeedle')
        const currentRatio = parseFloat(arc.attr('data-ratio'))
        arc.attr('data-ratio', ratio) // update the ratio

        return Snap.animate(
            currentRatio,
            ratio,
            function (val) {
                const path = svg_circle_arc_path(
                    500,
                    500,
                    450,
                    -90,
                    val * 180.0 - 90,
                )
                arc.attr('d', path)
            },
            1500,
            mina.easeinout,
        )
    }

    afterUpdate(() => {
        const perc = percentage <= 100 ? percentage : 100 // account for overage
        animate_arc(perc / 100, Snap(svgElem))
    })
</script>

<style global>
    .gauge {
        --tw-aspect-h: 1;
        --tw-aspect-w: 2;
        padding-bottom: calc(var(--tw-aspect-h) / var(--tw-aspect-w) * 100%);
        position: relative;
    }

    .gauge > * {
        bottom: 0;
        height: 100%;
        left: 0;
        position: absolute;
        right: 0;
        top: 0;
        width: 100%;
    }

    .gauge path {
        stroke-width: 75;
        fill: none;
    }

    .gauge.blue path {
        @apply stroke-sky-500;
    }

    .gauge.blue .percentage,
    .gauge.blue .count-text {
        @apply text-sky-500;
    }

    .gauge.green path {
        @apply stroke-green-500;
    }

    .gauge.green .percentage,
    .gauge.green .count-text {
        @apply text-green-500;
    }

    .gauge.red path {
        @apply stroke-red-500;
    }

    .gauge.red .percentage,
    .gauge.red .count-text {
        @apply text-red-500;
    }

    .gauge.purple path {
        @apply stroke-indigo-500;
    }

    .gauge.purple .percentage,
    .gauge.purple .count-text {
        @apply text-indigo-500;
    }
</style>

<div class="relative gauge {color}">
    {#if count !== ''}
        <div class="absolute text-right text-sm md:text-lg count-text">
            {count}
        </div>
    {/if}
    <svg viewBox="0 0 1000 500" bind:this="{svgElem}" class="max-w-full">
        <g class="opacity-10">
            <path d="M 950 500 A 450 450 0 0 0 50 500"></path>
        </g>
        <path
            d="M 50 500 A 450 450 0 0 0 50 500"
            class="gaugeNeedle"
            data-ratio="0"></path>
    </svg>
    <div class="absolute w-full h-full flex flex-col items-center justify-end">
        <h3
            class="block text-3xl md:text-5xl font-black ml-1.5 tracking-tight percentage"
        >
            <span>{stat}</span><span class="ml-0.5 text-base font-bold">%</span>
        </h3>
        <h4
            class="mt-0.5 font-bold text-gray-500 dark:text-gray-400 uppercase text-sm md:text-md tracking-wide"
        >
            {text}
        </h4>
    </div>
</div>
