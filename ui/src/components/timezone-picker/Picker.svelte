<script lang="ts">
  // Copied from https://github.com/tricinel/svelte-timezone-picker due to lack of svelte4 and typescript support
  // and to modify the appearance to fit dark/light modes
  import { createEventDispatcher, onMount } from 'svelte';
  import { slide } from 'svelte/transition';
  import groupedZones from './timezones';
  import {
    filter,
    keyCodes,
    pick,
    scrollIntoView,
    slugify,
    uid,
    ungroup,
  } from './utils';
  import TextInput from '../forms/TextInput.svelte';

  // ***** Public API *****

  // allow customizing the main btn styles
  export let btnClasses = '';

  // The timezone value comes from the consumer of the component
  // If it's not provided, we will set it in onMount to be the user's current timezone
  export let timezone = null;

  // Should the dropdown be expanded by default?
  export let expanded = false;

  // We can allow the user to filter the timezones displayed to only a few
  export let allowedTimezones = null;

  // ***** End Public API *****

  // What is the current zone?
  // Array ['Abidjan', '+00:00', '+00:00']
  // The first value is the display name for the zone, the second is the standard offset, the third the daylight saving time offset
  let currentZone;

  // We keep track of what the user is typing in the search box
  // String
  let userSearch;

  // What is the currently selected zone in the dropdown?
  // String 'Africa/Abidjan'
  let highlightedZone;

  // DOM nodes refs
  let toggleButtonRef;
  let searchInputRef;
  let clearButtonRef;
  let listBoxRef;
  let listBoxOptionRefs;

  // A few IDs that will we use for a11y
  const labelId = uid();
  const listBoxId = uid();
  const searchInputId = uid();

  // We ungroup the zones
  // e.g. { Africa: {'Africa/Abidjan': ['Abidjan', '+00:00', '+00:00']} }
  // => {'Africa/Abidjan': ['Abidjan', '+00:00', '+00:00']}
  const ungroupedZones = ungroup(groupedZones);

  const userTimezone =
    Intl.DateTimeFormat().resolvedOptions().timeZone || 'America/New_York';

  // We will only display the timezones the user passed in
  // and default to all the zones if that's empty or the wrong format
  let availableZones = ungroupedZones;

  if (allowedTimezones) {
    if (Array.isArray(allowedTimezones)) {
      availableZones = pick(ungroupedZones, [
        ...allowedTimezones,
        userTimezone,
      ]);
    } else {
      console.error(
        'You need to provide a list of timezones as an Array!',
        `You provided ${allowedTimezones}.`,
      );
    }
  }

  // We also want a list of all the valid zones
  // e.g. {'Africa/Abidjan': ['Abidjan', '+00:00', '+00:00'], 'Africa/Accra': ['Accra', '+00:00', '+00:00']}
  // => ['Africa/Abidjan', 'Africa/Accra']
  const validZones = Object.keys(availableZones);

  // Zones will be filtered as the user types, so we keep track of them internally here
  let filteredZones = [];

  // We take the ungroupedZones and create a list of just the user-visible labels
  // and add them to the refs
  // e.g. {'Africa/Abidjan': ['Abidjan', '+00:00', '+00:00'], 'Africa/Accra': ['Accra', '+00:00', '+00:00']}
  // => ['Abidjan', 'Accra']
  listBoxOptionRefs = Object.values(availableZones).map(([zone]) => ({
    [zone]: null,
  }));

  // We keep track of the initial state so we can reset to these values when needed
  const initialState = {
    expanded,
    userSearch: null,
  };

  // Reset the dropdown and all internal state to the initial values
  const reset = () => {
    expanded = initialState.expanded; // eslint-disable-line prefer-destructuring
    userSearch = initialState.userSearch; // eslint-disable-line prefer-destructuring
  };

  // We will use the dispatcher to send the update event
  const dispatch = createEventDispatcher();

  const dispatchUpdates = () => {
    const eventName = 'update';
    const eventData = { timezone };

    dispatch(eventName, eventData);
  };

  // Emit the event back to the consumer
  const handleTimezoneUpdate = (ev, zoneId) => {
    currentZone = ungroupedZones[zoneId];
    timezone = zoneId;
    dispatchUpdates();
    reset();
    toggleButtonRef.focus();
    ev.preventDefault();
  };

  // ***** Methods *****

  // Figure out if a grouped zone has any currently visible zones
  // We use this when the user searches in order to show/hide the group name in the list
  const groupHasVisibleChildren = (group, zones) =>
    Object.keys(groupedZones[group]).some(zone => zones.includes(zone));

  // Scroll the list to a specific element in it if that element is not already visible on screen
  const scrollList = zone => {
    const zoneElementRef = listBoxOptionRefs[zone];
    if (listBoxRef && zoneElementRef) {
      scrollIntoView(zoneElementRef, listBoxRef);
      zoneElementRef.focus({ preventScroll: true });
    }
  };

  // Every time the user uses their keyboard to move up or down in the list,
  // we need to figure out if their at the end/start of the list and scroll the correct elements
  // into view
  const moveSelection = direction => {
    const len = filteredZones.length;
    const zoneIndex = filteredZones.findIndex(zone => zone === highlightedZone);

    let index;

    if (direction === 'up') {
      index = (zoneIndex - 1 + len) % len;
    }

    if (direction === 'down') {
      index = (zoneIndex + 1) % len;
    }

    // We update the highlightedZone to be the one the user is currently on
    highlightedZone = filteredZones[index];
    // We make sure the highlightedZone is visible on screen, scrolling it into view if not
    scrollList(highlightedZone);
  };

  // We watch for when the user presses Escape, ArrowDown or ArrowUp and react accordingly
  const keyDown = ev => {
    // If the clearButton is focused, don't do anything else
    // We should only continue if the dropdown is expanded
    if (document.activeElement === clearButtonRef || !expanded) {
      return;
    }

    // If the user presses Escape, we dismiss the drodpown
    if (ev.keyCode === keyCodes.Escape) {
      reset();
    }

    // If the user presses the down arrow, start navigating the list
    if (ev.keyCode === keyCodes.ArrowDown) {
      ev.preventDefault();
      moveSelection('down');
    }
    // If the user presses the up arrow, start navigating the list
    if (ev.keyCode === keyCodes.ArrowUp) {
      ev.preventDefault();
      moveSelection('up');
    }
    // If the user presses Enter and the dropdown is expanded, select the current item
    if (ev.keyCode === keyCodes.Enter && highlightedZone) {
      handleTimezoneUpdate(ev, highlightedZone);
    }
    // If the user start to type letters or numbers, we focus on the Search field
    if (
      keyCodes.Characters.includes(ev.keyCode) ||
      ev.keyCode === keyCodes.Backspace
    ) {
      searchInputRef.focus();
    }
  };

  // When the user presses the clear button when searching,
  // we want to clear the text and refocus on the input
  const clearSearch = () => {
    userSearch = initialState.userSearch; // eslint-disable-line prefer-destructuring
    // Refocus to the search input
    searchInputRef.focus();
  };

  const setHighlightedZone = zone => {
    highlightedZone = zone;
  };

  const toggleExpanded = ev => {
    if (ev.keyCode) {
      // If it's a keyboard event, we should react only to certain keys
      // Enter and Space should show it
      if ([keyCodes.Enter, keyCodes.Space].includes(ev.keyCode)) {
        expanded = !expanded;
      }
      // Escape should just hide the menu
      if (ev.keyCode === keyCodes.Escape) {
        expanded = false;
      }
      // ArrowDown should show it
      if (ev.keyCode === keyCodes.ArrowDown) {
        expanded = true;
      }
    } else {
      // If there is no keyCode, it's not a keyboard event
      expanded = !expanded;
    }
  };

  const scrollToHighlighted = () => {
    if (expanded && highlightedZone) {
      scrollList(highlightedZone);
    }
  };

  // ***** Reactive *****

  // As the user types, we filter the available zones to show only those that should be visible
  $: filteredZones =
    userSearch && userSearch.length > 0
      ? filter(userSearch, availableZones)
      : validZones.slice();

  const setTimezone = tz => {
    if (!tz) {
      timezone = userTimezone;
    }

    if (tz && !validZones.includes(tz)) {
      // The timezone must be a valid timezone, so we check it against our list of values in flat
      console.warn(
        `The timezone provided is not valid: ${tz}!`,
        `Valid zones are: ${validZones}`,
      );
      timezone = userTimezone;
    }

    currentZone = ungroupedZones[timezone];
    setHighlightedZone(timezone);
  };

  // We want to properly handle any potential changes to the current timezone
  // that might come in from the consumer of the component.
  // This includes setting the proper timezone and dispatching the updated values
  // back up to the consumer
  $: setTimezone(timezone);

  // ***** Lifecycle methods *****
  onMount(() => {
    setTimezone(timezone);
    scrollToHighlighted();
  });
