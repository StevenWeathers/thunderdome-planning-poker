<script>
    import { afterUpdate } from 'svelte'
    import Snap from 'snapsvg-cjs'

    export let percentage = 100
    export let text = ''
    export let color = 'blue'
    export let stat = `${percentage}%`

    let svgElem
    let percElem

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

    animate_arc = function (ratio, svg, perc) {
        let arc = svg.path('')
        return Snap.animate(
            0,
            ratio,
            function (val) {
                let path
                arc.remove()
                path = svg_circle_arc_path(500, 500, 450, -90, val * 180.0 - 90)
                arc = svg.path(path)
                arc.attr({
                    class: 'data-arc',
                })
                perc.innerText = Math.round(val * 100) + '%'
            },
            Math.round(2000 * ratio),
            mina.easeinout,
        )
    }

    afterUpdate(() => {
        const perc = percentage <= 100 ? percentage : 100 // account for overage
        animate_arc(perc / 100, window.Snap(svgElem), percElem)
    })
</script>

<style global>
    .gauge path {
        stroke-width: 75;
        @apply stroke-slate-200;
        fill: none;
    }

    .gauge .title {
        @apply fill-gray-500;
    }

    .gauge.blue path.data-arc {
        @apply stroke-cyan-400;
    }

    .gauge.blue .percentage {
        @apply fill-cyan-600;
    }

    .gauge.green path.data-arc {
        @apply stroke-lime-400;
    }

    .gauge.green .percentage {
        @apply fill-lime-600;
    }

    .gauge.red path.data-arc {
        @apply stroke-red-500;
    }

    .gauge.red .percentage {
        @apply fill-red-700;
    }
</style>

<div class="gauge {color}">
    <svg viewBox="0 0 1000 500" bind:this="{svgElem}" class="max-w-full">
        <path d="M 950 500 A 450 450 0 0 0 50 500"></path>
        <text
            class="percentage"
            text-anchor="middle"
            alignment-baseline="middle"
            x="500"
            y="300"
            font-size="140"
            font-weight="bold"
            bind:this="{percElem}"
            >{stat}
        </text>
        <text
            class="title"
            text-anchor="middle"
            alignment-baseline="middle"
            x="500"
            y="450"
            font-size="90"
            font-weight="normal"
            >{text}
        </text>
    </svg>
</div>
