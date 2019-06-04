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
      await expect(page).toMatchElement('input[name="yourName"]')
    })

    it('should match an input with a "yourName" name then fill it with text', async () => {
      await expect(page).toFill('input[name="yourName"]', 'Thor')
    })

    it('should then submit the form and be redirected to Battles page', async () => {
      await expect(page).toClick('form[name="enlist"] button[type="submit"]')
      await page.waitFor(1000) // give the page a little time to redirect and render
      await expect(page).toMatchElement('h1', { text: 'My Battles' })
    })
  })

  describe('when on the My Battles page', () => {
    it('should match an input with a "battleName" name then fill it with text', async () => {
      await expect(page).toFill('input[name="battleName"]', 'Asgard')
    })

    it('should then submit the form and be redirected to that new Battle\'s page', async () => {
      await expect(page).toClick('form[name="createBattle"] button[type="submit"]')
      await page.waitFor(1000) // give the page a little time to redirect and render
      await expect(page).toMatchElement('h1', { text: '[Voting not started]' })
    })
  })
})
