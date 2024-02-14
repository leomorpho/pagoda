window.initializeJS = function initializeApp(targetElement) {
  const container = targetElement || document;

  // Call each initializer with the container
  initializeTextChanges(container);

  // As you add more functionalities, add their initializers here
};

// Initializes custom text changes
function initializeTextChanges(container) {
  const testElement = container.querySelector("#js-test-component");
  if (testElement) {
    testElement.textContent = "Content changed by static vanillajs!";
  }
}
