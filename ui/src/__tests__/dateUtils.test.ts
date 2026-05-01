import { describe, expect, it } from 'vitest';

import { formatDayForInput } from '../dateUtils';

describe('dateUtils', () => {
  describe('formatDayForInput', () => {
    it('formats the day for the selected timezone', () => {
      const instant = new Date('2026-05-01T01:30:00.000Z');

      expect(formatDayForInput(instant, 'UTC')).toBe('2026-05-01');
      expect(formatDayForInput(instant, 'America/Los_Angeles')).toBe('2026-04-30');
    });
  });
});
