import '@testing-library/jest-dom'
import { render } from '@testing-library/svelte'

import Login from '../Login.svelte'

describe('Login Page', () => {
    it('should render successfully', () => {
        render(Login, {
            xfetch: () => {},
            eventTag: () => {},
            notifications: () => {},
            router: () => {},
            battleId: null,
        })
    })
})
