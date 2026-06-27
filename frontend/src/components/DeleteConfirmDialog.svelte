<script lang="ts">
  let {
    open = $bindable(false),
    ruleName = '',
    onConfirm = () => {},
  }: {
    open?: boolean
    ruleName?: string
    onConfirm?: () => void
  } = $props()

  let dialogEl: HTMLDialogElement | undefined = $state()

  $effect(() => {
    if (open && dialogEl && !dialogEl.open) {
      dialogEl.showModal()
    } else if (!open && dialogEl?.open) {
      dialogEl.close()
    }
  })

  function closeDialog(): void {
    open = false
  }
</script>

<dialog class="modal" bind:this={dialogEl} onclose={closeDialog}>
  <div class="modal-box">
    <h3 class="text-lg font-bold">Delete Rule</h3>
    <p class="py-2 text-sm text-base-content/60">
      Are you sure you want to delete this rule? This cannot be undone.
    </p>

    {#if ruleName}
      <div class="rounded-box bg-base-300 px-3 py-2 text-sm">
        Rule: <span class="font-medium">{ruleName}</span>
      </div>
    {/if}

    <div class="modal-action">
      <button class="btn btn-ghost" onclick={closeDialog}>Cancel</button>
      <button class="btn btn-error" onclick={onConfirm}>Delete</button>
    </div>
  </div>

  <form method="dialog" class="modal-backdrop">
    <button onclick={closeDialog}>close</button>
  </form>
</dialog>
