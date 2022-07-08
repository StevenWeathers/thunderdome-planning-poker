describe('The Login Page', () => {
  beforeEach(() => {
    cy.task('db:teardown:registeredUser')
    cy.task('db:seed:registeredUser').as('currentUser')
  })

  it('should navigate to my battles page and reflect name in header', function () {
    // destructuring assignment of the currentUser object
    const { name, email, password } = this.currentUser

    cy.visit('/login')

    cy.get('input[name=yourEmail]').type(email)

    // {enter} causes the form to submit
    cy.get('input[name=yourPassword]').type(`${password}{enter}`)

    // we should be redirected to /battles
    cy.location('pathname').should('equal', '/battles')

    // our user cookie should be present
    cy.getCookie('warrior').should('exist')

    // UI should reflect this user being logged in
    cy.getByTestId('userprofile-link').should('contain', name)
  })
})