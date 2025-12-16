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

  // Generate a consistent color based on the user ID
  function generateColorFromId(id: string): string {
    const colors = [
      '#E53935',
      '#D81B60',
      '#8E24AA',
      '#5E35B1',
      '#3949AB',
      '#1E88E5',
      '#039BE5',
      '#00ACC1',
      '#00897B',
      '#43A047',
      '#7CB342',
      '#C0CA33',
      '#FDD835',
      '#FFB300',
      '#FB8C00',
      '#F4511E',
      '#6D4C41',
      '#757575',
      '#546E7A',
    ];
    let hash = 0;
    for (let i = 0; i < id.length; i++) {
      hash = id.charCodeAt(i) + ((hash << 5) - hash);
    }
    return colors[Math.abs(hash) % colors.length];
  }

  // Extract initials from the user name (up to 2 letters)
  function getInitials(name: string): string {
    if (!name) return '?';
    const words = name.trim().split(/\s+/);
    if (words.length === 1) {
      return words[0].substring(0, 2).toUpperCase();
    }
    return (words[0][0] + words[1][0]).toUpperCase();
  }

  let avatarColor = $derived(generateColorFromId(warriorId || userName));
  let initials = $derived(getInitials(userName));
</script>

{#if pictureUrl !== ''}
  <img src={pictureUrl} alt={$LL.avatarAltText()} class={klass} title={userName} {...options} />
{:else if avatar === 'none'}
  <div
    class={klass}
    style="background-color: {avatarColor}; width: {width}px; height: {width}px; display: flex; align-items: center; justify-content: center; font-weight: bold; color: white; font-size: {Math.max(
      width * 0.4,
      12,
    )}px;"
    title={userName}
    {...options}
  >
    {initials}
  </div>
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
