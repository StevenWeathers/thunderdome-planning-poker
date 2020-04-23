const appUrl = process.env.APP_URL || 'http://localhost:8080/'

describe('Thunderdome App', () => {
    beforeAll(async () => {
        await page.goto(appUrl) // @TODO - make this configurable for the endpoint
    })

    describe('when on the Landing page', () => {
        it('should match a link with a "Create a Battle" text inside', async () => {
            await expect(page).toMatchElement('a', { text: 'Create a Battle' })
        })
    })

    describe('when clicking the Create a Battle button as a new visitor', () => {
        it('should redirect to Enlist page', async () => {
            await expect(page).toClick('a', { text: 'Create a Battle' })
            await expect(page).toMatchElement('input[name="yourName1"]')
        })

        it('should match an input with a "yourName" name then fill it with text', async () => {
            await expect(page).toFill('input[name="yourName1"]', 'Thor')
        })

        it('should then submit the form and be redirected to Battles page', async () => {
            await expect(page).toClick(
                'form[name="registerGuest"] button[type="submit"]',
            )
            await page.waitFor(1000) // give the page a little time to redirect and render
            await expect(page).toMatchElement('h1', { text: 'My Battles' })
        })
    })

    describe('when on the My Battles page', () => {
        it('should match an input with a "battleName" name then fill it with text', async () => {
            await expect(page).toFill('input[name="battleName"]', 'Asgard')
        })

        it("should then submit the form and be redirected to that new Battle's page", async () => {
            await expect(page).toClick(
                'form[name="createBattle"] button[type="submit"]',
            )
            await page.waitFor(1000) // give the page a little time to redirect and render
            await expect(page).toMatchElement('h1', {
                text: '[Voting not started]',
            })
        })
    })

    describe('when on a Battle page as leader', () => {
        it('should have Add Plan Button', async () => {
            await expect(page).toMatchElement('button', { text: 'Add Plan' })
        })

        it('should have Delete Battle Button', async () => {
            await expect(page).toMatchElement('button', {
                text: 'Delete Battle',
            })
        })

        it('should have locked voting cards', async () => {
            await expect(page).toMatchElement(
                '[data-testId="pointCard"][data-locked="true"]',
            )
        })

        it('should have Unpointed and Pointed tabs count as 0', async () => {
            await expect(page).toMatchElement('button', {
                text: 'Unpointed (0)',
            })
            await expect(page).toMatchElement('button', { text: 'Pointed (0)' })
        })

        describe('when clicking the Add Plan button', () => {
            it('should display the addPlan form modal', async () => {
                await expect(page).toClick('button', { text: 'Add Plan' })
                await expect(page).toMatchElement('input[name="planName"]')
            })

            it('should then match an input with an id of planName then fill it with text', async () => {
                await expect(page).toFill('input[name="planName"]', 'Drink Ale')
            })

            it('should then submit the form', async () => {
                await expect(page).toClick(
                    'form[name="addPlan"] button[type="submit"]',
                )
            })

            it('should display the new Plan', async () => {
                await expect(page).toMatchElement(
                    '[data-testId="battlePlanName"]',
                    {
                        text: 'Drink Ale',
                    },
                )
            })

            it('should display control buttons', async () => {
                await expect(
                    page,
                ).toMatchElement(
                    '[data-testId="battlePlan"][data-planName="Drink Ale"] button',
                    { text: 'Delete' },
                )
                await expect(
                    page,
                ).toMatchElement(
                    '[data-testId="battlePlan"][data-planName="Drink Ale"] button',
                    { text: 'Edit' },
                )
                await expect(
                    page,
                ).toMatchElement(
                    '[data-testId="battlePlan"][data-planName="Drink Ale"] button',
                    { text: 'Activate' },
                )
            })

            it('should have Unpointed tab count as 1', async () => {
                await expect(page).toMatchElement('button', {
                    text: 'Unpointed (1)',
                })
            })
        })

        describe('when clicking the Edit Plan button', () => {
            it('should display the addPlan form modal', async () => {
                await expect(page).toClick('button', { text: 'Edit' })
                await expect(page).toMatchElement('input[name="planName"]')
            })

            it('should then match an input with an id of planName then fill it with text', async () => {
                await expect(page).toFill(
                    'input[name="planName"]',
                    'Drink Ale!',
                )
            })

            it('should then submit the form', async () => {
                await expect(page).toClick(
                    'form[name="addPlan"] button[type="submit"]',
                )
            })

            it('should display the updated Plan', async () => {
                await expect(page).toMatchElement(
                    '[data-testId="battlePlanName"]',
                    {
                        text: 'Drink Ale!',
                    },
                )
            })
        })

        describe('when clicking the Activate Plan button', () => {
            it('should enable the Point Cards', async () => {
                await expect(page).toClick('button', { text: 'Activate' })
                await expect(page).toMatchElement(
                    '[data-testId="pointCard"][data-locked="false"]',
                )
            })

            it('should enable the Voting Controls', async () => {
                await expect(page).toMatchElement(
                    '[data-testId="votingControls"]',
                )
            })

            it('should enable the Timer', async () => {
                await expect(page).toMatchElement('[data-testId="votingTimer"]')
            })
        })

        describe('when clicking a point card', () => {
            it('should activate that point card', async () => {
                await expect(page).toClick(
                    '[data-testId="pointCard"][data-point="1"]',
                )
                await expect(page).toMatchElement(
                    '[data-testId="pointCard"][data-point="1"][data-active="true"]',
                )
            })

            it('should flag the warrior as voted', async () => {
                await expect(page).toMatchElement(
                    '[data-testId="warriorCard"][data-warriorName="Thor"] [data-icon="vote-yea"]',
                )
            })
        })

        describe('when clicking the Finish Voting', () => {
            it('should enable the savePlanPoints form', async () => {
                await expect(page).toClick('button', { text: 'Finish Voting' })
                await expect(page).toMatchElement('form[name="savePlanPoints"]')
            })

            it('should display the point count on the card(s) and lock them', async () => {
                await expect(page).toMatchElement(
                    '[data-testId="pointCard"][data-point="1"][data-locked="true"]',
                )
                await expect(
                    page,
                ).toMatchElement(
                    '[data-testId="pointCard"][data-point="1"] [data-testId="pointCardCount"]',
                    { test: '1' },
                )
            })

            it('should display each warriors voted points', async () => {
                await expect(
                    page,
                ).toMatchElement(
                    '[data-testId="warriorCard"][data-warriorName="Thor"] [data-testId="warriorPoints"]',
                    { text: '1' },
                )
            })

            it('should display control buttons', async () => {
                await expect(
                    page,
                ).toMatchElement(
                    '[data-testId="battlePlan"][data-planName="Drink Ale!"] button',
                    { text: 'Delete' },
                )
                await expect(
                    page,
                ).toMatchElement(
                    '[data-testId="battlePlan"][data-planName="Drink Ale!"] button',
                    { text: 'Edit' },
                )
                await expect(
                    page,
                ).toMatchElement(
                    '[data-testId="battlePlan"][data-planName="Drink Ale!"] button',
                    { text: 'Activate' },
                )
            })
        })

        describe('when filling out the final points form', () => {
            it('should fill out the final points select', async () => {
                await expect(page).toSelect('select[name="planPoints"]', '1')
            })

            it('should submit the Final Points form', async () => {
                await expect(page).toClick(
                    'form[name="savePlanPoints"] button[type="submit"]',
                )
            })

            it('should update plan to Pointed Tab and its count', async () => {
                await expect(page).toMatchElement('button', {
                    text: 'Pointed (1)',
                })
            })
        })

        describe('when clicking the Delete Battle button', () => {
            it('should redirect to Battles page', async () => {
                await expect(page).toClick('button', { text: 'Delete Battle' })
                await page.waitFor(1000) // wait for page to redirect
                await expect(page).toMatchElement('h1', { text: 'My Battles' })
            })
        })
    })
})
