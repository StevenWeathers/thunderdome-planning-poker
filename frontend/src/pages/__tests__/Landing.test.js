import '@testing-library/jest-dom'
import { render } from '@testing-library/svelte'

import Landing from '../Landing.svelte'

describe('Landing Page', () => {
    it('should render successfully', () => {
        render(Landing, {
            xfetch: () => {},
            eventTag: () => {},
        })
    })
})
