module.exports = {
  content: [
    "./javascript/**/*.{js,svelte}",
    "./**/*.templ",
    "./node_modules/flowbite/**/*.js",
  ],
  daisyui: {
    themes: [
      {
        lightmode: {
          // Change to any existing daisyui theme or make your own
          ...require("daisyui/src/theming/themes")["cmyk"],
          // Edit styles if required
        },
      },
      {
        darkmode: {
          // Change to any existing daisyui theme or make your own
          ...require("daisyui/src/theming/themes")["business"],
          // Edit styles if required
        },
      },
    ],
  },
  plugins: [require("daisyui"), require("flowbite/plugin")],
};
