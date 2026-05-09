export const storyboardPointsMaxLength = 3;

export const hasStoryboardPoints = (points?: string | null) => (points ?? '').trim() !== '';

export const parseStoryboardPoints = (points?: string | null) => {
  const value = (points ?? '').trim();
  if (value === '') {
    return null;
  }

  if (value === '1/2') {
    return 0.5;
  }

  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : null;
};

export const sanitizeStoryboardPoints = (points?: string | null) =>
  (points ?? '').trim().slice(0, storyboardPointsMaxLength);
