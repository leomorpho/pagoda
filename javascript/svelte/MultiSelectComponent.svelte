<script>
  import MultiSelect from "svelte-multiselect";

  export let items = [`Svelte`, `React`, `Vue`, `Angular`, `...`];
  export let selected = [];
  export let placeholder = "Select options...";
  export let formInputName = "input_name";
  export let postURL = "/multiselect-post-request";
  export let csrfToken = "";

  // Ensure 'selected' is always an array
  $: selected = Array.isArray(selected) ? selected : [];

  // Function called when a new option is added
  async function handleAdd(event) {
    const newOption = event.detail.option;

    // Check if the option is already in the list to avoid duplicates
    if (!items.includes(newOption)) {
      // Optimistically add the new option to both items and selected
      items = [...items, newOption];
      selected = [...selected, newOption];

      try {
        // Attempt to create the category in the backend
        const response = await createCategory(newOption);
        if (!response.success) {
          // If the backend fails, remove the option and alert the user
          removeOption(newOption);
          alert("Failed to create new category. Please try again.");
        }
      } catch (error) {
        console.error("Error creating category:", error);
        removeOption(newOption);
        alert("Failed to create new category. Please try again.");
      }
    }
  }

  // Function to remove an option from items and selected
  function removeOption(option) {
    items = items.filter((item) => item !== option);
    selected = selected.filter((sel) => sel !== option);
  }

  // Function updated to post to the provided URL
  async function createCategory(categoryName) {
    try {
      const response = await fetch(postURL, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-CSRF-Token": csrfToken,
        },
        body: JSON.stringify({ category_name: categoryName }),
        credentials: "include",
      });

      if (!response.status == 201) {
        throw new Error("Network response was not ok");
      }

      const data = await response.json();
      return {
        success: response.status === 201, // Check if status code is 201 for success
        name: data.name, // Assuming the backend returns the name attribute in the response
      };
    } catch (error) {
      console.error("Error:", error);
      return { success: false }; // Assume failure on error
    }
  }
</script>

<MultiSelect
  bind:selected
  options={items}
  {placeholder}
  on:add={handleAdd}
  outerDivClass="!w-full !input !input-bordered !border !border-gray-300 !rounded-md !p-2 !bg-white !text-slate-500"
  liSelectedClass="!bg-orange-500 dark:!bg-blue-600 !p-2 !text-white"
  liOptionClass="!p-1 !m-1 !text-slate-600"
  liUserMsgClass="!p-1 !m-1 !text-slate-600"
  liActiveOptionClass="!bg-slate-200 dark:!bg-blue-300"
  allowUserOptions={true}
  createOptionMsg="Hit enter to create"
/>

<!-- Hidden inputs for form submission -->
{#each selected as option}
  <input type="hidden" name={formInputName} value={option} />
{/each}
