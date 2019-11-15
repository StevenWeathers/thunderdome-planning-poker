import { configure, addParameters, addDecorator } from '@storybook/svelte'
import { withA11y } from '@storybook/addon-a11y'

// automatically import all files ending in *.stories.js
const req = require.context('../storybook/stories', true, /\.stories\.js$/)
function loadStories() {
    req.keys().forEach(filename => req(filename))
}

configure(loadStories, module)
addDecorator(withA11y)
addParameters({ viewport: { viewports: newViewports } })
