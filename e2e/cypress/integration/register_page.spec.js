describe('The Register Page', () => {
  describe('Guest User', () => {
    it('should allow guest user signup', function () {
      cy.visit('/register')

      cy.get('h1').should('contain', 'Enlist to Battle')

      cy.get('input[name=yourName1]').type('TestGuestUser{enter}')

      cy.location('pathname').should('equal', '/battles')
    })
  })

  describe('Registered User', () => {
    beforeEach(() => {
      cy.task('db:teardown:registeredUser')
    })

    it('should allow user registration', function () {
      cy.visit('/register')

      cy.get('h1').should('contain', 'Enlist to Battle')

      cy.get('input[name=yourName2]').type('Registered Test User')

      cy.get('input[name=yourEmail]').type('registered@thunderdome.dev')

      cy.get('input[name=yourPassword1]').type('testreguserpassword')
      cy.get('input[name=yourPassword2]').type('testreguserpassword{enter}')

      cy.location('pathname').should('equal', '/battles')
    })

    it('should allow user registration from guest session', function () {
      cy.task('db:teardown:guestUser')
      cy.createGuestUser()

      cy.visit('/register')

      cy.get('h1').should('contain', 'Enlist to Battle')

      cy.get('input[name=yourEmail]').type('registered@thunderdome.dev')

      cy.get('input[name=yourPassword1]').type('testreguserpassword')
      cy.get('input[name=yourPassword2]').type('testreguserpassword{enter}')

      cy.location('pathname').should('equal', '/battles')
    })
  })
})