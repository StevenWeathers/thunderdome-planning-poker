<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../stores';
  import PageLayout from '../../components/PageLayout.svelte';
  import { appRoutes } from '../../config';
  import SolidButton from '../../components/global/SolidButton.svelte';
  import { Check, TriangleAlert } from '@lucide/svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  interface Props {
    router: any;
    xfetch: ApiClient;
    notifications: NotificationService;
    inviteType: any;
    inviteId: any;
  }

  let { router, xfetch, notifications, inviteType, inviteId }: Props = $props();

  let targetPage = $state('');

  let inviteDetails = $state({
    id: '',
    name: '',
    role: '',
    organization_id: '',
    department_id: '',
  });
  let inviteProcessed = $state(false);
  let inviteErr = $state('');

  function useInvite() {
    xfetch(`/api/users/${$user.id}/invite/${inviteType}/${inviteId}`, {
      method: 'POST',
      body: {},
    })
      .then(res => res.json())
      .then(function (result) {
        inviteDetails = result.data;
        if (inviteType === 'team') {
          if (inviteDetails.organization_id !== '') {
            targetPage = `${appRoutes.organization}/${inviteDetails.organization_id}/team`;
          }
          if (inviteDetails.department_id !== '') {
            targetPage = `${appRoutes.organization}/${inviteDetails.organization_id}/department/${inviteDetails.department_id}/team`;
          }
        }
        if (inviteType === 'department') {
          targetPage = `${appRoutes.organization}/${inviteDetails.organization_id}/department`;
        }
        if (inviteType === 'organization') {
          targetPage = `${appRoutes.organization}`;
        }
        inviteProcessed = true;
      })
      .catch(function (error) {
        if (Array.isArray(error)) {
          error[1].json().then(function (result: any) {
            inviteErr = result.error;
            inviteProcessed = true;
          });
        } else {
          inviteErr = 'Internal Error';
          inviteProcessed = true;
        }
      });
  }

  onMount(() => {
    if ($user.id && $user.rank !== 'GUEST') {
      useInvite();
    } else {
      router.route(`${appRoutes.login}/${inviteType}/${inviteId}`, true);
    }
  });
</script>

<PageLayout>
  <div class="flex justify-center">
    <div
      class="w-full max-w-lg relative p-4 lg:p-8 text-center bg-white rounded-lg shadow dark:bg-gray-800 sm:p-5 text-center mt-8 lg:mt-12"
    >
      {#if inviteProcessed}
        {#if inviteErr !== ''}
          <div
            class="w-12 h-12 rounded-full bg-red-100 dark:bg-red-900 p-2 flex items-center justify-center mx-auto mb-3.5"
          >
            <TriangleAlert class="w-8 h-8 text-red-500 dark:text-red-400" />
            <span class="sr-only">Failure</span>
          </div>
          <div class="text-lg font-semibold text-gray-900 dark:text-white">
            Error processing invite<br />
            <span class="text-red-500 dark:text-red-400">{inviteErr}</span>
          </div>
        {:else}
          <div
            class="w-12 h-12 rounded-full bg-green-100 dark:bg-green-900 p-2 flex items-center justify-center mx-auto mb-3.5"
          >
            <Check class="inline-block w-8 h-8 text-green-500 dark:text-green-400" />
            <span class="sr-only">Success</span>
          </div>
          <div class="mb-4 text-lg font-semibold text-gray-900 dark:text-white">
            You have successfully joined the {inviteType}<br />
            {inviteDetails.name}
          </div>
          <div>
            <SolidButton color="blue" href="{targetPage}/{inviteDetails.id}">Continue</SolidButton>
          </div>
        {/if}
      {/if}
    </div>
  </div>
</PageLayout>
