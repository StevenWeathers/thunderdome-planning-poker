import { render, fireEvent } from '@testing-library/svelte'

// Using Fixture Component due to inability to set Slot on Component render API
import TestButton from '../fixtures/TestHollowButton.svelte'

test('Renders the button with Slot text', () => {
    const { getByText } = render(TestButton, {})

    expect(getByText('Test Button')).toBeInTheDocument()
})

test('calls onClick when button is clicked', async () => {
    const mockCallback = jest.fn()
    const { getByText } = render(TestButton, {
        props: { onClick: mockCallback },
    })
    const button = getByText('Test Button')

    // Using await when firing events is unique to the svelte testing library because
    // we have to wait for the next `tick` so that Svelte flushes all pending state changes.
    await fireEvent.click(button)

    expect(mockCallback.mock.calls.length).toBe(1)
})
