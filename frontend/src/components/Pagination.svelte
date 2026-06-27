<script lang="ts">
  import { ChevronLeft, ChevronRight } from '@lucide/svelte'

  let {
    page = $bindable(1),
    pageSize = $bindable(10),
    total = 0,
    filteredTotal = 0,
  }: {
    page?: number
    pageSize?: number
    total?: number
    filteredTotal?: number
  } = $props()

  const totalPages = $derived(Math.max(1, Math.ceil(filteredTotal / pageSize)))

  const startItem = $derived(total === 0 ? 0 : (page - 1) * pageSize + 1)
  const endItem = $derived(Math.min(page * pageSize, filteredTotal))

  const visiblePages = $derived.by(() => {
    const maxVisible = 5
    let start = Math.max(1, page - Math.floor(maxVisible / 2))
    let end = Math.min(totalPages, start + maxVisible - 1)
    if (end - start + 1 < maxVisible) {
      start = Math.max(1, end - maxVisible + 1)
    }
    const pages: (number | 'ellipsis')[] = []
    if (start > 1) {
      pages.push(1)
      if (start > 2) pages.push('ellipsis')
    }
    for (let i = start; i <= end; i++) {
      pages.push(i)
    }
    if (end < totalPages) {
      if (end < totalPages - 1) pages.push('ellipsis')
      pages.push(totalPages)
    }
    return pages
  })

  function goToPage(p: number): void {
    if (p >= 1 && p <= totalPages) {
      page = p
    }
  }

  const pageSizes = [10, 15, 25, 50]
</script>

{#if total > 0 || filteredTotal > 0}
  <div class="flex items-center justify-between gap-4 border-t border-base-300 px-4 py-3 text-sm">
    <span class="text-base-content/60">
      {filteredTotal > 0 ? `${startItem}-${endItem} of ${filteredTotal}` : ''}
      {#if filteredTotal !== total}
        ({total} total)
      {/if}
    </span>

    <div class="flex items-center gap-2">
      <select class="select select-bordered select-sm" bind:value={pageSize}>
        {#each pageSizes as size}
          <option value={size}>{size}</option>
        {/each}
      </select>

      <div class="join">
        <button
          class="join-item btn btn-outline btn-sm"
          disabled={page <= 1}
          onclick={() => goToPage(page - 1)}
        >
          <ChevronLeft class="h-4 w-4" />
        </button>

        {#each visiblePages as p}
          {#if p === 'ellipsis'}
            <span class="join-item btn btn-ghost btn-sm pointer-events-none">...</span>
          {:else}
            <button
              class="join-item btn btn-sm {p === page ? 'btn-active' : 'btn-outline'}"
              onclick={() => goToPage(p)}
            >
              {p}
            </button>
          {/if}
        {/each}

        <button
          class="join-item btn btn-outline btn-sm"
          disabled={page >= totalPages}
          onclick={() => goToPage(page + 1)}
        >
          <ChevronRight class="h-4 w-4" />
        </button>
      </div>
    </div>
  </div>
{/if}
