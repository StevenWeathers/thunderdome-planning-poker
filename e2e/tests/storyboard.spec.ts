import { expect, test } from '../fixtures/user-sessions'
import { StoryboardPage } from '../fixtures/storyboard-page'

test.describe('Storyboard page', () => {
    let storyboard = { id: '', name: 'e2e storyboard page tests' }
    let storyboardLeave = { id: '' }
    let storyboardCancelDelete = { id: '' }
    let storyboardDelete = { id: '' }

    test.beforeAll(async ({ registeredPage, verifiedPage, adminPage }) => {
        const commonStoryboard = {
            storyboardName: `${storyboard.name}`,
            storyboardFacilitators: [`${adminPage.user.email}`],
        }
        storyboard = await registeredPage.createStoryboard({
            ...commonStoryboard,
        })
        storyboardLeave = await verifiedPage.createStoryboard({
            ...commonStoryboard,
        })
        storyboardCancelDelete = await registeredPage.createStoryboard({
            ...commonStoryboard,
        })
        storyboardDelete = await registeredPage.createStoryboard({
            ...commonStoryboard,
        })
    })

    test('unauthenticated user redirects to register', async ({ page }) => {
        const bp = new StoryboardPage(page)
        await bp.goto(storyboard.id)

        const title = bp.page.locator('h1')
        await expect(title).toHaveText('Register')
    })

    test('guest user successfully loads', async ({ guestPage }) => {
        const bp = new StoryboardPage(guestPage.page)
        await bp.goto(storyboard.id)

        await expect(bp.storyboardTitle).toHaveText(storyboard.name)
    })

    test('registered user successfully loads', async ({ registeredPage }) => {
        const bp = new StoryboardPage(registeredPage.page)
        await bp.goto(storyboard.id)

        await expect(bp.storyboardTitle).toHaveText(storyboard.name)
    })

    test('user can leave storyboard', async ({ registeredPage }) => {
        const bp = new StoryboardPage(registeredPage.page)
        await bp.goto(storyboardLeave.id)

        await bp.page.click('[data-testid="storyboard-leave"]')
        await expect(bp.page.locator('h1')).toHaveText('My Storyboards')
    })

    test('delete storyboard confirmation cancel does not delete storyboard', async ({
        registeredPage,
    }) => {
        const bp = new StoryboardPage(registeredPage.page)
        await bp.goto(storyboardCancelDelete.id)

        await bp.storyboardDeleteBtn.click()
        await bp.storyboardDeleteCancelBtn.click()

        await expect(bp.storyboardTitle).toHaveText(storyboard.name)
    })

    test('delete storyboard confirmation confirm deletes storyboard and redirects to storyboards page', async ({
        registeredPage,
    }) => {
        const bp = new StoryboardPage(registeredPage.page)
        await bp.goto(storyboardDelete.id)

        await bp.storyboardDeleteBtn.click()
        await bp.storyboardDeleteConfirmBtn.click()

        await expect(bp.page.locator('h1')).toHaveText('My Storyboards')
    })
})