</script>

<style>
  .overlay {
    background: transparent;
    height: 100vh;
    left: 0;
    position: fixed;
    top: 0;
    width: 100vw;
    z-index: 0;
  }

  button {
    background: transparent;
    border: 0;
    cursor: pointer;
  }

  svg polygon {
    fill: currentColor;
  }

  button[data-toggle] {
    align-content: flex-start;
    align-items: center;
    display: flex;
    padding: 0;
  }

  button[data-toggle] > span {
    font-weight: 500;
    margin-right: 0.4em;
    text-decoration: underline;
  }

  button[data-toggle] > span small {
    font-weight: 400;
    font-size: 0.8em;
  }

  .tz-dropdown {
    display: flex;
    flex-direction: column;
    min-width: 18em;
    max-width: 100vw;
    position: absolute;
    z-index: 50;
  }

  .tz-groups {
    height: 240px;
    max-height: 40vh;
    overflow: scroll;
  }

  ul {
    margin: 0;
    list-style: none inside none;
    padding: 0;
  }

  ul li {
    font-size: 0.9rem;
    display: block;
    margin: 0;
    padding: 0;
  }

  ul li > span {
    font-size: 0.8em;
    line-height: 1.4em;
    text-align: right;
  }

  ul li p {
    font-size: 0.92rem;
    font-weight: 600;
    letter-spacing: 0.08em;
    margin: 0;
    padding: 0;
    text-transform: uppercase;
  }

  ul li {
    border: 0;
    display: flex;
    justify-content: space-between;
    padding: 0.8em 1.2em;
    text-align: left;
  }

  ul li[aria-selected]:hover,
  ul li:focus,
  li[aria-selected='true'] {
    cursor: pointer;
  }

  .sr-only {
    border: 0;
    clip: 'rect(0, 0, 0, 0)';
    height: 1px;
    margin: -1px;
    opacity: 0;
    overflow: hidden;
    padding: 0;
    position: absolute;
    width: 1px;
  }
