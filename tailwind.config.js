module.exports = {
  content: [
    "./javascript/**/*.{js,svelte}",
    "./**/*.templ",
    "./node_modules/flowbite/**/*.js",
  ],
  daisyui: {
    themes: ["light", "dark"],
  },
  plugins: [require("daisyui"), require("flowbite/plugin")],
};
