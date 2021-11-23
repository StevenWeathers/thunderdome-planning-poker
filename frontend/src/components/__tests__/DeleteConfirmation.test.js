import '@testing-library/jest-dom'
import { render, fireEvent } from '@testing-library/svelte'

import DeleteConfirmation from '../DeleteConfirmation.svelte'

describe('DeleteConfirmation component', () => {
    it('should render successfully', () => {
        render(DeleteConfirmation, {})
    })

    it('should match snapshot', () => {
        const { container } = render(DeleteConfirmation, {})

        expect(container).toMatchSnapshot()
    })

    it('should match snapshot when permanent=false', () => {
        const { container } = render(DeleteConfirmation, { permanent: false })

        expect(container).toMatchSnapshot()
    })

    it('should fire handleDelete when confirmed', async () => {
        const stub = jest.fn()
        const { getByText } = render(DeleteConfirmation, { handleDelete: stub })
        const button = getByText('Confirm Delete')

        await fireEvent.click(button)

        expect(stub).toHaveBeenCalled()
    })

    it('should not fire handleDelete when cancel and instead fire toggleDelete', async () => {
        const handleDelete = jest.fn()
        const toggleDelete = jest.fn()
        const { getByText } = render(DeleteConfirmation, {
            handleDelete,
            toggleDelete,
        })
        const button = getByText('Cancel')

        await fireEvent.click(button)

        expect(handleDelete).not.toHaveBeenCalled()
        expect(toggleDelete).toHaveBeenCalled()
    })
})
