<script lang="ts">
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';
  import { onMount } from 'svelte';
  import AssociatedTeam from './AssociatedTeam.svelte';
  import AssociatedOrganization from './AssociatedOrganization.svelte';
  import HollowButton from '../global/HollowButton.svelte';

  export let eventTag;
  export let notifications;
  export let xfetch = async (url: string, ...options: any) => {};

  let userSubscriptions = [];

  function getApiKeys() {
    xfetch(`/api/users/${$user.id}/subscriptions`)
      .then((res: any) => res.json())
      .then(function (result) {
        userSubscriptions = result.data;
      })
      .catch(function () {
        notifications.danger('Error getting subscriptions');
        eventTag('fetch_users_subscriptions', 'engagement', 'failure');
      });
  }

  onMount(() => {
    getApiKeys();
  });
</script>

<div class="flex flex-col">
  <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
    <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
      <div
        class="shadow overflow-hidden border-b border-gray-200 dark:border-gray-700 sm:rounded-lg"
      >
        {#if $user.subscribed && userSubscriptions.length === 0}
          <p class="text-green-600 dark:text-lime-400 font-semibold">
            Active subscription(s) associated to Organization(s) or Team(s)
          </p>
        {:else}
          <table
            class="min-w-full divide-y divide-gray-200 dark:divide-gray-700"
          >
            <thead class="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th
                  scope="col"
                  class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                >
                  {$LL.type()}
                </th>
                <th
                  scope="col"
                  class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                >
                  {$LL.organization()}
                </th>
                <th
                  scope="col"
                  class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                >
                  {$LL.team()}
                </th>
                <th
                  scope="col"
                  class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                >
                  {$LL.expireDate()}
                </th>
              </tr>
            </thead>
            <tbody
              class="bg-white dark:bg-gray-700 divide-y divide-gray-200 dark:divide-gray-800 dark:text-white"
            >
              {#each userSubscriptions as sub, i}
                <tr
                  class:bg-slate-100="{i % 2 !== 0}"
                  class:dark:bg-gray-800="{i % 2 !== 0}"
                  data-testid="subscriptions"
                  data-apikeyid="{sub.id}"
                >
                  <td class="px-6 py-4 whitespace-nowrap" data-testid="sub-type"
                    >{sub.type}</td
                  >
                  <td class="px-6 py-4 whitespace-nowrap">
                    {#if sub.type === 'organization'}
                      {#if sub.organization_id === ''}
                        <HollowButton>Associate to Organization</HollowButton>
                      {:else}
                        <AssociatedOrganization
                          userId="{$user.id}"
                          organizationId="{sub.organization_id}"
                          xfetch="{xfetch}"
                          eventTag="{eventTag}"
                          notifications="{notifications}"
                        />
                      {/if}
                    {:else}
                      <span class="text-gray-300 dark:text-gray-400">N/A</span>
                    {/if}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    {#if sub.type === 'team'}
                      {#if sub.team_id === ''}
                        <HollowButton>Associate to Team</HollowButton>
                      {:else}
                        <AssociatedTeam
                          userId="{$user.id}"
                          teamId="{sub.team_id}"
                          xfetch="{xfetch}"
                          eventTag="{eventTag}"
                          notifications="{notifications}"
                        />
                      {/if}
                    {:else}
                      <span class="text-gray-300 dark:text-gray-400">N/A</span>
                    {/if}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    {new Date(sub.expires).toLocaleString()}
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        {/if}
      </div>
    </div>
  </div>
</div>
