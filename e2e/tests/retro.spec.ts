import { expect, test } from '../fixtures/user-sessions'
import { RetroPage } from '../fixtures/retro-page'

test.describe('Retro page', () => {
    let retro = { id: '', name: 'e2e retro page tests' }
    let retroLeave = { id: '' }
    let retroCancelDelete = { id: '' }
    let retroDelete = { id: '' }

    test.beforeAll(async ({ registeredPage, verifiedPage, adminPage }) => {
        const commonRetro = {
            retroName: `${retro.name}`,
            maxVotes: 3,
            brainstormVisibility: 'visible',
            retroFacilitators: [`${adminPage.user.email}`],
            format: 'worked_improve_question',
        }
        retro = await registeredPage.createRetro({ ...commonRetro })
        retroLeave = await verifiedPage.createRetro({ ...commonRetro })
        retroCancelDelete = await registeredPage.createRetro({
            ...commonRetro,
        })
        retroDelete = await registeredPage.createRetro({ ...commonRetro })
    })

    test('unauthenticated user redirects to register', async ({ page }) => {
        const bp = new RetroPage(page)
        await bp.goto(retro.id)

        const title = bp.page.locator('h1')
        await expect(title).toHaveText('Register')
    })

    test('guest user successfully loads', async ({ guestPage }) => {
        const bp = new RetroPage(guestPage.page)
        await bp.goto(retro.id)

        await expect(bp.retroTitle).toHaveText(retro.name)
    })

    test('registered user successfully loads', async ({ registeredPage }) => {
        const bp = new RetroPage(registeredPage.page)
        await bp.goto(retro.id)

        await expect(bp.retroTitle).toHaveText(retro.name)
    })

    test('user can leave retro', async ({ registeredPage }) => {
        const bp = new RetroPage(registeredPage.page)
        await bp.goto(retroLeave.id)

        await bp.page.click('[data-testid="retro-leave"]')
        await expect(bp.page.locator('h1')).toHaveText('My Retros')
    })

    test('delete retro confirmation cancel does not delete retro', async ({
        registeredPage,
    }) => {
        const bp = new RetroPage(registeredPage.page)
        await bp.goto(retroCancelDelete.id)

        await bp.retroDeleteBtn.click()
        await bp.retroDeleteCancelBtn.click()

        await expect(bp.retroTitle).toHaveText(retro.name)
    })

    test('delete retro confirmation confirm deletes retro and redirects to retros page', async ({
        registeredPage,
    }) => {
        const bp = new RetroPage(registeredPage.page)
        await bp.goto(retroDelete.id)

        await bp.retroDeleteBtn.click()
        await bp.retroDeleteConfirmBtn.click()

        await expect(bp.page.locator('h1')).toHaveText('My Retros')
    })
})
