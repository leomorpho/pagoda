import { initializeQuiz } from "./test_quiz";

window.initializeJS = function initializeApp(targetElement) {
  // Control zoom based on screen size
  controlZoom();
  initializeDarkModeSwitchers();
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

// Code taken from flowbite to switch between dark and light modes: https://flowbite.com/docs/customize/dark-mode/
function initializeDarkModeSwitcher() {
  var themeToggleDarkIcon = document.getElementById("theme-toggle-dark-icon");
  var themeToggleLightIcon = document.getElementById("theme-toggle-light-icon");

  // Change the icons inside the button based on previous settings
  if (
    localStorage.getItem("color-theme") === "dark" ||
    (!("color-theme" in localStorage) &&
      window.matchMedia("(prefers-color-scheme: dark)").matches)
  ) {
    themeToggleLightIcon.classList.remove("hidden");
  } else {
    themeToggleDarkIcon.classList.remove("hidden");
  }

  var themeToggleBtn = document.getElementById("theme-toggle");

  themeToggleBtn.addEventListener("click", function () {
    // toggle icons inside button
    themeToggleDarkIcon.classList.toggle("hidden");
    themeToggleLightIcon.classList.toggle("hidden");

    // if set via local storage previously
    if (localStorage.getItem("color-theme")) {
      if (localStorage.getItem("color-theme") === "light") {
        document.documentElement.classList.add("dark");
        localStorage.setItem("color-theme", "dark");
      } else {
        document.documentElement.classList.remove("dark");
        localStorage.setItem("color-theme", "light");
      }

      // if NOT set via local storage previously
    } else {
      if (document.documentElement.classList.contains("dark")) {
        document.documentElement.classList.remove("dark");
        localStorage.setItem("color-theme", "light");
      } else {
        document.documentElement.classList.add("dark");
        localStorage.setItem("color-theme", "dark");
      }
    }
  });
}

function initializeDarkModeSwitchers() {
  const themeToggleBtns = document.querySelectorAll(".theme-toggle");
  const themeToggleDarkIcons = document.querySelectorAll(
    ".theme-toggle-dark-icon"
  );
  const themeToggleLightIcons = document.querySelectorAll(
    ".theme-toggle-light-icon"
  );

  function updateIcons() {
    themeToggleDarkIcons.forEach((icon) => {
      if (
        localStorage.getItem("color-theme") === "dark" ||
        (!localStorage.getItem("color-theme") &&
          window.matchMedia("(prefers-color-scheme: dark)").matches)
      ) {
        icon.classList.add("hidden");
      } else {
        icon.classList.remove("hidden");
      }
    });
    themeToggleLightIcons.forEach((icon) => {
      if (
        localStorage.getItem("color-theme") === "dark" ||
        (!localStorage.getItem("color-theme") &&
          window.matchMedia("(prefers-color-scheme: dark)").matches)
      ) {
        icon.classList.remove("hidden");
      } else {
        icon.classList.add("hidden");
      }
    });
  }

  // Initial icon update based on the current theme
  updateIcons();

  themeToggleBtns.forEach((btn) => {
    btn.addEventListener("click", () => {
      // Toggle theme and update icons
      if (
        localStorage.getItem("color-theme") === "dark" ||
        (!localStorage.getItem("color-theme") &&
          window.matchMedia("(prefers-color-scheme: dark)").matches)
      ) {
        document.documentElement.classList.remove("dark");
        localStorage.setItem("color-theme", "light");
      } else {
        document.documentElement.classList.add("dark");
        localStorage.setItem("color-theme", "dark");
      }
      updateIcons();
    });
  });
}
