import MultiSelectComponent from "./MultiSelectComponent.svelte";
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

// Utility function to render any Svelte component by ID
export function renderSvelteComponent(Component, id, props = {}) {
  const rootElement = document.getElementById(id);
  if (!rootElement) {
    throw new Error(`Could not find element with id ${id}`);
  }

  // If an instance already exists, destroy it first to clean up
  if (window.svelteInstances[id]) {
    window.svelteInstances[id].$destroy();
  }

  // Instantiate the new Svelte component with the target and props
  window.svelteInstances[id] = new Component({
    target: rootElement,
    props: props,
  });
}

window.renderMultiSelect = function (id) {
  renderSvelteComponent(MultiSelectComponent, id);
};
