import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import MultiSelectInputWithOptions from './__tests__/MultiSelectInputWithOptions.svelte';

const renderWithOptions = (props: Record<string, any> = {}) => {
  return render(MultiSelectInputWithOptions, props);
};

describe('MultiSelectInput component', () => {
  it('should render successfully', () => {
    renderWithOptions({ name: 'multi-example' });

    const select = page.getByRole('listbox');
    expect(select).toBeTruthy();
  });

  it('should update value when multiple selections change', async () => {
    const { container } = renderWithOptions({ name: 'choices' });

    const select = container.querySelector('select[name="choices"]') as HTMLSelectElement;

    await userEvent.selectOptions(select, ['one', 'three']);

    expect(Array.from(select.selectedOptions, option => option.value)).toEqual(['one', 'three']);
  });

  it('should reflect initial selected values', async () => {
    const { container } = renderWithOptions({
      name: 'choices',
      value: ['two', 'three'],
    });

    const select = container.querySelector('select[name="choices"]') as HTMLSelectElement;

    await expect.poll(() => Array.from(select.selectedOptions, option => option.value)).toEqual(['two', 'three']);
  });

  it('should call onchange handler when selection changes', async () => {
    const onchange = vi.fn();
    const { container } = renderWithOptions({ name: 'choices', onchange });

    const select = container.querySelector('select[name="choices"]') as HTMLSelectElement;

    await userEvent.selectOptions(select, ['two']);

    expect(onchange).toHaveBeenCalled();
  });

  it('should be focusable via component method', () => {
    const { component } = renderWithOptions({ name: 'focus-multi-select' });

    component.focus();

    const select = page.getByRole('listbox');
    const element = select.element() as HTMLSelectElement;
    expect(document.activeElement).toBe(element);
  });
});
