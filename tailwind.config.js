module.exports = {
  content: [
    "./javascript/**/*.{js,svelte}",
    "./**/*.templ",
    "./node_modules/flowbite/**/*.js",
  ],
  // https://themes.ionevolve.com/
  daisyui: {
    themes: [
      {
        lightmode: {
          primary: "#FFD500",
          "primary-focus": "#FFDC2C",
          "primary-content": "#1b1c22",

          secondary: "#9c00f0",
          "secondary-focus": "#b429ff",
          "secondary-content": "#ffffff",

          accent: "#37cdbe",
          "accent-focus": "#2ba69a",
          "accent-content": "#ffffff",

          neutral: "#3b424e",
          "neutral-focus": "#2a2e37",
          "neutral-content": "#ffffff",

          "base-100": "#ffffff",
          "base-200": "#f9fafb",
          "base-300": "#ced3d9",
          "base-content": "#1e2734",

          info: "#1c92f2",
          success: "#009485",
          warning: "#ff9900",
          error: "#ff5724",

          "--rounded-box": "1rem",
          "--rounded-btn": ".5rem",
          "--rounded-badge": "1.9rem",

          "--animation-btn": ".25s",
          "--animation-input": ".2s",

          "--btn-text-case": "uppercase",
          "--navbar-padding": ".5rem",
          "--border-btn": "1px",
        },
      },
      {
        darkmode: {
          primary: "#FFD500",
          "primary-focus": "#FFDC2C",
          "primary-content": "#1b1c22",

          secondary: "#b9ffb3",
          "secondary-focus": "#8aff80",
          "secondary-content": "#1b1c22",

          accent: "#ffffb3",
          "accent-focus": "#ffff80",
          "accent-content": "#1b1c22",

          neutral: "#22212c",
          "neutral-focus": "#1b1c22",
          "neutral-content": "#d5ccff",

          "base-100": "#302f3d",
          "base-200": "#22212c",
          "base-300": "#1b1c22",
          "base-content": "#d5ccff",

          info: "#1c92f2",
          success: "#009485",
          warning: "#ff9900",
          error: "#ff5724",

          "--rounded-box": "1rem",
          "--rounded-btn": ".5rem",
          "--rounded-badge": "1.9rem",

          "--animation-btn": ".25s",
          "--animation-input": ".2s",

          "--btn-text-case": "uppercase",
          "--navbar-padding": ".5rem",
          "--border-btn": "1px",
        },
      },
    ],
  },
  plugins: [require("daisyui"), require("flowbite/plugin")],
};
