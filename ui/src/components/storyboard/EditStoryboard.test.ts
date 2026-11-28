import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import EditStoryboard from './EditStoryboard.svelte';

describe('EditStoryboard component', () => {
  it('should render successfully', () => {
    const toggleEditStoryboard = vi.fn();
    const handleStoryboardEdit = vi.fn();
    render(EditStoryboard, {
      toggleEditStoryboard,
      handleStoryboardEdit,
    });

    const form = page.getByRole('form', { name: 'createStoryboard' });
    expect(form).toBeTruthy();
  });

  it('should render the modal dialog', () => {
    const toggleEditStoryboard = vi.fn();
    const handleStoryboardEdit = vi.fn();
    const { container } = render(EditStoryboard, {
      toggleEditStoryboard,
      handleStoryboardEdit,
    });

    const modal = container.querySelector('[role="dialog"]');
    expect(modal).toBeTruthy();
  });

  it('should render storyboard name input field', () => {
    const toggleEditStoryboard = vi.fn();
    const handleStoryboardEdit = vi.fn();
    const { container } = render(EditStoryboard, {
      toggleEditStoryboard,
      handleStoryboardEdit,
    });

    const input = container.querySelector('input[name="storyboardName"]');
    expect(input).toBeTruthy();
    expect(input?.getAttribute('id')).toBe('storyboardName');
  });

  it('should render join code input field', () => {
    const toggleEditStoryboard = vi.fn();
    const handleStoryboardEdit = vi.fn();
    const { container } = render(EditStoryboard, {
      toggleEditStoryboard,
      handleStoryboardEdit,
    });

    const input = container.querySelector('input[name="joinCode"]');
    expect(input).toBeTruthy();
    expect(input?.getAttribute('id')).toBe('joinCode');
  });

  it('should render facilitator code input field', () => {
    const toggleEditStoryboard = vi.fn();
    const handleStoryboardEdit = vi.fn();
    const { container } = render(EditStoryboard, {
      toggleEditStoryboard,
      handleStoryboardEdit,
    });

    const input = container.querySelector('input[name="facilitatorCode"]');
    expect(input).toBeTruthy();
    expect(input?.getAttribute('id')).toBe('facilitatorCode');
  });

  it('should fire handleStoryboardEdit with form values when submitted', async () => {
    const toggleEditStoryboard = vi.fn();
    const handleStoryboardEdit = vi.fn();
    const { container } = render(EditStoryboard, {
      toggleEditStoryboard,
      handleStoryboardEdit,
      storyboardName: 'Initial Storyboard',
      joinCode: 'join-1',
      facilitatorCode: 'fac-1',
    });

    const nameInput = container.querySelector('input[name="storyboardName"]') as HTMLInputElement;
    const joinInput = container.querySelector('input[name="joinCode"]') as HTMLInputElement;
    const facInput = container.querySelector('input[name="facilitatorCode"]') as HTMLInputElement;
    const button = page.getByRole('button', { name: /save/i });

    await userEvent.clear(nameInput);
    await userEvent.fill(nameInput, 'Updated Storyboard');
    await userEvent.clear(joinInput);
    await userEvent.fill(joinInput, 'join-2');
    await userEvent.clear(facInput);
    await userEvent.fill(facInput, 'fac-2');
    await button.click();

    expect(handleStoryboardEdit).toHaveBeenCalledWith({
      storyboardName: 'Updated Storyboard',
      joinCode: 'join-2',
      facilitatorCode: 'fac-2',
    });
    expect(handleStoryboardEdit).toHaveBeenCalledTimes(1);
  });

  it('should submit form on Enter key when input is focused', async () => {
    const toggleEditStoryboard = vi.fn();
    const handleStoryboardEdit = vi.fn();
    const { container } = render(EditStoryboard, {
      toggleEditStoryboard,
      handleStoryboardEdit,
    });

    const nameInput = container.querySelector('input[name="storyboardName"]') as HTMLInputElement;

    await userEvent.fill(nameInput, 'Keyboard Submit');
    await userEvent.keyboard('{Enter}');

    expect(handleStoryboardEdit).toHaveBeenCalledWith({
      storyboardName: 'Keyboard Submit',
      joinCode: '',
      facilitatorCode: '',
    });
    expect(handleStoryboardEdit).toHaveBeenCalledTimes(1);
  });

  it('should update input values when changed', async () => {
    const toggleEditStoryboard = vi.fn();
    const handleStoryboardEdit = vi.fn();
    const { container } = render(EditStoryboard, {
      toggleEditStoryboard,
      handleStoryboardEdit,
    });

    const nameInput = container.querySelector('input[name="storyboardName"]') as HTMLInputElement;
    const joinInput = container.querySelector('input[name="joinCode"]') as HTMLInputElement;
    const facInput = container.querySelector('input[name="facilitatorCode"]') as HTMLInputElement;

    await userEvent.fill(nameInput, 'New Name');
    await userEvent.fill(joinInput, 'join-code');
    await userEvent.fill(facInput, 'fac-code');

    await expect.element(nameInput).toHaveValue('New Name');
    await expect.element(joinInput).toHaveValue('join-code');
    await expect.element(facInput).toHaveValue('fac-code');
  });

  it.each([
    ['another scale', 'scale-2'],
    ['no scale', ''],
  ])('should submit a change to %s', async (_, selected) => {
    const handleStoryboardEdit = vi.fn();
    const current = { id: 'scale-1', name: 'Sizes', values: ['S', 'M'] };
    const { container } = render(EditStoryboard, {
      handleStoryboardEdit,
      storyboardName: 'Board',
      estimationScale: current,
      getEstimationScales: () => Promise.resolve([current, { id: 'scale-2', name: 'Fibonacci', values: ['1', '2'] }]),
    });
    const select = container.querySelector('select[name="estimationScale"]') as HTMLSelectElement;

    await expect.element(select).toHaveValue('scale-1');
    await vi.waitFor(() => expect(select.options.length).toBe(3));
    await userEvent.selectOptions(select, selected);
    await page.getByRole('button', { name: /save/i }).click();

    expect(handleStoryboardEdit).toHaveBeenCalledWith({
      storyboardName: 'Board',
      joinCode: '',
      facilitatorCode: '',
      estimationScaleId: selected,
    });
  });

  it('should keep the current scale selected when the scales fail to load', async () => {
    const onEstimationScalesError = vi.fn();
    const { container } = render(EditStoryboard, {
      estimationScale: { id: 'scale-1', name: 'Sizes', values: ['S', 'M'] },
      getEstimationScales: () => Promise.reject(new Error('boom')),
      onEstimationScalesError,
    });

    await vi.waitFor(() => expect(onEstimationScalesError).toHaveBeenCalled());
    await expect
      .element(container.querySelector('select[name="estimationScale"]') as HTMLSelectElement)
      .toHaveValue('scale-1');
  });
});
