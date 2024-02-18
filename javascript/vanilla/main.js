import { initializeQuiz } from "./test_quiz";

window.initializeJS = function initializeApp(targetElement) {
  console.log("initialize js");
  // Control zoom based on screen size
  controlZoom();

  window.darkModeSwitchersInitialized =
    window.darkModeSwitchersInitialized || false;
  initializeDarkModeSwitchers();

  const container = targetElement || document;

  // Check if the quiz container exists before trying to access its attributes
  const quizContainer = container.querySelector("#js-quiz-container");
  if (quizContainer && !quizContainer.hasAttribute("data-initialized")) {
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

// Code derived from flowbite to switch between dark and light modes: https://flowbite.com/docs/customize/dark-mode/
export function initializeDarkModeSwitchers() {
  if (window.darkModeSwitchersInitialized) {
    return;
  }

  const themeToggleBtns = document.querySelectorAll(".theme-toggle");
  const themeToggleDarkIcons = document.querySelectorAll(
    ".theme-toggle-dark-icon"
  );
  const themeToggleLightIcons = document.querySelectorAll(
    ".theme-toggle-light-icon"
  );

  // Define a function to update the icons based on the theme
  function updateIcons() {
    themeToggleDarkIcons.forEach((icon) => {
      icon.classList.toggle(
        "hidden",
        !(document.documentElement.getAttribute("data-theme") === "dark")
      );
    });
    themeToggleLightIcons.forEach((icon) => {
      icon.classList.toggle(
        "hidden",
        document.documentElement.getAttribute("data-theme") === "dark"
      );
    });
  }

  // Define a function to set the initial theme based on local storage or browser setting
  function setInitialTheme() {
    const storedTheme = localStorage.getItem("color-theme");
    const prefersDarkMode = window.matchMedia(
      "(prefers-color-scheme: dark)"
    ).matches;

    if (storedTheme) {
      document.documentElement.setAttribute("data-theme", storedTheme);
      document.documentElement.classList.toggle("dark", storedTheme === "dark");
    } else if (prefersDarkMode) {
      document.documentElement.setAttribute("data-theme", "dark");
      document.documentElement.classList.add("dark");
    } else {
      document.documentElement.setAttribute("data-theme", "light");
      document.documentElement.classList.remove("dark");
    }
    updateIcons();
  }

  // Set the initial theme
  setInitialTheme();

  // Set up the event listeners for the theme toggle buttons
  themeToggleBtns.forEach((btn) => {
    btn.addEventListener("click", () => {
      const currentTheme = document.documentElement.getAttribute("data-theme");
      const newTheme = currentTheme === "dark" ? "light" : "dark";
      document.documentElement.setAttribute("data-theme", newTheme);
      document.documentElement.classList.toggle("dark", newTheme === "dark");
      localStorage.setItem("color-theme", newTheme);
      updateIcons();
    });
  });

  // Mark the dark mode switchers as initialized to prevent duplicate initializations
  window.darkModeSwitchersInitialized = true;
}
