// build.mjs
import esbuild from "esbuild";
import sveltePlugin from "esbuild-svelte";

// Define an asynchronous function to handle the build process
async function build() {
  try {
    // Bundle Svelte components
    await esbuild.build({
      entryPoints: ["javascript/svelte/main.js"], // Entry point for your Svelte app
      bundle: true,
      outfile: "static/svelte_bundle.js", // Output file for Svelte
      minify: true,
      sourcemap: true,
      plugins: [sveltePlugin()], // Use the Svelte plugin
    });

    // Bundle vanilla JS
    await esbuild.build({
      entryPoints: ["javascript/vanilla/example.js"], // Entry point for your vanilla JS
      bundle: true,
      outfile: "static/vanilla_bundle.js", // Output file for vanilla JS
      minify: true,
      sourcemap: true,
    });
  } catch (error) {
    console.error(error);
    process.exit(1);
  }
}

// Execute the build function
build();
