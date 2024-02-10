<!-- javascript/svelte/TestSvelteComponent.svelte -->

<script>
    let count = 0;

    function increment() {
        count += 1;
    }
</script>

<main>
    <h2>Svelte Component Test</h2>
    <p>This is a test Svelte component.</p>
    <p>Count: {count}</p>
    <button on:click={increment}>Increment</button>
</main>

<style>
    main {
        padding: 1em;
        background-color: #f0f0f0;
        border: 1px solid #ccc;
        margin: 1em 0;
    }
    button {
        padding: 0.5em 1em;
        font-size: 1rem;
    }
</style>
