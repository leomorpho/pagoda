// build.js
const esbuild = require("esbuild");
const sveltePlugin = require("esbuild-svelte");

// Bundle Svelte components
esbuild
  .build({
    entryPoints: ["javascript/svelte/main.js"], // Entry point for your Svelte app
    bundle: true,
    outfile: "static/svelte_bundle.js", // Output file for Svelte
    minify: true,
    sourcemap: true,
    plugins: [sveltePlugin()], // Use the Svelte plugin
  })
  .catch(() => process.exit(1));

// Bundle vanilla JS
esbuild
  .build({
    entryPoints: ["javascript/vanilla/example.js"], // Entry point for your vanilla JS
    bundle: true,
    outfile: "static/vanilla_bundle.js", // Output file for vanilla JS
    minify: true,
    sourcemap: true,
  })
  .catch(() => process.exit(1));
