<script lang="ts">
  import { Search } from '@lucide/svelte'

  let {
    value = $bindable(''),
    placeholder = 'Search rules...',
  }: {
    value?: string
    placeholder?: string
  } = $props()
</script>

<div class="relative">
  <Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-base-content/60" />
  <input
    bind:value
    {placeholder}
    class="input input-bordered w-full pl-9"
  />
</div>
