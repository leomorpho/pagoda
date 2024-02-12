/** @type {import('tailwindcss').Config} */

module.exports = {
  content: [
    "./javascript/**/*.{js,svelte}",
    "./**/*.templ",
    "./node_modules/flowbite/**/*.js",
  ],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui"), require("flowbite/plugin")],
};
