<script lang="ts">
  import Snap from 'snapsvg-cjs';
  import { Check, Minus } from '@lucide/svelte';
  import UserAvatar from './user/UserAvatar.svelte';

  interface Props {
    percentage?: number;
    text?: string;
    color?: string;
    stat?: number;
    count?: string;
    details?: Array<{
      id: string;
      name: string;
      avatar?: string;
      gravatarHash?: string;
      pictureUrl?: string;
      met: boolean;
    }>;
  }

  let { percentage = 100, text = '', color = 'blue', stat: providedStat, count = '', details = [] }: Props = $props();

  let detailsList = $derived(details);
  let displayStat = $derived(providedStat ?? percentage);

  type GaugeDetail = {
    id: string;
    name: string;
    avatar?: string;
    gravatarHash?: string;
    pictureUrl?: string;
    met: boolean;
  };

  let svgElem: SVGSVGElement;
  let showDetails = $state(false);

  let polar_to_cartesian: (cx: number, cy: number, radius: number, angle: number) => number[],
    svg_circle_arc_path: (x: number, y: number, radius: number, start_angle: number, end_angle: number) => string,
    animate_arc: (ratio: number, svg: any) => any;

  polar_to_cartesian = function (cx: number, cy: number, radius: number, angle: number) {
    let radians;
    radians = ((angle - 90) * Math.PI) / 180.0;
    return [
      Math.round((cx + radius * Math.cos(radians)) * 100) / 100,
      Math.round((cy + radius * Math.sin(radians)) * 100) / 100,
    ];
  };

  svg_circle_arc_path = function (x: number, y: number, radius: number, start_angle: number, end_angle: number) {
    let end_xy, start_xy;
    start_xy = polar_to_cartesian(x, y, radius, end_angle);
    end_xy = polar_to_cartesian(x, y, radius, start_angle);
    return (
      'M ' + start_xy[0] + ' ' + start_xy[1] + ' A ' + radius + ' ' + radius + ' 0 0 0 ' + end_xy[0] + ' ' + end_xy[1]
    );
  };

  animate_arc = function (ratio: number, svg: any) {
    const arc = svg.select('.gaugeNeedle');
    const currentRatio = parseFloat(arc.attr('data-ratio'));
    arc.attr('data-ratio', ratio); // update the ratio

    return Snap.animate(
      currentRatio,
      ratio,
      function (val: number) {
        const path = svg_circle_arc_path(500, 500, 450, -90, val * 180.0 - 90);
        arc.attr('d', path);
      },
      1500,
      (globalThis as any).mina.easeinout,
    );
  };

  $effect(() => {
    const perc = percentage <= 100 ? percentage : 100; // account for overage
    animate_arc(perc / 100, Snap(svgElem));
  });
</script>

<div class="relative gauge-wrap">
  <div
    class="relative gauge {color} transition-transform duration-200 hover:-translate-y-1 focus-within:-translate-y-1"
  >
    <div class="gauge-surface">
      {#if count !== ''}
        <div class="absolute text-right text-sm md:text-lg count-text">
          {count}
        </div>
      {/if}
      <svg viewBox="0 0 1000 500" bind:this={svgElem} class="max-w-full">
        <g class="opacity-10">
          <path d="M 950 500 A 450 450 0 0 0 50 500"></path>
        </g>
        <path d="M 50 500 A 450 450 0 0 0 50 500" class="gaugeNeedle" data-ratio="0"></path>
      </svg>
      <div class="absolute w-full h-full flex flex-col items-center justify-end">
        <h3 class="block text-3xl md:text-5xl font-black ms-1.5 tracking-tight percentage">
          <span>{displayStat}</span><span class="ms-0.5 text-base font-bold">%</span>
        </h3>
        <h4 class="mt-0.5 font-bold text-gray-500 dark:text-gray-400 uppercase text-sm md:text-md tracking-wide">
          {text}
        </h4>
      </div>

      <button
        type="button"
        class="absolute inset-0 z-10 cursor-pointer appearance-none border-0 bg-transparent p-0"
        onmouseenter={() => {
          showDetails = true;
        }}
        onmouseleave={() => {
          showDetails = false;
        }}
        onfocus={() => {
          showDetails = true;
        }}
        onblur={() => {
          showDetails = false;
        }}
        aria-label={text}
      ></button>
    </div>
  </div>

  {#if detailsList.length > 0}
    <div
      class="absolute inset-x-0 top-full z-30 mt-3 flex justify-center px-2 transition duration-150 {showDetails
        ? 'pointer-events-none opacity-100 translate-y-0'
        : 'pointer-events-none opacity-0 -translate-y-1'}"
      aria-hidden={!showDetails}
    >
      <div
        class="w-full max-w-xs rounded-2xl border border-slate-300 bg-white p-3 shadow-2xl {showDetails
          ? 'pointer-events-auto'
          : 'pointer-events-none'} dark:border-gray-600 dark:bg-gray-900"
      >
        <div class="mb-2 flex items-center justify-between gap-3">
          <h5 class="truncate text-sm font-semibold text-slate-900 dark:text-slate-100">{text}</h5>
          <span
            class="rounded-full bg-slate-100 px-2 py-1 text-[11px] font-semibold uppercase tracking-wide text-slate-700 dark:bg-gray-700 dark:text-slate-200"
          >
            {detailsList.filter((item: GaugeDetail) => item.met).length}/{detailsList.length}
          </span>
        </div>

        <div class="max-h-60 space-y-2 overflow-y-auto pe-1">
          {#each detailsList as item}
            <div
              class="flex items-center gap-2 rounded-xl px-2 py-1.5 {item.met
                ? 'bg-emerald-50/80 dark:bg-emerald-500/10'
                : 'bg-slate-50 dark:bg-gray-700/60'}"
            >
              <UserAvatar
                warriorId={item.id}
                userName={item.name}
                avatar={item.avatar || ''}
                gravatarHash={item.gravatarHash || ''}
                pictureUrl={item.pictureUrl || ''}
                width={28}
                class="h-7 w-7 rounded-full"
              />
              <span
                class="min-w-0 flex-1 truncate text-sm font-medium {item.met
                  ? 'text-slate-900 dark:text-slate-100'
                  : 'text-slate-500 dark:text-slate-300'}"
              >
                {item.name}
              </span>
              <span
                class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full {item.met
                  ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-500/20 dark:text-emerald-300'
                  : 'bg-slate-200 text-slate-500 dark:bg-gray-600 dark:text-slate-300'}"
                aria-hidden="true"
              >
                {#if item.met}
                  <Check class="h-3.5 w-3.5" />
                {:else}
                  <Minus class="h-3.5 w-3.5" />
                {/if}
              </span>
            </div>
          {/each}
        </div>
      </div>
    </div>
  {/if}
</div>

<style global lang="postcss">
  .gauge-wrap {
    position: relative;
  }

  .gauge {
    --tw-aspect-h: 1;
    --tw-aspect-w: 2;
    padding-bottom: calc(var(--tw-aspect-h) / var(--tw-aspect-w) * 100%);
    position: relative;
  }

  .gauge-surface > * {
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
