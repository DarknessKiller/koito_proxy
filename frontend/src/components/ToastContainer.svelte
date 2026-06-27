<script lang="ts">
  import { toastManager } from '$lib/toast.svelte'
</script>

{#each toastManager.items as toast (toast.id)}
  <div class="fixed bottom-4 right-4 z-[100001]">
    <div
      class="alert {toast.type === 'error' ? 'alert-error' : toast.type === 'success' ? 'alert-success' : 'alert-info'} shadow-lg"
    >
      <span>{toast.message}</span>
    </div>
  </div>
{/each}
