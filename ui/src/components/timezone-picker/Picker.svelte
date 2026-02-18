<script lang="ts">
  // Copied from https://github.com/tricinel/svelte-timezone-picker due to lack of svelte4 and typescript support
  // and to modify the appearance to fit dark/light modes
  import { slide } from 'svelte/transition';
  import groupedZones from './timezones';
  import { filter, keyCodes, pick, scrollIntoView, slugify, uid, ungroup } from './utils';
  import TextInput from '../forms/TextInput.svelte';

  // ***** Types *****
  type TimezoneDetails = [string, string, string]; // [displayName, gmtOffset, dstOffset]
  type UngroupedZones = Record<string, TimezoneDetails>;
  type GroupedZones = Record<string, UngroupedZones>;
  type ListBoxRefs = Record<string, HTMLLIElement | null>;

  // ***** Public API *****
  interface Props {
    timezone?: string | null;
    expanded?: boolean;
    allowedTimezones?: string[] | null;
    onUpdate?: (timezone: string | null) => void;
  }

  let { timezone = $bindable(null), expanded = false, allowedTimezones = null, onUpdate }: Props = $props();

  // ***** End Public API *****

  // State variables
  let currentZone = $state<TimezoneDetails | undefined>();
  let userSearch = $state<string | null>(null);
  let highlightedZone = $state<string | undefined>();

  // DOM refs
  let toggleButtonRef = $state<HTMLButtonElement | undefined>();
  let searchInputRef = $state<any>();
  let clearButtonRef = $state<HTMLButtonElement | undefined>();
  let listBoxRef = $state<HTMLUListElement | undefined>();

  // Constants
  const labelId = uid();
  const listBoxId = uid();
  const searchInputId = uid();

  // We ungroup the zones
  const ungroupedZones = ungroup(groupedZones) as UngroupedZones;
  const userTimezone = Intl.DateTimeFormat().resolvedOptions().timeZone || 'America/New_York';

  // Process allowed timezones
  const availableZones = $derived.by<UngroupedZones>(() => {
    let zones = ungroupedZones;
    if (allowedTimezones) {
      if (Array.isArray(allowedTimezones)) {
        zones = pick(ungroupedZones, [...allowedTimezones, userTimezone]) as UngroupedZones;
      } else {
        console.error('You need to provide a list of timezones as an Array!', `You provided ${allowedTimezones}.`);
      }
    }
    return zones;
  });

  const validZones = $derived(Object.keys(availableZones));

  // Initialize option refs
  let listBoxOptionRefs = $state<ListBoxRefs>({});

  // Initial state
  let internalExpanded = $state(false);

  $effect(() => {
    internalExpanded = expanded;
  });

  const initialState = {
    expanded: false,
    userSearch: null,
  };

  // Derived state
  let filteredZones = $derived<string[]>(
    userSearch && userSearch.length > 0 ? (filter(userSearch, availableZones) as string[]) : validZones.slice(),
  );

  // ***** Methods *****

  // Reset the dropdown and all internal state to the initial values
  const reset = () => {
    expanded = initialState.expanded;
    userSearch = initialState.userSearch;
  };

  // Custom event dispatch and callback
  const dispatchUpdates = () => {
    // Call the onUpdate callback if provided
    if (onUpdate) {
      onUpdate(timezone);
    }

    // Also dispatch a custom event for compatibility
    const event = new CustomEvent('timezoneUpdate', {
      detail: { timezone },
      bubbles: true,
    });

    if (toggleButtonRef) {
      toggleButtonRef.dispatchEvent(event);
    }
  };

  // Emit the event back to the consumer
  const handleTimezoneUpdate = (ev: Event, zoneId: string) => {
    currentZone = ungroupedZones[zoneId];
    timezone = zoneId;
    dispatchUpdates();
    reset();
    toggleButtonRef?.focus();
    ev.preventDefault();
  };

  // Figure out if a grouped zone has any currently visible zones
  const groupHasVisibleChildren = (group: string, zones: string[]) =>
    Object.keys((groupedZones as unknown as GroupedZones)[group]).some(zone => zones.includes(zone));

  // Scroll the list to a specific element
  const scrollList = (zone: string) => {
    const zoneElementRef = listBoxOptionRefs[zone];
    if (listBoxRef && zoneElementRef) {
      scrollIntoView(zoneElementRef, listBoxRef);
      zoneElementRef.focus({ preventScroll: true });
    }
  };

  // Move selection up or down
  const moveSelection = (direction: 'up' | 'down') => {
    const len = filteredZones.length;
    const zoneIndex = filteredZones.findIndex(zone => zone === highlightedZone);

    let index: number;
    if (direction === 'up') {
      index = (zoneIndex - 1 + len) % len;
    } else {
      index = (zoneIndex + 1) % len;
    }

    highlightedZone = filteredZones[index];
    scrollList(highlightedZone);
  };

  // Handle keyboard navigation
  const keyDown = (ev: KeyboardEvent) => {
    if (document.activeElement === clearButtonRef || !expanded) {
      return;
    }

    if (ev.keyCode === keyCodes.Escape) {
      reset();
    }
    if (ev.keyCode === keyCodes.ArrowDown) {
      ev.preventDefault();
      moveSelection('down');
    }
    if (ev.keyCode === keyCodes.ArrowUp) {
      ev.preventDefault();
      moveSelection('up');
    }
    if (ev.keyCode === keyCodes.Enter && highlightedZone) {
      handleTimezoneUpdate(ev, highlightedZone);
    }
    if (keyCodes.Characters.includes(ev.keyCode) || ev.keyCode === keyCodes.Backspace) {
      searchInputRef?.focus();
    }
  };

  // Clear search
  const clearSearch = () => {
    userSearch = initialState.userSearch;
    searchInputRef?.focus();
  };

  const setHighlightedZone = (zone: string) => {
    highlightedZone = zone;
  };

  const toggleExpanded = (ev: MouseEvent | KeyboardEvent) => {
    if ('keyCode' in ev && ev.keyCode) {
      if ([keyCodes.Enter, keyCodes.Space].includes(ev.keyCode)) {
        expanded = !expanded;
      }
      if (ev.keyCode === keyCodes.Escape) {
        expanded = false;
      }
      if (ev.keyCode === keyCodes.ArrowDown) {
        expanded = true;
      }
    } else {
      expanded = !expanded;
    }
  };

  const scrollToHighlighted = () => {
    if (expanded && highlightedZone) {
      scrollList(highlightedZone);
    }
  };

  const setTimezone = (tz: string | null) => {
    if (!tz) {
      timezone = userTimezone;
    }

    if (tz && !validZones.includes(tz)) {
      console.warn(`The timezone provided is not valid: ${tz}!`, `Valid zones are: ${validZones}`);
      timezone = userTimezone;
    }

    if (timezone) {
      currentZone = ungroupedZones[timezone];
      setHighlightedZone(timezone);
    }
  };

  // Effects - replace reactive statements and onMount
  $effect(() => {
    setTimezone(timezone);
  });

  $effect(() => {
    if (expanded && highlightedZone) {
      scrollToHighlighted();
    }
  });
