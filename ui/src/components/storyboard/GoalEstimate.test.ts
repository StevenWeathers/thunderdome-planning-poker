import { describe, it, expect } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import GoalEstimate from './GoalEstimate.svelte';
import type { StoryboardColumn } from '../../types/storyboard';

describe('GoalEstimate component', () => {
  it('should render successfully with empty columns', () => {
    const { container } = render(GoalEstimate, {
      columns: [],
    });

    const element = container.querySelector('span[title="Estimated Total Story Points"]');
    expect(element).toBeTruthy();
  });

  it('should display 0 points when no columns are provided', () => {
    const { container } = render(GoalEstimate, {
      columns: [],
    });

    const element = container.querySelector('span[title="Estimated Total Story Points"]');
    expect(element?.textContent).toContain('0');
    expect(element?.textContent).toContain('pts');
  });

  it('should calculate total points from a single column with one story', () => {
    const columns: StoryboardColumn[] = [
      {
        id: 'col-1',
        name: 'Column 1',
        sort_order: '1',
        personas: [],
        stories: [
          {
            id: 'story-1',
            name: 'Story 1',
            content: 'Content 1',
            color: '',
            points: 5,
            closed: false,
            link: '',
            sort_order: '1',
            annotations: [],
            comments: [],
          },
        ],
      },
    ];

    const { container } = render(GoalEstimate, {
      columns,
    });

    const element = container.querySelector('span[title="Estimated Total Story Points"]');
    expect(element?.textContent).toContain('5');
  });

  it('should calculate total points from multiple stories in a single column', () => {
    const columns: StoryboardColumn[] = [
      {
        id: 'col-1',
        name: 'Column 1',
        sort_order: '1',
        personas: [],
        stories: [
          {
            id: 'story-1',
            name: 'Story 1',
            content: 'Content 1',
            color: '',
            points: 3,
            closed: false,
            link: '',
            sort_order: '1',
            annotations: [],
            comments: [],
          },
          {
            id: 'story-2',
            name: 'Story 2',
            content: 'Content 2',
            color: '',
            points: 5,
            closed: false,
            link: '',
            sort_order: '2',
            annotations: [],
            comments: [],
          },
          {
            id: 'story-3',
            name: 'Story 3',
            content: 'Content 3',
            color: '',
            points: 8,
            closed: false,
            link: '',
            sort_order: '3',
            annotations: [],
            comments: [],
          },
        ],
      },
    ];

    const { container } = render(GoalEstimate, {
      columns,
    });

    const element = container.querySelector('span[title="Estimated Total Story Points"]');
    expect(element?.textContent).toContain('16'); // 3 + 5 + 8 = 16
  });

  it('should calculate total points from multiple columns with multiple stories', () => {
    const columns: StoryboardColumn[] = [
      {
        id: 'col-1',
        name: 'Column 1',
        sort_order: '1',
        personas: [],
        stories: [
          {
            id: 'story-1',
            name: 'Story 1',
            content: 'Content 1',
            color: '',
            points: 2,
            closed: false,
            link: '',
            sort_order: '1',
            annotations: [],
            comments: [],
          },
          {
            id: 'story-2',
            name: 'Story 2',
            content: 'Content 2',
            color: '',
            points: 3,
            closed: false,
            link: '',
            sort_order: '2',
            annotations: [],
            comments: [],
          },
        ],
      },
      {
        id: 'col-2',
        name: 'Column 2',
        sort_order: '2',
        personas: [],
        stories: [
          {
            id: 'story-3',
            name: 'Story 3',
            content: 'Content 3',
            color: '',
            points: 5,
            closed: false,
            link: '',
            sort_order: '1',
            annotations: [],
            comments: [],
          },
          {
            id: 'story-4',
            name: 'Story 4',
            content: 'Content 4',
            color: '',
            points: 8,
            closed: false,
            link: '',
            sort_order: '2',
            annotations: [],
            comments: [],
          },
        ],
      },
      {
        id: 'col-3',
        name: 'Column 3',
        sort_order: '3',
        personas: [],
        stories: [
          {
            id: 'story-5',
            name: 'Story 5',
            content: 'Content 5',
            color: '',
            points: 13,
            closed: false,
            link: '',
            sort_order: '1',
            annotations: [],
            comments: [],
          },
        ],
      },
    ];

    const { container } = render(GoalEstimate, {
      columns,
    });

    const element = container.querySelector('span[title="Estimated Total Story Points"]');
    expect(element?.textContent).toContain('31'); // 2 + 3 + 5 + 8 + 13 = 31
  });

  it('should handle columns with no stories', () => {
    const columns: StoryboardColumn[] = [
      {
        id: 'col-1',
        name: 'Column 1',
        sort_order: '1',
        personas: [],
        stories: [],
      },
      {
        id: 'col-2',
        name: 'Column 2',
        sort_order: '2',
        personas: [],
        stories: [],
      },
    ];

    const { container } = render(GoalEstimate, {
      columns,
    });

    const element = container.querySelector('span[title="Estimated Total Story Points"]');
    expect(element?.textContent).toContain('0');
  });

  it('should handle stories with zero points', () => {
    const columns: StoryboardColumn[] = [
      {
        id: 'col-1',
        name: 'Column 1',
        sort_order: '1',
        personas: [],
        stories: [
          {
            id: 'story-1',
            name: 'Story 1',
            content: 'Content 1',
            color: '',
            points: 0,
            closed: false,
            link: '',
            sort_order: '1',
            annotations: [],
            comments: [],
          },
          {
            id: 'story-2',
            name: 'Story 2',
            content: 'Content 2',
            color: '',
            points: 0,
            closed: false,
            link: '',
            sort_order: '2',
            annotations: [],
            comments: [],
          },
        ],
      },
    ];

    const { container } = render(GoalEstimate, {
      columns,
    });

    const element = container.querySelector('span[title="Estimated Total Story Points"]');
    expect(element?.textContent).toContain('0');
  });

  it('should handle mixed zero and non-zero points', () => {
    const columns: StoryboardColumn[] = [
      {
        id: 'col-1',
        name: 'Column 1',
        sort_order: '1',
        personas: [],
        stories: [
          {
            id: 'story-1',
            name: 'Story 1',
            content: 'Content 1',
            color: '',
            points: 0,
            closed: false,
            link: '',
            sort_order: '1',
            annotations: [],
            comments: [],
          },
          {
            id: 'story-2',
            name: 'Story 2',
            content: 'Content 2',
            color: '',
            points: 5,
            closed: false,
            link: '',
            sort_order: '2',
            annotations: [],
            comments: [],
          },
          {
            id: 'story-3',
            name: 'Story 3',
            content: 'Content 3',
            color: '',
            points: 0,
            closed: false,
            link: '',
            sort_order: '3',
            annotations: [],
            comments: [],
          },
          {
            id: 'story-4',
            name: 'Story 4',
            content: 'Content 4',
            color: '',
            points: 3,
            closed: false,
            link: '',
            sort_order: '4',
            annotations: [],
            comments: [],
          },
        ],
      },
    ];

    const { container } = render(GoalEstimate, {
      columns,
    });

    const element = container.querySelector('span[title="Estimated Total Story Points"]');
    expect(element?.textContent).toContain('8'); // 0 + 5 + 0 + 3 = 8
  });

  it('should display the "pts" label alongside the total', () => {
    const columns: StoryboardColumn[] = [
      {
        id: 'col-1',
        name: 'Column 1',
        sort_order: '1',
        personas: [],
        stories: [
          {
            id: 'story-1',
            name: 'Story 1',
            content: 'Content 1',
            color: '',
            points: 10,
            closed: false,
            link: '',
            sort_order: '1',
            annotations: [],
            comments: [],
          },
        ],
      },
    ];

    const { container } = render(GoalEstimate, {
      columns,
    });

    const element = container.querySelector('span[title="Estimated Total Story Points"]');
    expect(element?.textContent).toMatch(/10.*pts/);
  });
});
