export type EmojiPickerOption = {
  key: string;
  value: string;
  label: string;
};

export type EmojiPickerItem = EmojiPickerOption & {
  count: number;
  disabled?: boolean;
  pickerDisabled?: boolean;
  selected?: boolean;
};

export const DEVELOPER_REACTION_OPTIONS: EmojiPickerOption[] = [
  { key: 'rocket', value: '🚀', label: 'rocket' },
  { key: 'fire', value: '🔥', label: 'fire' },
  { key: 'thumbsdown', value: '👎', label: 'thumbs down' },
  { key: 'thumbsup', value: '👍', label: 'thumbs up' },
  { key: 'laugh', value: '😄', label: 'laugh' },
  { key: 'hooray', value: '🎉', label: 'hooray' },
  { key: 'confused', value: '😕', label: 'confused' },
  { key: 'heart', value: '❤️', label: 'heart' },
  { key: 'eyes', value: '👀', label: 'eyes' },
];
