import { initializeQuiz } from "./test_quiz";

window.initializeJS = function initializeApp(targetElement) {
  console.log("initialize js");
  // Control zoom based on screen size
  controlZoom();

  window.darkModeSwitchersInitialized =
    window.darkModeSwitchersInitialized || false;
  initializeDarkModeSwitchers();
  adjustHoverColors();

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

const LIGHT = "lightmode";
const DARK = "darkmode";

// Code derived from flowbite to switch between dark and light modes: https://flowbite.com/docs/customize/dark-mode/
function initializeDarkModeSwitchers() {
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

  // Function to update icons based on the actual current theme
  function updateIcons() {
    const isDarkMode =
      document.documentElement.getAttribute("data-theme") === DARK ||
      document.documentElement.classList.contains(DARK);
    themeToggleDarkIcons.forEach((icon) => {
      icon.classList.toggle("hidden", isDarkMode); // Hide if dark mode
    });
    themeToggleLightIcons.forEach((icon) => {
      icon.classList.toggle("hidden", !isDarkMode); // Hide if not dark mode
    });
  }

  // Function to set the initial theme and update icons accordingly
  function setInitialTheme() {
    const storedTheme = localStorage.getItem("color-theme");
    const prefersDarkMode = window.matchMedia(
      "(prefers-color-scheme: dark)"
    ).matches;
    let initialTheme = LIGHT;

    if (storedTheme) {
      initialTheme = storedTheme;
    } else if (prefersDarkMode) {
      initialTheme = DARK;
    }

    document.documentElement.setAttribute("data-theme", initialTheme);
    document.documentElement.classList.toggle(DARK, initialTheme === DARK);
    updateIcons(); // Ensure icons are updated based on the initial theme
  }

  // Set the initial theme based on local storage or browser setting
  setInitialTheme();

  // Setup event listeners for theme toggle buttons
  themeToggleBtns.forEach((btn) => {
    btn.addEventListener("click", () => {
      const newTheme =
        document.documentElement.getAttribute("data-theme") === DARK
          ? LIGHT
          : DARK;
      document.documentElement.setAttribute("data-theme", newTheme);
      document.documentElement.classList.toggle(DARK, newTheme === DARK);
      localStorage.setItem("color-theme", newTheme);
      updateIcons(); // Update icons every time the theme is toggled
      adjustHoverColors();
    });
  });

  // Mark the dark mode switchers as initialized
  window.darkModeSwitchersInitialized = true;
}

function adjustHoverColors() {
  const isDarkMode =
    document.documentElement.getAttribute("data-theme") === DARK; // Ensure this matches how you set the theme
  console.log(`adjustHoverColors, isDark? ${isDarkMode}`);

  const rootStyle = document.documentElement.style;

  if (isDarkMode) {
    rootStyle.setProperty("--brightness-hover", "var(--brightness-hover-dark)");
  } else {
    rootStyle.setProperty(
      "--brightness-hover",
      "var(--brightness-hover-light)"
    );
  }
}