</style>

{#if expanded}
  <div class="overlay" on:click="{reset}"></div>
{/if}

<button
  bind:this="{toggleButtonRef}"
  type="button"
  aria-label="{`${currentZone[0]} is currently selected. Change timezone`}"
  aria-haspopup="listbox"
  data-toggle="true"
  aria-expanded="{expanded}"
  on:click="{toggleExpanded}"
  on:keydown="{toggleExpanded}"
  class="{btnClasses}"
>
  <span>{currentZone[0]} <small>GMT {currentZone[1]}</small></span>
  <svg width="10" height="16" viewBox="0 0 16 16" fill="currentColor">
    <polygon
      x="0"
      y="0"
      points="8, 8, 16, 16, 0, 16"
      transform="{expanded ? 'rotate(0)' : 'rotate(180, 8, 8)'} translate(0 -4)"
    ></polygon>
  </svg>
</button>
{#if expanded}
  <div
    class="tz-dropdown rounded shadow bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-400"
    transition:slide
    on:introend="{scrollToHighlighted}"
    on:keydown="{keyDown}"
  >
    <span class="sr-only" id="{labelId}">
      Select a timezone from the list. Start typing to filter or use the arrow
      keys to navigate the list
    </span>
    <!-- svelte-ignore a11y-autofocus -->
    <div class="w-full p-2">
      <TextInput
        id="{searchInputId}"
        bind:this="{searchInputRef}"
        type="search"
        aria-autocomplete="list"
        aria-controls="{listBoxId}"
        aria-labelledby="{labelId}"
        autocomplete="off"
        autocorrect="off"
        placeholder="Search..."
        bind:value="{userSearch}"
        autofocus
        class="searchField"
      />
    </div>

    <ul
      tabindex="-1"
      class="tz-groups"
      id="{listBoxId}"
      role="listbox"
      bind:this="{listBoxRef}"
      aria-labelledby="{labelId}"
      aria-activedescendant="{currentZone && `tz-${slugify(currentZone[0])}`}"
    >
      {#each Object.keys(groupedZones) as group}
        {#if groupHasVisibleChildren(group, filteredZones)}
          <li role="option" aria-hidden="true">
            <p>{group}</p>
          </li>
          {#each Object.entries(groupedZones[group]) as [zoneLabel, zoneDetails]}
            {#if filteredZones.includes(zoneLabel)}
              <li
                role="option"
                tabindex="0"
                id="{`tz-${slugify(zoneLabel)}`}"
                bind:this="{listBoxOptionRefs[zoneLabel]}"
                aria-label="{`Select ${zoneDetails[0]}`}"
                aria-selected="{highlightedZone === zoneDetails[0]}"
                on:mouseover="{() => setHighlightedZone(zoneDetails[0])}"
                on:click="{ev => handleTimezoneUpdate(ev, zoneLabel)}"
                class="hover:bg-blue-500 hover:text-white dark:hover:bg-sky-300 dark:hover:text-gray-800 focus:bg-blue-500 focus:text-white dark:focus:bg-sky-300 dark:focus:text-gray-800"
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
