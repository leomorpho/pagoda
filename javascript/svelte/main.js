import MultiSelectComponent from "./MultiSelectComponent.svelte";
import SvelteTodoComponent from "./SvelteTodoComponent.svelte";

// Define a registry object that maps names to Svelte component classes
const SvelteComponentRegistry = {
  SvelteTodoComponent,
  MultiSelectComponent,
};

// Assuming `window.svelteInstances` is a map to track component instances
window.svelteInstances = window.svelteInstances || {};

// Utility function to render any Svelte component by its registry name and ID
function renderSvelteComponentByName(componentName, id, props = {}) {
  const Component = SvelteComponentRegistry[componentName];
  if (!Component) {
    throw new Error(`Component ${componentName} not found in registry`);
  }

  const rootElement = document.getElementById(id);
  if (!rootElement) {
    throw new Error(`Could not find element with id ${id}`);
  }

  // Check for an existing instance and destroy it if present
  if (window.svelteInstances[id]) {
    window.svelteInstances[id].$destroy();
  }

  // Instantiate the new Svelte component with the target and props
  window.svelteInstances[id] = new Component({
    target: rootElement,
    props: props,
  });
}

window.renderSvelteComponent = function (componentName, id, props = {}) {
  renderSvelteComponentByName(componentName, id, props);
};
