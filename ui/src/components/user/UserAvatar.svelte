<script lang="ts">
  import { AppConfig } from '../../config';
  import LL from '../../i18n/i18n-svelte';

  const { PathPrefix, AvatarService } = AppConfig;

  interface Props {
    class?: string;
    pictureUrl?: string;
    warriorId?: string;
    gravatarHash?: string;
    avatar?: string;
    userName?: string;
    width?: number;
    options?: any;
  }

  let {
    class: klass = '',
    pictureUrl = '',
    warriorId = '',
    gravatarHash = '',
    avatar = '',
    userName = '',
    width = 48,
    options = {},
  }: Props = $props();
</script>

{#if pictureUrl !== ''}
  <img src={pictureUrl} alt={$LL.avatarAltText()} class={klass} title={userName} {...options} />
{:else if AvatarService === 'gravatar'}
  {#if gravatarHash !== ''}
    <img
      src="https://gravatar.com/avatar/{gravatarHash}?s={width}&d={avatar !== '' ? avatar : 'mp'}&r=g"
      alt={$LL.avatarAltText()}
      class={klass}
      title={userName}
      {...options}
    />
  {:else}
    <img
      src="https://gravatar.com/avatar/{warriorId}?s={width}&d={avatar !== '' ? avatar : 'mp'}&r=g"
      alt={$LL.avatarAltText()}
      class={klass}
      title={userName}
      {...options}
    />
  {/if}
{:else if AvatarService === 'robohash'}
  <img
    src="https://robohash.org/{warriorId}.png?set={avatar}&size={width}x{width}"
    alt={$LL.avatarAltText()}
    class={klass}
    title={userName}
    {...options}
  />
{:else if AvatarService === 'govatar'}
  <img
    src="{PathPrefix}/avatar/{width}/{warriorId}/{avatar}"
    alt={$LL.avatarAltText()}
    class={klass}
    title={userName}
    {...options}
  />
{:else if AvatarService === 'goadorable'}
  <img
    src="{PathPrefix}/avatar/{width}/{warriorId}"
    alt={$LL.avatarAltText()}
    class={klass}
    title={userName}
    {...options}
  />
{/if}
