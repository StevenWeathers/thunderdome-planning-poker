import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import CommentForm from './CommentForm.svelte';

describe('CommentForm component', () => {
  it('should render successfully', () => {
    render(CommentForm, { onSubmit: vi.fn() });

    const form = page.getByRole('form', { name: 'checkinComment' });
    const textarea = page.getByRole('textbox');

    expect(form).toBeTruthy();
    expect(textarea).toBeTruthy();
  });

  it('should render placeholder text', () => {
    const { container } = render(CommentForm, { onSubmit: vi.fn() });

    const textarea = container.querySelector('textarea');
    expect(textarea?.getAttribute('placeholder')).toBe('Write a comment...');
  });

  it('should disable buttons when comment is empty', () => {
    render(CommentForm, { onSubmit: vi.fn() });

    const clearButton = page.getByRole('button', { name: 'Clear' });
    const submitButton = page.getByRole('button', { name: 'Post Comment' });

    expect(clearButton.element()).toBeDisabled();
    expect(submitButton.element()).toBeDisabled();
  });

  it('should enable buttons when comment is provided', async () => {
    render(CommentForm, { onSubmit: vi.fn() });

    const textarea = page.getByRole('textbox');
    const clearButton = page.getByRole('button', { name: 'Clear' });
    const submitButton = page.getByRole('button', { name: 'Post Comment' });

    await userEvent.fill(textarea, 'Hello world');

    expect(clearButton.element()).not.toBeDisabled();
    expect(submitButton.element()).not.toBeDisabled();
  });

  it('should show and update the character counter', async () => {
    const { container } = render(CommentForm, { onSubmit: vi.fn() });

    const textarea = page.getByRole('textbox');

    await userEvent.fill(textarea, 'Hello');

    const counter = container.querySelector('[data-testid="comment-counter"]');
    expect(counter?.textContent).toBe('5');
  });

  it('should submit comment and clear the textarea', async () => {
    const onSubmit = vi.fn();
    const { container } = render(CommentForm, { onSubmit });

    const textarea = page.getByRole('textbox');
    const submitButton = page.getByRole('button', { name: 'Post Comment' });

    await userEvent.fill(textarea, 'Ready to ship');
    await submitButton.click();

    expect(onSubmit).toHaveBeenCalledWith('Ready to ship');
    expect(onSubmit).toHaveBeenCalledTimes(1);

    await expect.element(textarea).toHaveValue('');
    expect(container.querySelector('[data-testid="comment-counter"]')).toBeNull();
  });

  it('should clear comment when clear button is clicked', async () => {
    const { container } = render(CommentForm, { onSubmit: vi.fn() });

    const textarea = page.getByRole('textbox');
    const clearButton = page.getByRole('button', { name: 'Clear' });

    await userEvent.fill(textarea, 'Clear me');
    await clearButton.click();

    await expect.element(textarea).toHaveValue('');
    expect(clearButton.element()).toBeDisabled();
    expect(container.querySelector('[data-testid="comment-counter"]')).toBeNull();
  });

  it('should not enable submit for whitespace-only comment', async () => {
    const onSubmit = vi.fn();
    render(CommentForm, { onSubmit });

    const textarea = page.getByRole('textbox');
    const submitButton = page.getByRole('button', { name: 'Post Comment' });

    await userEvent.fill(textarea, '   ');

    expect(submitButton.element()).toBeDisabled();
  });
});
