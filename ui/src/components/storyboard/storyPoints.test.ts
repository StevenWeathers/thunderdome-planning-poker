import { describe, it, expect } from 'vitest';

import { sanitizeStoryboardPoints } from './storyPoints';

describe('sanitizeStoryboardPoints', () => {
  it.each([
    [' XXL ', 'XXL'],
    ['☕️', '☕️'],
    ['123456789', '12345678'],
    ['🚀🚀🚀🚀🚀🚀🚀🚀🚀', '🚀🚀🚀🚀🚀🚀🚀🚀'],
    [undefined, ''],
  ])('sanitizes %s', (points, expected) => {
    expect(sanitizeStoryboardPoints(points)).toBe(expected);
  });
});
