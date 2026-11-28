import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';
import { createRawSnippet } from 'svelte';

import Modal from './Modal.svelte';

const children = createRawSnippet(() => ({ render: () => '<input aria-label="Name" />' }));

describe('Modal Escape key', () => {
  it('closes after other keys have been pressed', async () => {
    const closeModal = vi.fn();
    render(Modal, { closeModal, children, ariaLabel: 'Test' });

    await userEvent.type(page.getByLabelText('Name'), 'abc');
    expect(closeModal).not.toHaveBeenCalled();

    await userEvent.keyboard('{Escape}');
    expect(closeModal).toHaveBeenCalledOnce();
  });

  it('saves the edit in the focused field before closing', async () => {
    const calls: string[] = [];
    render(Modal, { closeModal: () => calls.push('close'), children, ariaLabel: 'Test' });
    const input = page.getByLabelText('Name');
    input.element().addEventListener('change', () => calls.push('change'));

    await userEvent.type(input, 'abc');
    await userEvent.keyboard('{Escape}');

    expect(calls).toEqual(['change', 'close']);
  });

  it('closes only the topmost modal', async () => {
    const closeFirst = vi.fn();
    const closeSecond = vi.fn();
    render(Modal, { closeModal: closeFirst, ariaLabel: 'First', children });
    const second = render(Modal, { closeModal: closeSecond, ariaLabel: 'Second', children });

    await userEvent.keyboard('{Escape}');
    expect(closeSecond).toHaveBeenCalledOnce();
    expect(closeFirst).not.toHaveBeenCalled();

    second.unmount();
    await userEvent.keyboard('{Escape}');
    expect(closeFirst).toHaveBeenCalledOnce();
  });
});

describe('Modal page scroll', () => {
  it('is locked until the last modal closes', () => {
    const pageOverflow = () => document.documentElement.style.overflow;
    const first = render(Modal, { ariaLabel: 'First', children });
    const second = render(Modal, { ariaLabel: 'Second', children });
    expect(pageOverflow()).toBe('hidden');

    second.unmount();
    expect(pageOverflow()).toBe('hidden');

    first.unmount();
    expect(pageOverflow()).toBe('');
  });
});
