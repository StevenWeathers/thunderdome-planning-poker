import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import ColumnPersonasModal from './ColumnPersonasModal.svelte';

describe('ColumnPersonasModal component', () => {
  const mockColumn = {
    id: 'column-1',
    name: 'Testing',
    personas: [],
  };

  const mockPersonas = [
    { id: 'persona-1', name: 'Alice', role: 'Designer' },
    { id: 'persona-2', name: 'Bob', role: 'Developer' },
    { id: 'persona-3', name: 'Charlie', role: 'Product Manager' },
  ];

  it('should render the modal successfully', () => {
    const { container } = render(ColumnPersonasModal, {
      onPersonaAdd: vi.fn(), // Provide a no-op for add to avoid errors
      onPersonaRemove: vi.fn(), // Provide a no-op for remove to avoid errors
      column: mockColumn,
      personas: mockPersonas,
    });

    const modal = container.querySelector('[role="dialog"]');
    expect(modal).toBeTruthy();
  });

  it('should display the section title', () => {
    const { container } = render(ColumnPersonasModal, {
      onPersonaAdd: vi.fn(), // Provide a no-op for add to avoid errors
      onPersonaRemove: vi.fn(), // Provide a no-op for remove to avoid errors
      column: mockColumn,
      personas: mockPersonas,
    });

    const title = container.querySelector('.sectionTitle');
    expect(title?.textContent).toContain('Column Personas');
  });

  it('should show empty state when no personas are added to the column', () => {
    const { container } = render(ColumnPersonasModal, {
      onPersonaAdd: vi.fn(), // Provide a no-op for add to avoid errors
      onPersonaRemove: vi.fn(), // Provide a no-op for remove to avoid errors
      column: mockColumn,
      personas: mockPersonas,
    });

    const emptyState = Array.from(container.querySelectorAll('.emptyState')).find(el =>
      el.textContent?.includes('No personas added to this column'),
    );
    expect(emptyState).toBeTruthy();
  });

  it('should display all available personas when none are added', () => {
    const { container } = render(ColumnPersonasModal, {
      onPersonaAdd: vi.fn(), // Provide a no-op for add to avoid errors
      onPersonaRemove: vi.fn(), // Provide a no-op for remove to avoid errors
      column: mockColumn,
      personas: mockPersonas,
    });

    mockPersonas.forEach(persona => {
      const personaElement = Array.from(container.querySelectorAll('.personaName')).find(
        el => el.textContent === persona.name,
      );
      expect(personaElement).toBeTruthy();
    });
  });

  it('should display persona role correctly', () => {
    const { container } = render(ColumnPersonasModal, {
      onPersonaAdd: vi.fn(), // Provide a no-op for add to avoid errors
      onPersonaRemove: vi.fn(), // Provide a no-op for remove to avoid errors
      column: mockColumn,
      personas: mockPersonas,
    });

    mockPersonas.forEach(persona => {
      const roleElement = Array.from(container.querySelectorAll('.personaRole')).find(
        el => el.textContent === persona.role,
      );
      expect(roleElement).toBeTruthy();
    });
  });

  it('should call onPersonaAdd when adding a persona', async () => {
    const onPersonaAdd = vi.fn();
    const { container } = render(ColumnPersonasModal, {
      onPersonaRemove: vi.fn(), // Provide a no-op for remove to avoid errors
      column: mockColumn,
      personas: mockPersonas,
      onPersonaAdd,
    });

    // Find the first persona card hoverable (available persona) and click its button
    const personaCards = container.querySelectorAll('.personaCardHoverable');
    const firstCardButton = personaCards[0].querySelector('button');

    if (firstCardButton) {
      await userEvent.click(firstCardButton);
    }

    expect(onPersonaAdd).toHaveBeenCalledWith({
      column_id: 'column-1',
      persona_id: 'persona-1',
    });
  });

  it('should display added personas in the column personas section', () => {
    const columnWithPersonas = {
      id: 'column-1',
      name: 'Testing',
      personas: [mockPersonas[0], mockPersonas[1]],
    };

    const { container } = render(ColumnPersonasModal, {
      onPersonaAdd: vi.fn(), // Provide a no-op for add to avoid errors
      onPersonaRemove: vi.fn(), // Provide a no-op for remove to avoid errors
      column: columnWithPersonas,
      personas: mockPersonas,
    });

    expect(container.textContent).toContain('Alice');
    expect(container.textContent).toContain('Bob');
    expect(container.textContent).not.toContain('No personas added to this column');
  });

  it('should call onPersonaRemove when removing a persona', async () => {
    const onPersonaRemove = vi.fn();
    const columnWithPersonas = {
      id: 'column-1',
      name: 'Testing',
      personas: [mockPersonas[0], mockPersonas[1]],
    };

    const { container } = render(ColumnPersonasModal, {
      onPersonaAdd: vi.fn(), // Provide a no-op for add to avoid errors
      column: columnWithPersonas,
      personas: mockPersonas,
      onPersonaRemove,
    });

    // Find the first persona card (added personas) and click its delete button
    const personaCards = container.querySelectorAll('.personaCard');
    const firstCardButton = personaCards[0].querySelector('button');

    if (firstCardButton) {
      await userEvent.click(firstCardButton);
    }

    expect(onPersonaRemove).toHaveBeenCalledWith({
      column_id: 'column-1',
      persona_id: 'persona-1',
    });
  });

  it('should filter out added personas from available personas list', () => {
    const columnWithPersonas = {
      id: 'column-1',
      name: 'Testing',
      personas: [mockPersonas[0]],
    };

    const { container } = render(ColumnPersonasModal, {
      onPersonaAdd: vi.fn(), // Provide a no-op for add to avoid errors
      onPersonaRemove: vi.fn(), // Provide a no-op for remove to avoid errors
      column: columnWithPersonas,
      personas: mockPersonas,
    });

    // Get all persona name elements
    const personaNames = Array.from(container.querySelectorAll('.personaName')).map(el => el.textContent);

    // Alice should appear once (in added section)
    expect(personaNames.filter(name => name === 'Alice').length).toBe(1);
    // Bob and Charlie should appear once (in available section)
    expect(personaNames.includes('Bob')).toBe(true);
    expect(personaNames.includes('Charlie')).toBe(true);
  });

  it('should show "All personas have been added" message when all are added', () => {
    const columnWithAllPersonas = {
      id: 'column-1',
      name: 'Testing',
      personas: mockPersonas,
    };

    const { container } = render(ColumnPersonasModal, {
      onPersonaAdd: vi.fn(), // Provide a no-op for add to avoid errors
      onPersonaRemove: vi.fn(), // Provide a no-op for remove to avoid errors
      column: columnWithAllPersonas,
      personas: mockPersonas,
    });

    const emptyState = Array.from(container.querySelectorAll('.emptyState')).find(el =>
      el.textContent?.includes('All personas have been added'),
    );
    expect(emptyState).toBeTruthy();
  });

  it('should show "No personas available" when personas array is empty', () => {
    const { container } = render(ColumnPersonasModal, {
      onPersonaAdd: vi.fn(), // Provide a no-op for add to avoid errors
      onPersonaRemove: vi.fn(), // Provide a no-op for remove to avoid errors
      column: mockColumn,
      personas: [],
    });

    const emptyState = Array.from(container.querySelectorAll('.emptyState')).find(el =>
      el.textContent?.includes('No personas available'),
    );
    expect(emptyState).toBeTruthy();
  });

  it('should call toggleModal when modal is closed', async () => {
    const toggleModal = vi.fn();
    const { container } = render(ColumnPersonasModal, {
      onPersonaAdd: vi.fn(), // Provide a no-op for add to avoid errors
      onPersonaRemove: vi.fn(), // Provide a no-op for remove to avoid errors
      column: mockColumn,
      personas: mockPersonas,
      toggleModal,
    });

    // Find the close button in the modal
    const closeButton = container.querySelector('[class*="close"]') || page.getByRole('button', { name: /close/i });

    if (closeButton) {
      await userEvent.click(closeButton);
      expect(toggleModal).toHaveBeenCalled();
    }
  });

  it('should render persona icons for all personas', () => {
    const columnWithPersonas = {
      id: 'column-1',
      name: 'Testing',
      personas: [mockPersonas[0]],
    };

    const { container } = render(ColumnPersonasModal, {
      onPersonaAdd: vi.fn(), // Provide a no-op for add to avoid errors
      onPersonaRemove: vi.fn(), // Provide a no-op for remove to avoid errors
      column: columnWithPersonas,
      personas: mockPersonas,
    });

    const iconContainers = container.querySelectorAll('.personaIconContainer');
    expect(iconContainers.length).toBeGreaterThan(0);
  });

  it('should have correct aria label for the modal', async () => {
    const { container } = render(ColumnPersonasModal, {
      onPersonaAdd: vi.fn(), // Provide a no-op for add to avoid errors
      onPersonaRemove: vi.fn(), // Provide a no-op for remove to avoid errors
      column: mockColumn,
      personas: mockPersonas,
    });

    const modal = container.querySelector('[role="dialog"]');
    expect(modal?.getAttribute('aria-label')).toBeTruthy();
  });

  it('should display section header for available personas', () => {
    const { container } = render(ColumnPersonasModal, {
      onPersonaAdd: vi.fn(), // Provide a no-op for add to avoid errors
      onPersonaRemove: vi.fn(), // Provide a no-op for remove to avoid errors
      column: mockColumn,
      personas: mockPersonas,
    });

    const headers = Array.from(container.querySelectorAll('.sectionHeader'));
    const availableHeader = headers.find(el => el.textContent?.includes('Available'));
    expect(availableHeader).toBeTruthy();
  });

  it('should handle multiple add operations in sequence', async () => {
    const onPersonaAdd = vi.fn();
    const columnWithPersonas = {
      id: 'column-1',
      name: 'Testing',
      personas: [mockPersonas[0]],
    };

    const { container } = render(ColumnPersonasModal, {
      onPersonaRemove: vi.fn(), // Provide a no-op for remove to avoid errors
      column: columnWithPersonas,
      personas: mockPersonas,
      onPersonaAdd,
    });

    // Find available persona cards and click their buttons sequentially
    const personaCards = container.querySelectorAll('.personaCardHoverable');
    const buttons = Array.from(personaCards)
      .map(card => card.querySelector('button'))
      .filter(Boolean) as HTMLButtonElement[];

    // Click the first two add buttons
    if (buttons.length >= 2) {
      await userEvent.click(buttons[0]);
      await userEvent.click(buttons[1]);
    }

    expect(onPersonaAdd).toHaveBeenCalledTimes(2);
  });

  it('should use default values when props are not provided', () => {
    const { container } = render(ColumnPersonasModal, {});

    const modal = container.querySelector('[role="dialog"]');
    expect(modal).toBeTruthy();

    const emptyState = Array.from(container.querySelectorAll('.emptyState')).find(el =>
      el.textContent?.includes('No personas available'),
    );
    expect(emptyState).toBeTruthy();
  });
});