</script>

{#if expanded}
  <div
    class="overlay"
    onclick={reset}
    onkeydown={e => {
      if (e.key === 'Enter' || e.key === ' ') reset();
    }}
    aria-label="Close timezone picker"
    aria-modal="true"
    tabindex="-1"
    role="dialog"
  ></div>
{/if}

<div class="timezone-picker-container">
  <button
    bind:this={toggleButtonRef}
    type="button"
    aria-label={`${currentZone?.[0]} is currently selected. Change timezone`}
    aria-haspopup="listbox"
    data-toggle="true"
    aria-expanded={expanded}
    onclick={toggleExpanded}
    onkeydown={toggleExpanded}
  >
    <span>
      {currentZone?.[0]}
      <small>GMT {currentZone?.[1]}</small>
    </span>
    <svg width="10" height="10" viewBox="0 0 12 12" fill="none" class="chevron">
      <path
        d="M2.5 4.5L6 8L9.5 4.5"
        stroke="currentColor"
        stroke-width="1.5"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
    </svg>
  </button>

  {#if expanded}
    <div
      class="tz-dropdown"
      transition:slide={{ duration: 200, easing: t => t * (2 - t) }}
      onintroend={scrollToHighlighted}
      onkeydown={keyDown}
      role="dialog"
      tabindex="0"
    >
      <span class="sr-only" id={labelId}>
        Select a timezone from the list. Start typing to filter or use the arrow keys to navigate the list
      </span>

      <div class="search-container">
        <TextInput
          id={searchInputId}
          bind:this={searchInputRef}
          type="search"
          aria-autocomplete="list"
          aria-controls={listBoxId}
          aria-labelledby={labelId}
          autocomplete="off"
          autocorrect="off"
          placeholder="Search timezones..."
          bind:value={userSearch}
          autofocus
          class="searchField"
        />
      </div>

      <ul
        tabindex="-1"
        class="tz-groups"
        id={listBoxId}
        role="listbox"
        bind:this={listBoxRef}
        aria-labelledby={labelId}
        aria-activedescendant={currentZone && `tz-${slugify(currentZone[0])}`}
      >
        {#each Object.keys(groupedZones as unknown as GroupedZones) as group}
          {#if groupHasVisibleChildren(group, filteredZones)}
            <li role="option" aria-hidden="true" aria-selected="false">
              <p>{group}</p>
            </li>
            {#each Object.entries((groupedZones as unknown as GroupedZones)[group]) as [zoneLabel, zoneDetails]: [string, TimezoneDetails]}
              {#if filteredZones.includes(zoneLabel)}
                <li
                  role="option"
                  tabindex="0"
                  id={`tz-${slugify(zoneLabel)}`}
                  bind:this={listBoxOptionRefs[zoneLabel]}
                  aria-label={`Select ${zoneDetails[0]}`}
                  aria-selected={highlightedZone === zoneDetails[0]}
                  onmouseover={() => setHighlightedZone(zoneDetails[0])}
                  onfocus={() => setHighlightedZone(zoneDetails[0])}
                  onclick={ev => handleTimezoneUpdate(ev, zoneLabel)}
                  onkeydown={ev => {
                    if (ev.key === 'Enter' || ev.key === ' ') {
                      handleTimezoneUpdate(ev, zoneLabel);
                    }
                  }}
                >
                  {zoneDetails[0]} <span>GMT {zoneDetails[1]}</span>
                </li>
              {/if}
            {/each}
          {/if}
        {/each}
      </ul>
    </div>
  {/if}
</div>

<style>
  .overlay {
    background: rgba(0, 0, 0, 0.1);
    backdrop-filter: blur(2px);
    height: 100vh;
    left: 0;
    position: fixed;
    top: 0;
    width: 100vw;
    z-index: 40;
    transition: all 0.2s ease-in-out;
  }

  button {
    background: transparent;
    border: 0;
    cursor: pointer;
    transition: all 0.2s ease-in-out;
  }

  button[data-toggle] {
    align-items: center;
    background: transparent;
    border: none;
    border-radius: 6px;
    display: inline-flex;
    font-family:
      'Inter',
      -apple-system,
      BlinkMacSystemFont,
      'Segoe UI',
      Roboto,
      sans-serif;
    font-size: 0.9rem;
    padding: 4px 8px;
    position: relative;
    transition: all 0.15s ease-in-out;
    color: #374151;
    text-decoration: underline;
    text-decoration-color: #d1d5db;
    text-underline-offset: 3px;
    text-decoration-thickness: 1px;
  }

  button[data-toggle]:hover {
    color: #1f2937;
    text-decoration-color: #6b7280;
    background: rgba(59, 130, 246, 0.05);
  }

  button[data-toggle]:focus {
    outline: none;
    color: #3b82f6;
    text-decoration-color: #3b82f6;
    box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
  }

  button[data-toggle]:active {
    color: #1e40af;
    text-decoration-color: #1e40af;
  }

  /* Dark mode styles for button */
  :global(.dark) button[data-toggle] {
    color: #d1d5db;
    text-decoration-color: #4b5563;
  }

  :global(.dark) button[data-toggle]:hover {
    color: #f9fafb;
    text-decoration-color: #9ca3af;
    background: rgba(96, 165, 250, 0.1);
  }

  :global(.dark) button[data-toggle]:focus {
    color: #60a5fa;
    text-decoration-color: #60a5fa;
    box-shadow: 0 0 0 2px rgba(96, 165, 250, 0.3);
  }

  button[data-toggle] > span {
    font-weight: 500;
    text-align: left;
    line-height: 1.3;
  }

  button[data-toggle] > span small {
    font-weight: 400;
    font-size: 0.85em;
    color: #6b7280;
    margin-left: 6px;
    opacity: 0.8;
  }

  :global(.dark) button[data-toggle] > span small {
    color: #9ca3af;
  }

  .chevron {
    margin-left: 6px;
    transition: transform 0.15s ease-in-out;
    color: #9ca3af;
    opacity: 0.7;
  }

  :global(.dark) .chevron {
    color: #6b7280;
  }

  button[data-toggle][aria-expanded='true'] .chevron {
    transform: rotate(180deg);
    opacity: 1;
  }

  .timezone-picker-container {
    position: relative;
    display: inline-block;
  }

  .tz-dropdown {
    display: flex;
    flex-direction: column;
    min-width: 320px;
    max-width: 400px;
    position: absolute;
    z-index: 50;
    top: 100%;
    left: 0;
    margin-top: 4px;
    background: white;
    border: 1px solid #e2e8f0;
    border-radius: 16px;
    box-shadow:
      0 20px 25px -5px rgba(0, 0, 0, 0.1),
      0 10px 10px -5px rgba(0, 0, 0, 0.04);
    overflow: hidden;
    backdrop-filter: blur(8px);
  }

  :global(.dark) .tz-dropdown {
    background: #1f2937;
    border-color: #374151;
    box-shadow:
      0 20px 25px -5px rgba(0, 0, 0, 0.4),
      0 10px 10px -5px rgba(0, 0, 0, 0.2);
  }

  .search-container {
    padding: 20px 20px 16px;
    border-bottom: 1px solid #f1f5f9;
    background: #fafbfc;
  }

  :global(.dark) .search-container {
    background: #111827;
    border-bottom-color: #374151;
  }

  .tz-groups {
    height: 280px;
    max-height: 50vh;
    overflow-y: auto;
    scrollbar-width: thin;
    scrollbar-color: #cbd5e1 transparent;
  }

  .tz-groups::-webkit-scrollbar {
    width: 6px;
  }

  .tz-groups::-webkit-scrollbar-track {
    background: transparent;
  }

  .tz-groups::-webkit-scrollbar-thumb {
    background-color: #cbd5e1;
    border-radius: 3px;
  }

  :global(.dark) .tz-groups::-webkit-scrollbar-thumb {
    background-color: #4b5563;
  }

  ul {
    margin: 0;
    list-style: none;
    padding: 8px 0;
  }

  ul li {
    font-size: 0.95rem;
    display: block;
    margin: 0;
    padding: 0;
    font-family:
      'Inter',
      -apple-system,
      BlinkMacSystemFont,
      'Segoe UI',
      Roboto,
      sans-serif;
  }

  ul li > span {
    font-size: 0.8em;
    line-height: 1.4em;
    text-align: right;
    color: #64748b;
    font-weight: 500;
    font-family: 'JetBrains Mono', 'SF Mono', Monaco, Consolas, monospace;
  }

  :global(.dark) ul li > span {
    color: #94a3b8;
  }

  ul li p {
    font-size: 0.8rem;
    font-weight: 700;
    letter-spacing: 0.1em;
    margin: 0;
    padding: 12px 20px 8px;
    text-transform: uppercase;
    color: #475569;
    background: #f8fafc;
    border-bottom: 1px solid #f1f5f9;
    position: sticky;
    top: 0;
    z-index: 10;
  }

  :global(.dark) ul li p {
    color: #94a3b8;
    background: #0f172a;
    border-bottom-color: #1e293b;
  }

  ul li[role='option']:not([aria-hidden='true']) {
    border: 0;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 20px;
    text-align: left;
    transition: all 0.15s ease-in-out;
    cursor: pointer;
    position: relative;
    color: #1f2937;
  }

  :global(.dark) ul li[role='option']:not([aria-hidden='true']) {
    color: #f9fafb;
  }

  ul li[role='option']:not([aria-hidden='true']):hover {
    background: linear-gradient(135deg, #eff6ff 0%, #dbeafe 100%);
    color: #1e40af;
  }

  ul li[role='option']:not([aria-hidden='true']):focus,
  ul li[aria-selected='true']:not([aria-hidden='true']) {
    background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
    color: white;
    outline: none;
    position: relative;
  }

  ul li[role='option']:not([aria-hidden='true']):focus > span,
  ul li[aria-selected='true']:not([aria-hidden='true']) > span {
    color: rgba(255, 255, 255, 0.8);
  }

  :global(.dark) ul li[role='option']:not([aria-hidden='true']):hover {
    background: linear-gradient(135deg, #1e3a8a 0%, #1e40af 100%);
    color: #dbeafe;
  }

  :global(.dark) ul li[role='option']:not([aria-hidden='true']):focus,
  :global(.dark) ul li[aria-selected='true']:not([aria-hidden='true']) {
    background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
    color: white;
  }

  .sr-only {
    border: 0;
    clip: rect(0, 0, 0, 0);
    height: 1px;
    margin: -1px;
    opacity: 0;
    overflow: hidden;
    padding: 0;
    position: absolute;
    width: 1px;
  }

  /* Enhanced search input styling */
  :global(.searchField) {
    width: 100% !important;
    padding: 12px 16px !important;
    border: 1px solid #e2e8f0 !important;
    border-radius: 10px !important;
    font-size: 0.95rem !important;
    font-family:
      'Inter',
      -apple-system,
      BlinkMacSystemFont,
      'Segoe UI',
      Roboto,
      sans-serif !important;
    transition: all 0.2s ease-in-out !important;
    background: white !important;
  }

  :global(.searchField:focus) {
    outline: none !important;
    border-color: #3b82f6 !important;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1) !important;
  }

  :global(.dark .searchField) {
    background: #374151 !important;
    border-color: #4b5563 !important;
    color: #f9fafb !important;
  }

  :global(.dark .searchField:focus) {
    border-color: #60a5fa !important;
    box-shadow: 0 0 0 3px rgba(96, 165, 250, 0.2) !important;
  }

  :global(.searchField::placeholder) {
    color: #9ca3af !important;
    font-weight: 400 !important;
  }

  :global(.dark .searchField::placeholder) {
    color: #6b7280 !important;
  }
</style>
