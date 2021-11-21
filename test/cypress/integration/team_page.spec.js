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
      // seed a user in the DB that we can control from our tests
      cy.createUser()
    })

    it('successfully loads for authenticated registered user', function () {
      cy.login(this.currentUser)

      cy.createUserTeam(this.currentUser).then(() => {
        cy.visit(`/team/${this.currentTeam.id}`)

        cy.get('h1').should('contain', `Team: Test Team`)
      })

      // cleanup our user (for some reason can't access this context in after utility
      cy.logout(this.currentUser)
    })
  })
})