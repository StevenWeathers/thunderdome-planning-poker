import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import AddUser from './AddUser.svelte';

import type { NotificationService } from '../../types/notifications';
import type { ApiClient } from '../../types/apiclient';

function createNotifications(): NotificationService {
  return {
    show: vi.fn(),
    success: vi.fn(),
    danger: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
    removeToast: vi.fn(),
  };
}

function createJsonResponse(data: unknown): Response {
  return new Response(JSON.stringify(data), {
    status: 200,
    headers: {
      'Content-Type': 'application/json',
    },
  });
}

describe('AddUser component', () => {
  it('should disable submit button when no users or emails are provided', () => {
    const { container } = render(AddUser, {
      toggleAdd: vi.fn(),
      handleAdd: vi.fn(),
      handleInvite: vi.fn(),
      notifications: createNotifications(),
      xfetch: vi.fn() as unknown as ApiClient,
    });

    const button = container.querySelector('button[type="submit"]') as HTMLButtonElement;
    expect(button.disabled).toBe(true);
  });

  it('should submit multiple invite emails at once', async () => {
    const handleInvite = vi.fn();
    const { container } = render(AddUser, {
      toggleAdd: vi.fn(),
      handleAdd: vi.fn(),
      handleInvite,
      notifications: createNotifications(),
      xfetch: vi.fn() as unknown as ApiClient,
    });

    const roleSelect = container.querySelector('select[name="userRole"]') as HTMLSelectElement;
    const textarea = container.querySelector('textarea[name="userEmail"]') as HTMLTextAreaElement;
    const button = page.getByRole('button', { name: /add user/i });

    await userEvent.selectOptions(roleSelect, 'MEMBER');
    await userEvent.fill(textarea, 'first@example.com\nsecond@example.com, third@example.com');
    await userEvent.click(button);

    expect(handleInvite).toHaveBeenCalledWith([
      { email: 'first@example.com', role: 'MEMBER' },
      { email: 'second@example.com', role: 'MEMBER' },
      { email: 'third@example.com', role: 'MEMBER' },
    ]);
    expect(handleInvite).toHaveBeenCalledTimes(1);
  });

  it('should block submit when any invite value is not a valid email', async () => {
    const handleInvite = vi.fn();
    const { container } = render(AddUser, {
      toggleAdd: vi.fn(),
      handleAdd: vi.fn(),
      handleInvite,
      notifications: createNotifications(),
      xfetch: vi.fn() as unknown as ApiClient,
    });

    const roleSelect = container.querySelector('select[name="userRole"]') as HTMLSelectElement;
    const textarea = container.querySelector('textarea[name="userEmail"]') as HTMLTextAreaElement;
    const button = page.getByRole('button', { name: /add user/i });

    await userEvent.selectOptions(roleSelect, 'MEMBER');
    await userEvent.fill(textarea, 'valid@example.com\nnot-an-email');
    await userEvent.click(button);

    expect(handleInvite).not.toHaveBeenCalled();
    expect(textarea.validationMessage).toBe('Please enter valid email addresses only.');

    const invalidMessage = container.querySelector('[data-testid="invalid-email-message"]');
    expect(invalidMessage?.textContent).toContain('not-an-email');
  });

  it('should submit multiple selected users at once', async () => {
    const handleAdd = vi.fn();
    const xfetch: ApiClient = vi.fn((url: string) => {
      if (url.includes('/departments/')) {
        return Promise.resolve(
          createJsonResponse({ data: [{ id: 'dept-user-1', name: 'Dept User', email: 'dept@example.com' }] }),
        );
      }

      return Promise.resolve(
        createJsonResponse({ data: [{ id: 'org-user-1', name: 'Org User', email: 'org@example.com' }] }),
      );
    }) as unknown as ApiClient;

    const { container } = render(AddUser, {
      toggleAdd: vi.fn(),
      handleAdd,
      handleInvite: vi.fn(),
      pageType: 'team',
      orgId: 'org-1',
      deptId: 'dept-1',
      notifications: createNotifications(),
      xfetch,
    });

    const roleSelect = container.querySelector('select[name="userRole"]') as HTMLSelectElement;
    const orgSelect = container.querySelector('select[name="orgUser"]') as HTMLSelectElement;
    const deptSelect = container.querySelector('select[name="deptUser"]') as HTMLSelectElement;
    const button = page.getByRole('button', { name: /add user/i });

    await userEvent.selectOptions(roleSelect, 'ADMIN');
    await userEvent.selectOptions(orgSelect, ['org-user-1']);
    await userEvent.selectOptions(deptSelect, ['dept-user-1']);
    await userEvent.click(button);

    expect(handleAdd).toHaveBeenCalledWith([
      { user_id: 'org-user-1', role: 'ADMIN' },
      { user_id: 'dept-user-1', role: 'ADMIN' },
    ]);
    expect(handleAdd).toHaveBeenCalledTimes(1);
  });
});
