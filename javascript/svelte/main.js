import TestSvelteComponent from "./TestSvelteComponent.svelte";

// Assuming `window.svelteInstances` is a map to track component instances
window.svelteInstances = window.svelteInstances || {};

window.initializeAppSvelte = function (targetElement) {
  const container = targetElement || document;

  // Initialize specific Svelte components
  initializeTestSvelteComponent(container);

  // Initialize other Svelte components in a similar manner
};

function initializeTestSvelteComponent(container) {
  const testComponentTarget = container.querySelector("#svelte-test-component");

  // If an instance already exists, destroy it
  if (window.svelteInstances.testComponent) {
    window.svelteInstances.testComponent.$destroy();
  }

  // Initialize the new instance
  if (testComponentTarget) {
    window.svelteInstances.testComponent = new TestSvelteComponent({
      target: testComponentTarget,
    });
  }
}
