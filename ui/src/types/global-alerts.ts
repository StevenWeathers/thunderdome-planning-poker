export type GlobalAlert = {
  id: string;
  type: 'NEW' | 'ERROR' | 'INFO' | 'SUCCESS' | 'WARNING';
  content: string;
  allowDismiss?: boolean;
  registeredOnly?: boolean;
  active?: boolean;
  createdDate?: string;
  name?: string;
  updatedDate?: string;
};
