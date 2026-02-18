export type NotificationService = {
  show: (message: string, timeout?: number, theme?: string) => void;
  success: (message: string, timeout?: number) => void;
  danger: (message: string, timeout?: number) => void;
  warning: (message: string, timeout?: number) => void;
  info: (message: string, timeout?: number) => void;
  removeToast: (id: number) => void;
};
