import { describe, it, expect, vi } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import Personas from './Personas.svelte';

const template = {
  id: 'template-1',
  name: 'Online Store',
  description: 'Common users of an e-commerce site',
  isPublic: true,
  personas: [{ name: 'Shopper', role: 'Visitor' }, { name: 'Manager' }],
};

const loadButton = () => page.getByRole('button', { name: 'Load from Template' });

describe('Personas component templates', () => {
  it('lets a facilitator load the personas of a template', async () => {
    const onLoadTemplate = vi.fn();
    render(Personas, { isFacilitator: true, loadTemplates: () => Promise.resolve([template]), onLoadTemplate });

    await loadButton().click();
    await page.getByRole('button', { name: 'Add 2 personas' }).click();

    expect(onLoadTemplate).toHaveBeenCalledWith('template-1');
    expect(page.getByTestId('persona-template-option').elements()).toHaveLength(0);
  });

  it('counts only the personas not yet on the storyboard', async () => {
    render(Personas, {
      isFacilitator: true,
      personas: [{ id: 'p1', name: 'Shopper' }],
      loadTemplates: () => Promise.resolve([template]),
    });

    await loadButton().click();

    await expect.element(page.getByRole('button', { name: 'Add 1 persona' })).toBeEnabled();
  });

  it('disables a template whose personas are all on the storyboard', async () => {
    render(Personas, {
      isFacilitator: true,
      personas: [
        { id: 'p1', name: 'Shopper' },
        { id: 'p2', name: 'Manager' },
      ],
      loadTemplates: () => Promise.resolve([template]),
    });

    await loadButton().click();

    await expect.element(page.getByRole('button', { name: 'Already added' })).toBeDisabled();
  });

  it('reports a failure to get templates', async () => {
    render(Personas, { isFacilitator: true, loadTemplates: () => Promise.reject(new Error('nope')) });

    await loadButton().click();

    await expect.element(page.getByRole('alert')).toHaveTextContent('Failed to get persona templates');
  });

  it('does not offer templates to participants', () => {
    render(Personas, { isFacilitator: false, loadTemplates: () => Promise.resolve([template]) });

    expect(loadButton().elements()).toHaveLength(0);
  });
});
