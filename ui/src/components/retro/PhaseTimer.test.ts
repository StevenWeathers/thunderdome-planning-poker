import { afterEach, describe, expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';

import PhaseTimer from './PhaseTimer.svelte';

describe('PhaseTimer component', () => {
  afterEach(() => {
    vi.useRealTimers();
  });

  it('should call onEnded when the countdown reaches zero', async () => {
    vi.useFakeTimers();

    const timeStart = new Date('2026-05-09T12:00:00.000Z');
    const onEnded = vi.fn();

    vi.setSystemTime(timeStart);

    render(PhaseTimer, {
      retroId: 'retro-1',
      timeLimitMin: 0,
      timeStart,
      onEnded,
    });

    expect(onEnded).toHaveBeenCalledTimes(1);
  });

  it('should update the countdown when timeLimitMin changes after mount', async () => {
    vi.useFakeTimers();

    const timeStart = new Date('2026-05-09T12:00:00.000Z');
    vi.setSystemTime(timeStart);

    const { container, rerender } = render(PhaseTimer, {
      retroId: 'retro-1',
      timeLimitMin: 10,
      timeStart,
    });

    expect(container.textContent?.replace(/\s+/g, '')).toContain('10m00s');

    await rerender({
      retroId: 'retro-1',
      timeLimitMin: 5,
      timeStart,
    });

    expect(container.textContent?.replace(/\s+/g, '')).toContain('05m00s');

    await vi.advanceTimersByTimeAsync(1000);

    expect(container.textContent?.replace(/\s+/g, '')).toContain('04m59s');
  });
});
