import { getContext, setContext } from 'svelte';
import type { RetroItem } from '../../types/retro';

// keep in step with RetroTagMaxLength and RetroItemMaxTags on the server
export const RETRO_TAG_MAX_LENGTH = 64;
export const RETRO_ITEM_MAX_TAGS = 10;

const contextKey = Symbol('retroTagSuggestions');

/**
 * Preset tags followed by any custom tags already used on the retro's items, without case-insensitive duplicates
 */
export function buildTagSuggestions(presets: string[] = [], items: RetroItem[] = []): string[] {
  const suggestions = new Map<string, string>();

  for (const tag of [...presets, ...items.flatMap(item => item.tags ?? [])]) {
    const key = tag.toLowerCase();
    if (!suggestions.has(key)) {
      suggestions.set(key, tag);
    }
  }

  return [...suggestions.values()].sort((a, b) => a.localeCompare(b));
}

export function setTagSuggestions(getSuggestions: () => string[]) {
  setContext(contextKey, getSuggestions);
}

export function getTagSuggestions(): () => string[] {
  return getContext<(() => string[]) | undefined>(contextKey) ?? (() => []);
}
