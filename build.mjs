import esbuild from "esbuild";
import sveltePlugin from "esbuild-svelte";

// Define an asynchronous function to handle the build process
async function build() {
  try {
    // Bundle Svelte components
    await esbuild.build({
      entryPoints: ["javascript/svelte/main.js"], // Entry point for your Svelte app
      mainFields: ["svelte", "browser", "module", "main"],
      conditions: ["svelte", "browser"],
      bundle: true,
      outfile: "static/svelte_bundle.js", // Output file for Svelte
      minify: true,
      sourcemap: true,
      format: "esm", // Output format - 'esm' for ECMAScript modules
      plugins: [
        sveltePlugin(), // Use the Svelte plugin
      ],
    });

    // Bundle vanilla JS or other assets as needed
    await esbuild.build({
      entryPoints: ["javascript/vanilla/main.js"], // Entry point for your vanilla JS
      bundle: true,
      outfile: "static/vanilla_bundle.js", // Output file for vanilla JS
      minify: true,
      sourcemap: true,
    });

    console.log("Build completed successfully");
  } catch (error) {
    console.error("Build failed:", error);
    process.exit(1);
  }
}

// Execute the build function
build();
