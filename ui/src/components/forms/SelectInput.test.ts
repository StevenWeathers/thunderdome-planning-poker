import { describe, it, expect } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import SelectInputWithOptions from './__tests__/SelectInputWithOptions.svelte';

const renderWithOptions = (props: Record<string, any> = {}) => {
  return render(SelectInputWithOptions, props);
};

describe('SelectInput component', () => {
  it('should render successfully', () => {
    renderWithOptions({ name: 'example' });

    const select = page.getByRole('combobox');
    expect(select).toBeTruthy();
  });

  it('should render options from slot', () => {
    const { container } = renderWithOptions();

    const options = Array.from(container.querySelectorAll('option'));
    expect(options.length).toBe(3);
    expect(options.map(option => option.getAttribute('value'))).toEqual(['', 'one', 'two']);
  });

  it('should update value when selection changes', async () => {
    const { container } = renderWithOptions({ name: 'choice' });

    const select = container.querySelector('select[name="choice"]') as HTMLSelectElement;

    await userEvent.selectOptions(select, 'two');

    await expect.element(select).toHaveValue('two');
  });

  it('should pass through attributes to the select element', () => {
    const { container } = renderWithOptions({
      id: 'select-id',
      name: 'select-name',
      required: true,
      disabled: true,
    });

    const select = container.querySelector('select[name="select-name"]') as HTMLSelectElement;
    expect(select.getAttribute('id')).toBe('select-id');
    expect(select.getAttribute('required')).toBe('');
    expect(select.disabled).toBe(true);
  });

  it('should be focusable via component method', () => {
    const { component } = renderWithOptions({ name: 'focus-select' });

    component.focus();

    const select = page.getByRole('combobox');
    const element = select.element() as HTMLSelectElement;
    expect(document.activeElement).toBe(element);
  });
});
