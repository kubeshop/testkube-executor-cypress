describe('The Home Page', () => {
  it('successfully loads', () => {
    cy.visit('https://kubtest.io') 
    ey.contains('Efficient testing of k8s applications')

  })
})