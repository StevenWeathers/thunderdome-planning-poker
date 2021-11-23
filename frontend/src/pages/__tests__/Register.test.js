import '@testing-library/jest-dom'
import { render } from '@testing-library/svelte'

import Register from '../Register.svelte'

describe('Register Page', () => {
    it('should render successfully', () => {
        render(Register, {
            xfetch: () => {},
            eventTag: () => {},
            notifications: () => {},
            router: () => {},
            battleId: null,
        })
    })
})
