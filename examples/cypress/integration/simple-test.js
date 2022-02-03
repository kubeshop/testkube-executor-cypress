describe('The Home Page', () => {
  it('successfully loads', () => {
    cy.visit('https://testkube.io');
    expect(Cypress.env('testparam')).to.equal('testvalue');
    cy.contains('Efficient testing of k8s applications');
  })
})