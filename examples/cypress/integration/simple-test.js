describe('The Home Page', () => {
  it('successfully loads', () => {
    cy.visit('https://testkube.kubeshop.io');

    expect(Cypress.env('testparam')).to.equal('testvalue');

    cy.contains(
      'Testkube provides a Kubernetes-native framework for test definition, execution and results'
    );
  });
});
