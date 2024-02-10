import TestSvelteComponent from "./TestSvelteComponent.svelte";

document.addEventListener("DOMContentLoaded", () => {
  const testComponentTarget = document.getElementById("svelte-test-component");
  if (testComponentTarget) {
    new TestSvelteComponent({
      target: testComponentTarget,
    });
  }
});
