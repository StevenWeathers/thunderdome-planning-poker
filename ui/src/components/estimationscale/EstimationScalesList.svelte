<script lang="ts">
  import Table from '../table/Table.svelte';
  import HeadCol from '../table/HeadCol.svelte';
  import TableRow from '../table/TableRow.svelte';
  import RowCol from '../table/RowCol.svelte';

  import LL from '../../i18n/i18n-svelte';
  import DeleteConfirmation from '../global/DeleteConfirmation.svelte';
  import TableContainer from '../table/TableContainer.svelte';
  import TableNav from '../table/TableNav.svelte';
  import CrudActions from '../table/CrudActions.svelte';
  import { createEventDispatcher } from 'svelte';
  import { validateUserIsAdmin } from '../../validationUtils';
  import { user } from '../../stores';
  import CreateEstimationScale from './CreateEstimationScale.svelte';
  import UpdateEstimationScale from './UpdateEstimationScale.svelte';
  import BooleanDisplay from '../global/BooleanDisplay.svelte';
  import TableFooter from '../table/TableFooter.svelte';

  interface Props {
    xfetch: any;
    notifications: any;
    organizationId: any;
    teamId: any;
    departmentId: any;
    isEntityAdmin?: boolean;
    scales?: any;
    apiPrefix?: string;
    scaleCount?: number;
    scalesPage?: number;
    scalesPageLimit?: number;
    changePage?: any;
    getScales?: any;
  }

  let {
    xfetch,
    notifications,
    organizationId,
    teamId,
    departmentId,
    isEntityAdmin = false,
    scales = [],
    apiPrefix = '/api',
    scaleCount = 0,
    scalesPage = $bindable(1),
    scalesPageLimit = 10,
    changePage = () => {},
    getScales = () => {}
  }: Props = $props();

  const dispatch = createEventDispatcher();

  let showAddScale = $state(false);
  let showUpdateScale = $state(false);
  let updateScale = $state({});
  let showRemoveScale = $state(false);
  let removeScaleId = null;

  function handleCreateScale() {
    getScales();
    toggleCreateScale();
  }

  function handleScaleUpdate() {
    getScales();
    toggleUpdateScale({})();
  }

  function handleScaleRemove() {
    xfetch(`${apiPrefix}/estimation-scales/${removeScaleId}`, {
      method: 'DELETE',
    })
      .then(function () {
        toggleRemoveScale(null)();
        notifications.success($LL.estimationScaleRemoveSuccess());
        getScales();
      })
      .catch(function () {
        notifications.danger($LL.estimationScaleRemoveError());
      });
  }

  function toggleCreateScale() {
    showAddScale = !showAddScale;
  }

  const toggleUpdateScale = scale => () => {
    updateScale = scale;
    showUpdateScale = !showUpdateScale;
  };

  const toggleRemoveScale = scaleId => () => {
    showRemoveScale = !showRemoveScale;
    removeScaleId = scaleId;
  };

  let isAdmin = $derived(validateUserIsAdmin($user));
</script>

<div class="w-full">
  <TableContainer>
    <TableNav
      title={$LL.estimationScales()}
      createBtnEnabled={isAdmin || isEntityAdmin}
      createBtnText={$LL.estimationScaleCreate()}
      createButtonHandler={toggleCreateScale}
      createBtnTestId="scale-create"
    />
    <Table>
      {#snippet header()}
            <tr >
          <HeadCol>{$LL.name()}</HeadCol>
          <HeadCol>{$LL.description()}</HeadCol>
          <HeadCol>{$LL.scaleValues()}</HeadCol>
          {#if isAdmin && !teamId && !organizationId}
            <HeadCol>{$LL.scaleType()}</HeadCol>
            <HeadCol>{$LL.isPublic()}</HeadCol>
          {/if}
          <HeadCol>{$LL.default()}</HeadCol>
          <HeadCol type="action">
            <span class="sr-only">{$LL.actions()}</span>
          </HeadCol>
        </tr>
          {/snippet}
      {#snippet body({ class: className })}
            <tbody   class="{className}">
          {#each scales as scale, i}
            <TableRow itemIndex={i}>
              <RowCol>
                <div class="font-medium text-gray-900 dark:text-gray-200">
                  <span data-testid="scale-name">{scale.name}</span>
                </div>
              </RowCol>
              <RowCol>
                <span data-testid="scale-description">{scale.description}</span>
              </RowCol>
              <RowCol>
                <span data-testid="scale-values">{scale.values.join(', ')}</span>
              </RowCol>
              {#if isAdmin && !teamId && !organizationId}
                <RowCol>
                  <span data-testid="scale-type">{scale.scaleType}</span>
                </RowCol>
                <RowCol>
                  <span data-testid="scale-is-public"
                    ><BooleanDisplay boolValue={scale.isPublic} /></span
                  >
                </RowCol>
              {/if}
              <RowCol>
                <span data-testid="scale-is-default"
                  ><BooleanDisplay boolValue={scale.defaultScale} />
                </span>
              </RowCol>
              <RowCol type="action">
                {#if isAdmin || isEntityAdmin}
                  <CrudActions
                    editBtnClickHandler={toggleUpdateScale(scale)}
                    deleteBtnClickHandler={toggleRemoveScale(scale.id)}
                  />
                {/if}
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
          {/snippet}
    </Table>
    <TableFooter
      bind:current={scalesPage}
      num_items={scaleCount}
      per_page={scalesPageLimit}
      on:navigate={changePage}
    />
  </TableContainer>

  {#if showAddScale}
    <CreateEstimationScale
      toggleCreate={toggleCreateScale}
      handleCreate={handleCreateScale}
      organizationId={organizationId}
      departmentId={departmentId}
      teamId={teamId}
      apiPrefix={apiPrefix}
      xfetch={xfetch}
      notifications={notifications}
    />
  {/if}

  {#if showUpdateScale}
    <UpdateEstimationScale
      toggleUpdate={toggleUpdateScale({})}
      handleUpdate={handleScaleUpdate}
      scaleId={updateScale.id}
      name={updateScale.name}
      description={updateScale.description}
      values={updateScale.values}
      scaleType={updateScale.scaleType}
      isPublic={updateScale.isPublic}
      defaultScale={updateScale.defaultScale}
      organizationId={organizationId}
      departmentId={departmentId}
      teamId={teamId}
      apiPrefix={apiPrefix}
      xfetch={xfetch}
      notifications={notifications}
    />
  {/if}

  {#if showRemoveScale}
    <DeleteConfirmation
      toggleDelete={toggleRemoveScale(null)}
      handleDelete={handleScaleRemove}
      permanent={true}
      confirmText={$LL.removeEstimationScaleConfirmText()}
      confirmBtnText={$LL.removeEstimationScale()}
    />
  {/if}
</div>
