document.addEventListener("DOMContentLoaded", () => {
  const testElement = document.getElementById("js-test-component");
  if (testElement) {
    testElement.textContent = "Content changed by static vanillajs!";
  }
});
