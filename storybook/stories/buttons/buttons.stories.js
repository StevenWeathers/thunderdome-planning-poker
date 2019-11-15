import '../../css/utils.css'

import { storiesOf } from '@storybook/svelte'
import { action } from '@storybook/addon-actions'
import {
    select,
    boolean,
    number,
    text,
    withKnobs,
} from '@storybook/addon-knobs'

import ButtonHollow from '../../components/buttons/ButtonHollow.svelte'

const stories = storiesOf('Buttons | Buttons', module)
stories.addDecorator(withKnobs)

const availableColors = ['green', 'red', 'blue', 'teal', 'purple']

stories.add('Button Hollow', () => ({
    Component: ButtonHollow,
    props: {
        text: text('text', 'Button'),
        href: text('href', ''),
        disabled: boolean('disabled', false),
        color: select('color', availableColors, 'green'),
    },
}))
