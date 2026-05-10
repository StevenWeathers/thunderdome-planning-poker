import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import ColorLegendForm from './ColorLegendForm.svelte';

describe('ColorLegendForm component', () => {
  const legend = [
    { color: 'red', legend: 'High priority' },
    { color: 'blue', legend: 'Low priority' },
  ];

  const notifications = {
    success: vi.fn(),
    danger: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
  } as any;

  const buildXfetch = ({ subscribed = false, orgId = '', templates = [] as any[], orgTemplates = [] as any[] } = {}) =>
    vi.fn((endpoint: string) => {
      if (endpoint === '/api/teams/team-1') {
        return Promise.resolve({
          json: () =>
            Promise.resolve({
              data: {
                team: { subscribed, organization_id: orgId },
              },
            }),
        } as Response);
      }

      if (endpoint === '/api/teams/team-1/color-legend-templates') {
        return Promise.resolve({
          json: () => Promise.resolve({ data: templates }),
        } as Response);
      }

      if (endpoint === `/api/organizations/${orgId}/color-legend-templates`) {
        return Promise.resolve({
          json: () => Promise.resolve({ data: orgTemplates }),
        } as Response);
      }

      return Promise.reject(new Error(`Unexpected endpoint: ${endpoint}`));
    });

  it('should render successfully', () => {
    const handleLegendRevision = vi.fn();
    const toggleEditLegend = vi.fn();
    render(ColorLegendForm, {
      handleLegendRevision,
      toggleEditLegend,
      colorLegend: legend,
      xfetch: buildXfetch(),
      teamId: 'team-1',
      notifications,
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
      xfetch: buildXfetch(),
      teamId: 'team-1',
      notifications,
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
      xfetch: buildXfetch(),
      teamId: 'team-1',
      notifications,
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
      xfetch: buildXfetch(),
      teamId: 'team-1',
      notifications,
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
      xfetch: buildXfetch(),
      teamId: 'team-1',
      notifications,
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
      xfetch: buildXfetch(),
      teamId: 'team-1',
      notifications,
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
      xfetch: buildXfetch(),
      teamId: 'team-1',
      notifications,
    });

    const input = container.querySelector('input[name="legend-blue"]') as HTMLInputElement;

    await userEvent.clear(input);
    await userEvent.fill(input, 'Changed');

    await expect.element(input).toHaveValue('Changed');
  });

  it('should show the copy from template button only when subscribed', async () => {
    const handleLegendRevision = vi.fn();
    const toggleEditLegend = vi.fn();

    render(ColorLegendForm, {
      handleLegendRevision,
      toggleEditLegend,
      colorLegend: legend,
      isFacilitator: true,
      xfetch: buildXfetch({ subscribed: true }),
      teamId: 'team-1',
      notifications,
    });

    await expect.element(page.getByRole('button', { name: 'Copy From Template' })).toBeVisible();
  });

  it('should copy legend values from a selected template', async () => {
    const handleLegendRevision = vi.fn();
    const toggleEditLegend = vi.fn();
    const xfetch = buildXfetch({
      subscribed: true,
      orgId: 'org-1',
      templates: [
        {
          id: 'team-template',
          name: 'Team Template',
          description: '',
          teamId: 'team-1',
          colorLegend: [
            { color: 'red', legend: 'Urgent' },
            { color: 'blue', legend: 'Backlog' },
          ],
        },
      ],
      orgTemplates: [
        {
          id: 'org-template',
          name: 'Org Template',
          description: 'Shared legend',
          organizationId: 'org-1',
          colorLegend: [
            { color: 'red', legend: 'Critical' },
            { color: 'blue', legend: 'Nice to have' },
          ],
        },
      ],
    });

    const { container } = render(ColorLegendForm, {
      handleLegendRevision,
      toggleEditLegend,
      colorLegend: legend,
      isFacilitator: true,
      xfetch,
      teamId: 'team-1',
      notifications,
    });

    await userEvent.click(page.getByRole('button', { name: 'Copy From Template' }));

    await vi.waitFor(() => {
      expect(page.getByRole('button', { name: 'Select a color legend template...' })).toBeTruthy();
    });
    await userEvent.click(page.getByRole('button', { name: 'Select a color legend template...' }));
    await userEvent.click(page.getByRole('button', { name: /Team: Team Template/i }));

    await expect.element(container.querySelector('input[name="legend-red"]') as HTMLInputElement).toHaveValue('Urgent');
    await expect
      .element(container.querySelector('input[name="legend-blue"]') as HTMLInputElement)
      .toHaveValue('Backlog');
  });
});
