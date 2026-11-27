import { describe, it, expect, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';

import { user } from '../../stores';
import GroupPhase from './GroupPhase.svelte';

const item = (id: string, groupId: string) => ({
  id,
  type: 'worked_well',
  content: `Item ${id}`,
  userId: 'user-1',
  groupId,
  comments: [],
  reactions: [],
});

const group = (id: string, itemIds: string[] = [], name = '') => ({
  id,
  name,
  items: itemIds.map(itemId => item(itemId, id)),
});

const setup = (initialGroups: ReturnType<typeof group>[]) => {
  // deeply reactive, as it is when passed in by the retro page
  const groups = $state(initialGroups);
  user.create({ id: 'user-1', name: 'Alex Doe', rank: 'warrior', locale: 'en' } as any);

  const handleItemChange = vi.fn();
  const { container } = render(GroupPhase, { groups, handleItemChange });

  const shownIds = () =>
    Array.from(container.querySelectorAll<HTMLElement>('[data-testid="retro-group"] [data-groupid]')).map(
      el => el.dataset.groupid,
    );
  const newGroup = () => {
    const tile = container.querySelector('[data-testid="retro-new-group"]');

    return tile
      ? {
          id: tile.querySelector<HTMLElement>('[data-groupid]')!.dataset.groupid,
          visible: !tile.classList.contains('invisible'),
        }
      : undefined;
  };
  const dispatch = (
    type: 'consider' | 'finalize',
    groupId: string,
    items: any[],
    itemId: string,
    trigger = '',
    source = 'pointer',
  ) =>
    container
      .querySelector(`[data-groupid="${groupId}"]`)!
      .dispatchEvent(new CustomEvent(type, { detail: { items, info: { id: itemId, trigger, source } } }));

  return { handleItemChange, shownIds, newGroup, dispatch };
};

describe('GroupPhase', () => {
  it('hides empty groups, keeping one invisible as the new group target', () => {
    const { shownIds, newGroup } = setup([group('a'), group('b', ['1']), group('c'), group('d', ['2', '3'])]);

    expect(shownIds()).toEqual(['b', 'd']);
    expect(newGroup()).toEqual({ id: 'c', visible: false });
  });

  it('keeps empty groups that have been named', () => {
    const { shownIds } = setup([group('a', [], 'Tooling'), group('b', ['1']), group('c')]);

    expect(shownIds()).toEqual(['a', 'b']);
  });

  it('offers the new group target while dragging and moves the item into it', async () => {
    const { handleItemChange, shownIds, newGroup, dispatch } = setup([group('a', ['1', '2']), group('b')]);

    dispatch('consider', 'a', [item('2', 'a')], '1');
    await expect.poll(() => newGroup()?.visible).toBe(true);

    dispatch('finalize', 'b', [item('1', 'b')], '1');
    dispatch('finalize', 'a', [item('2', 'a')], '1');

    await expect.poll(shownIds).toEqual(['a', 'b']);
    expect(newGroup()).toBeUndefined();
    expect(handleItemChange).toHaveBeenCalledExactlyOnceWith('1', 'b');
  });

  it('removes an emptied group only after the drag is finished', async () => {
    const { shownIds, newGroup, dispatch } = setup([group('a', ['1']), group('b', ['2']), group('c')]);

    dispatch('consider', 'a', [], '1');
    await expect.poll(() => newGroup()?.visible).toBe(true);
    expect(shownIds()).toEqual(['a', 'b']);

    dispatch('finalize', 'b', [item('2', 'b'), item('1', 'b')], '1');
    dispatch('finalize', 'a', [], '1');

    await expect.poll(shownIds).toEqual(['b']);
    expect(newGroup()?.visible).toBe(false);
  });

  it('stays in drag mode between keyboard moves', async () => {
    const { newGroup, dispatch } = setup([group('a', ['1', '2']), group('b', ['3']), group('c')]);

    dispatch('consider', 'a', [item('1', 'a'), item('2', 'a')], '1', 'dragStarted', 'keyboard');
    dispatch('finalize', 'b', [item('3', 'b'), item('1', 'b')], '1', 'droppedIntoZone', 'keyboard');
    await expect.poll(() => newGroup()?.visible).toBe(true);

    dispatch('consider', 'b', [item('3', 'b'), item('1', 'b')], '1', 'dragStopped', 'keyboard');
    await expect.poll(() => newGroup()?.visible).toBe(false);
  });
});
