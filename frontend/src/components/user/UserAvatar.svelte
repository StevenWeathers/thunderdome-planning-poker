<script lang="ts">
  import { AppConfig } from '../../config';
  import LL from '../../i18n/i18n-svelte';

  const { PathPrefix, AvatarService } = AppConfig;
  let klass = '';

  export let warriorId = '';
  export let gravatarHash = '';
  export let avatar = '';
  export let width = 48;
  export { klass as class };
  export let options = {};
</script>

{#if AvatarService === 'dicebear'}
  <img
    src="https://avatars.dicebear.com/api/{avatar}/{warriorId}.svg?w={width}"
    alt="{$LL.avatarAltText()}"
    class="{klass}"
    {...options}
  />
{:else if AvatarService === 'gravatar'}
  {#if gravatarHash !== ''}
    <img
      src="https://gravatar.com/avatar/{gravatarHash}?s={width}&d={avatar}&r=g"
      alt="{$LL.avatarAltText()}"
      class="{klass}"
      {...options}
    />
  {:else}
    <img
      src="https://gravatar.com/avatar/{warriorId}?s={width}&d={avatar}&r=g"
      alt="{$LL.avatarAltText()}"
      class="{klass}"
      {...options}
    />
  {/if}
{:else if AvatarService === 'robohash'}
  <img
    src="https://robohash.org/{warriorId}.png?set={avatar}&size={width}x{width}"
    alt="{$LL.avatarAltText()}"
    class="{klass}"
    {...options}
  />
{:else if AvatarService === 'govatar'}
  <img
    src="{PathPrefix}/avatar/{width}/{warriorId}/{avatar}"
    alt="{$LL.avatarAltText()}"
    class="{klass}"
    {...options}
  />
{:else if AvatarService === 'goadorable'}
  <img
    src="{PathPrefix}/avatar/{width}/{warriorId}"
    alt="{$LL.avatarAltText()}"
    class="{klass}"
    {...options}
  />
{/if}
