import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import RetroFeedbackGroup from './RetroFeedbackGroup.svelte';

const group = { id: 'group-1', name: 'Deploys', voteCount: 4, items: [] };

const setup = (props: { allowCumulativeVoting?: boolean; userVotesOnThisGroup?: number; userVotesUsed?: number }) => {
  const handleVote = vi.fn();
  const handleVoteSubtract = vi.fn();

  render(RetroFeedbackGroup, {
    phase: 'vote',
    group,
    voteLimit: 3,
    handleVote,
    handleVoteSubtract,
    ...props,
  });

  return { handleVote, handleVoteSubtract };
};

const button = (name: string) => page.getByRole('button', { name, exact: true });

describe('RetroFeedbackGroup voting', () => {
  it('adds a vote with the single vote button', async () => {
    const { handleVote, handleVoteSubtract } = setup({});

    await userEvent.click(button('Vote for Deploys'));

    expect(handleVote).toHaveBeenCalledWith('group-1');
    expect(handleVoteSubtract).not.toHaveBeenCalled();
  });

  it('removes the vote with the same button once voted, even at the limit', async () => {
    const { handleVote, handleVoteSubtract } = setup({ userVotesOnThisGroup: 1, userVotesUsed: 3 });

    await userEvent.click(button('Remove your vote from Deploys'));

    expect(handleVoteSubtract).toHaveBeenCalledWith('group-1');
    expect(handleVote).not.toHaveBeenCalled();
  });

  it('disables voting for a new group when no votes remain', async () => {
    setup({ userVotesUsed: 3 });

    await expect.element(button('Cannot vote for Deploys, no votes remaining')).toBeDisabled();
  });

  it('adds and removes votes with the cumulative stepper', async () => {
    const { handleVote, handleVoteSubtract } = setup({
      allowCumulativeVoting: true,
      userVotesOnThisGroup: 2,
      userVotesUsed: 2,
    });

    await expect.element(page.getByTestId('user-group-votes')).toHaveTextContent('2');
    await userEvent.click(button('Add a vote to Deploys'));
    await userEvent.click(button('Remove a vote from Deploys'));

    expect(handleVote).toHaveBeenCalledWith('group-1');
    expect(handleVoteSubtract).toHaveBeenCalledWith('group-1');
  });

  it('disables adding when over the limit and removing when there are no votes', async () => {
    setup({ allowCumulativeVoting: true, userVotesUsed: 4 });

    await expect.element(button('Cannot add vote to Deploys, no votes remaining')).toBeDisabled();
    await expect.element(button('Cannot remove vote, you have no votes on Deploys')).toBeDisabled();
  });
});
