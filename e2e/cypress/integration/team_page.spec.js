describe('The Team Page', () => {
  describe('Unauthenicated User', () => {
    it('redirects to login for unauthenticated user', function () {
      cy.visit('/team/bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa')

      cy.url().should('include', '/login')
    })
  })

  describe('Guest User', () => {})

  describe('Registered User', () => {
    beforeEach(() => {
      cy.task('db:teardown:registeredUser')
      cy.task('db:seed:registeredUser').as('currentUser')
    })

    it('successfully loads for authenticated registered user', function () {
      cy.login(this.currentUser)

      cy.createUserTeam(this.currentUser).then(() => {
        cy.visit(`/team/${this.currentTeam.id}`)

        cy.get('h1').should('contain', `Test Team`)
      })
    })
  })
})