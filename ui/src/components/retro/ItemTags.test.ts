import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import ItemTags from './ItemTags.svelte';
import { buildTagSuggestions } from './retroTags';

const setup = (props: { tags?: string[]; disabled?: boolean } = {}) => {
  const sendSocketEvent = vi.fn();
  render(ItemTags, { itemId: 'item-1', sendSocketEvent, ...props });

  return { sendSocketEvent };
};

describe('ItemTags component', () => {
  it('adds a custom tag with collapsed whitespace', async () => {
    const { sendSocketEvent } = setup();

    await userEvent.click(page.getByTestId('retro-item-tag-toggle'));
    await userEvent.fill(page.getByLabelText('Search or create a tag'), '  Waiting on   vendor ');
    await userEvent.click(page.getByTestId('retro-item-tag-create'));

    expect(sendSocketEvent).toHaveBeenCalledWith(
      'item_tag_add',
      JSON.stringify({ item_id: 'item-1', tag: 'Waiting on vendor' }),
    );
  });

  it('removes a tag', async () => {
    const { sendSocketEvent } = setup({ tags: ['Blocked'] });

    await userEvent.click(page.getByLabelText('Remove tag Blocked'));

    expect(sendSocketEvent).toHaveBeenCalledWith(
      'item_tag_remove',
      JSON.stringify({ item_id: 'item-1', tag: 'Blocked' }),
    );
  });

  it('does not add a tag the item already has in another case', async () => {
    const { sendSocketEvent } = setup({ tags: ['Blocked'] });

    await userEvent.click(page.getByTestId('retro-item-tag-toggle'));
    await userEvent.fill(page.getByLabelText('Search or create a tag'), 'blocked');
    await userEvent.keyboard('{Enter}');

    expect(sendSocketEvent).not.toHaveBeenCalled();
  });

  it('shows tags without editing controls when disabled', () => {
    setup({ tags: ['Blocked'], disabled: true });

    expect(page.getByTestId('retro-item-tag').element()).toHaveTextContent('Blocked');
    expect(page.getByTestId('retro-item-tag-toggle').query()).toBeNull();
    expect(page.getByLabelText('Remove tag Blocked').query()).toBeNull();
  });
});

describe('buildTagSuggestions', () => {
  it('merges presets with tags in use, ignoring case', () => {
    const items = [
      { id: '1', content: '', groupId: '', type: '', userId: '', tags: ['blocked', 'Vendor'] },
      { id: '2', content: '', groupId: '', type: '', userId: '' },
    ];

    expect(buildTagSuggestions(['Quick win', 'Blocked'], items)).toEqual(['Blocked', 'Quick win', 'Vendor']);
  });
});
