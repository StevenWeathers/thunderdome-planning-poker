import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import UpdatePasswordForm from './UpdatePasswordForm.svelte';

const createNotifications = () => ({
  show: vi.fn(),
  success: vi.fn(),
  danger: vi.fn(),
  warning: vi.fn(),
  info: vi.fn(),
  removeToast: vi.fn(),
});

describe('UpdatePasswordForm component', () => {
  it('should render successfully', () => {
    const handleUpdate = vi.fn();
    const toggleForm = vi.fn();
    render(UpdatePasswordForm, {
      handleUpdate,
      toggleForm,
      notifications: createNotifications(),
    });

    const form = page.getByRole('form', { name: 'updatePassword' });
    expect(form).toBeTruthy();
  });

  it('should render password inputs', () => {
    const { container } = render(UpdatePasswordForm, {
      handleUpdate: vi.fn(),
      toggleForm: vi.fn(),
      notifications: createNotifications(),
    });

    const password1 = container.querySelector('input[name="yourPassword1"]');
    const password2 = container.querySelector('input[name="yourPassword2"]');

    expect(password1).toBeTruthy();
    expect(password2).toBeTruthy();
  });

  it('should disable submit button when fields are empty', () => {
    const { container } = render(UpdatePasswordForm, {
      handleUpdate: vi.fn(),
      toggleForm: vi.fn(),
      notifications: createNotifications(),
    });

    const button = container.querySelector('button[type="submit"]') as HTMLButtonElement;
    expect(button).toBeTruthy();
    expect(button.disabled).toBe(true);
  });

  it('should enable submit button when both passwords are provided', async () => {
    const { container } = render(UpdatePasswordForm, {
      handleUpdate: vi.fn(),
      toggleForm: vi.fn(),
      notifications: createNotifications(),
    });

    const password1 = container.querySelector('input[name="yourPassword1"]') as HTMLInputElement;
    const password2 = container.querySelector('input[name="yourPassword2"]') as HTMLInputElement;
    const button = container.querySelector('button[type="submit"]') as HTMLButtonElement;

    await userEvent.fill(password1, 'longenough');
    await userEvent.fill(password2, 'longenough');

    expect(button.disabled).toBe(false);
  });

  it('should call toggleForm when cancel is clicked', async () => {
    const toggleForm = vi.fn();
    render(UpdatePasswordForm, {
      handleUpdate: vi.fn(),
      toggleForm,
      notifications: createNotifications(),
    });

    const cancelButton = page.getByRole('button', { name: /cancel/i });
    await cancelButton.click();

    expect(toggleForm).toHaveBeenCalledTimes(1);
  });

  it('should show error and not submit when passwords do not match', async () => {
    const handleUpdate = vi.fn();
    const notifications = createNotifications();
    const { container } = render(UpdatePasswordForm, {
      handleUpdate,
      toggleForm: vi.fn(),
      notifications,
    });

    const password1 = container.querySelector('input[name="yourPassword1"]') as HTMLInputElement;
    const password2 = container.querySelector('input[name="yourPassword2"]') as HTMLInputElement;
    const button = page.getByRole('button', { name: /update/i });

    await userEvent.fill(password1, 'longenough');
    await userEvent.fill(password2, 'different');
    await button.click();

    expect(notifications.danger).toHaveBeenCalledTimes(1);
    expect(handleUpdate).not.toHaveBeenCalled();
  });

  it('should show error and not submit when passwords are too short', async () => {
    const handleUpdate = vi.fn();
    const notifications = createNotifications();
    const { container } = render(UpdatePasswordForm, {
      handleUpdate,
      toggleForm: vi.fn(),
      notifications,
    });

    const password1 = container.querySelector('input[name="yourPassword1"]') as HTMLInputElement;
    const password2 = container.querySelector('input[name="yourPassword2"]') as HTMLInputElement;
    const button = page.getByRole('button', { name: /update/i });

    await userEvent.fill(password1, 'short');
    await userEvent.fill(password2, 'short');
    await button.click();

    expect(notifications.danger).toHaveBeenCalledTimes(1);
    expect(handleUpdate).not.toHaveBeenCalled();
  });

  it('should submit update when passwords are valid', async () => {
    const handleUpdate = vi.fn();
    const notifications = createNotifications();
    const { container } = render(UpdatePasswordForm, {
      handleUpdate,
      toggleForm: vi.fn(),
      notifications,
    });

    const password1 = container.querySelector('input[name="yourPassword1"]') as HTMLInputElement;
    const password2 = container.querySelector('input[name="yourPassword2"]') as HTMLInputElement;
    const button = page.getByRole('button', { name: /update/i });

    await userEvent.fill(password1, 'longenough');
    await userEvent.fill(password2, 'longenough');
    await button.click();

    expect(handleUpdate).toHaveBeenCalledWith('longenough', 'longenough');
    expect(handleUpdate).toHaveBeenCalledTimes(1);
    expect(notifications.danger).not.toHaveBeenCalled();
  });
});
