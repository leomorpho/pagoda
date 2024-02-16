import { initializeQuiz } from "./test_quiz";

window.initializeJS = function initializeApp(targetElement) {
  // Control zoom based on screen size
  controlZoom();
  const container = targetElement || document;

  // Initialize the quiz only if it hasn't been initialized yet
  if (
    !container
      .querySelector("#js-quiz-container")
      .hasAttribute("data-initialized")
  ) {
    initializeQuiz(container);
  }
};

// controlZoom prevents zooming on mobile devices to improve user experience. Note that it is
// not enforced by all browsers, so is not guaranteed to work on all of them. TODO: we may want to
// remove that entirely as Brave/Chrome ignore these directives, and it may be globally useless if
// all browsers ignore it.
function controlZoom() {
  const viewportMeta = document.querySelector('meta[name="viewport"]');

  function updateViewport() {
    if (window.innerWidth < 1024) {
      // Disable zooming on small screens
      viewportMeta.setAttribute(
        "content",
        "width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no"
      );
    } else {
      // Allow zooming on large screens
      viewportMeta.setAttribute(
        "content",
        "width=device-width, initial-scale=1.0"
      );
    }
  }

  // Update on initial load
  updateViewport();

  // Update on window resize
  window.addEventListener("resize", updateViewport);
}
