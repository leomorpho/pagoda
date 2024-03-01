<script>
  import { Badge, MultiSelect } from 'flowbite-svelte';

  let selected = [];
  let countries = [
    { value: 'us', name: 'United States', color: 'indigo' },
    { value: 'ca', name: 'Canada', color: 'green' },
    { value: 'fr', name: 'France', color: 'blue' },
    { value: 'jp', name: 'Japan', color: 'red' },
    { value: 'en', name: 'England', color: 'yellow' }
  ];
</script>

<MultiSelect items={countries} bind:value={selected} let:item let:clear>
  <Badge color={item.color} dismissable params={{ duration: 100 }} on:close={clear}>
    {item.name}
  </Badge>
</MultiSelect>