import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import ColorLegendForm from './ColorLegendForm.svelte';

describe('ColorLegendForm component', () => {
  const legend = [
    { color: 'red', legend: 'High priority' },
    { color: 'blue', legend: 'Low priority' },
  ];

  it('should render successfully', () => {
    const handleLegendRevision = vi.fn();
    const toggleEditLegend = vi.fn();
    render(ColorLegendForm, {
      handleLegendRevision,
      toggleEditLegend,
      colorLegend: legend,
    });

    const form = page.getByRole('form', { name: 'colorLegend' });
    expect(form).toBeTruthy();
  });

  it('should render the modal dialog', () => {
    const handleLegendRevision = vi.fn();
    const toggleEditLegend = vi.fn();
    const { container } = render(ColorLegendForm, {
      handleLegendRevision,
      toggleEditLegend,
      colorLegend: legend,
    });

    const modal = container.querySelector('[role="dialog"]');
    expect(modal).toBeTruthy();
  });

  it('should render inputs for each legend color', () => {
    const handleLegendRevision = vi.fn();
    const toggleEditLegend = vi.fn();
    const { container } = render(ColorLegendForm, {
      handleLegendRevision,
      toggleEditLegend,
      colorLegend: legend,
    });

    const inputs = container.querySelectorAll('input[name^="legend-"]');
    expect(inputs.length).toBe(2);
    expect(container.querySelector('input[name="legend-red"]')).toBeTruthy();
    expect(container.querySelector('input[name="legend-blue"]')).toBeTruthy();
  });

  it('should disable legend inputs when user is not facilitator', () => {
    const handleLegendRevision = vi.fn();
    const toggleEditLegend = vi.fn();
    const { container } = render(ColorLegendForm, {
      handleLegendRevision,
      toggleEditLegend,
      colorLegend: legend,
      isFacilitator: false,
    });

    const inputs = Array.from(container.querySelectorAll('input[name^="legend-"]'));
    expect(inputs.length).toBe(2);
    expect(inputs.every(input => (input as HTMLInputElement).disabled)).toBe(true);
  });

  it('should enable legend inputs when user is facilitator', () => {
    const handleLegendRevision = vi.fn();
    const toggleEditLegend = vi.fn();
    const { container } = render(ColorLegendForm, {
      handleLegendRevision,
      toggleEditLegend,
      colorLegend: legend,
      isFacilitator: true,
    });

    const inputs = Array.from(container.querySelectorAll('input[name^="legend-"]'));
    expect(inputs.length).toBe(2);
    expect(inputs.every(input => !(input as HTMLInputElement).disabled)).toBe(true);
  });

  it('should submit legend updates and close the modal', async () => {
    const handleLegendRevision = vi.fn();
    const toggleEditLegend = vi.fn();
    const { container } = render(ColorLegendForm, {
      handleLegendRevision,
      toggleEditLegend,
      colorLegend: legend,
      isFacilitator: true,
    });

    const input = container.querySelector('input[name="legend-red"]') as HTMLInputElement;
    const button = page.getByRole('button', { name: /save/i });

    await userEvent.clear(input);
    await userEvent.fill(input, 'Updated legend');
    await userEvent.click(button);

    expect(handleLegendRevision).toHaveBeenCalledWith(legend);
    expect(handleLegendRevision).toHaveBeenCalledTimes(1);
    expect(toggleEditLegend).toHaveBeenCalledTimes(1);
  });

  it('should update input value when changed', async () => {
    const handleLegendRevision = vi.fn();
    const toggleEditLegend = vi.fn();
    const { container } = render(ColorLegendForm, {
      handleLegendRevision,
      toggleEditLegend,
      colorLegend: legend,
      isFacilitator: true,
    });

    const input = container.querySelector('input[name="legend-blue"]') as HTMLInputElement;

    await userEvent.clear(input);
    await userEvent.fill(input, 'Changed');

    await expect.element(input).toHaveValue('Changed');
  });
});
