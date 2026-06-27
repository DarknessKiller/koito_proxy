<script lang="ts">
  import { onMount } from 'svelte'
  import { Shield, Plus, X } from '@lucide/svelte'
  import type { Rule, RuleFormData } from '$lib/types'
  import { fetchRules, createRule, updateRule, deleteRule } from '$lib/api'
  import { toastManager } from '$lib/toast.svelte'
  import { getStores } from '$lib/stores.svelte'
  import RuleTable from './components/RuleTable.svelte'
  import RuleFormDialog from './components/RuleFormDialog.svelte'
  import DeleteConfirmDialog from './components/DeleteConfirmDialog.svelte'
  import SearchInput from './components/SearchInput.svelte'
  import Pagination from './components/Pagination.svelte'
  import ToastContainer from './components/ToastContainer.svelte'

  const store = getStores()

  let formDialogOpen = $state(false)
  let deleteDialogOpen = $state(false)
  let editingRule = $state<Rule | null>(null)
  let deletingRuleId = $state<string | null>(null)
  let deletingRuleName = $derived.by(() => {
    if (!deletingRuleId) return ''
    const rule = store.rules.find(r => r.id === deletingRuleId)
    if (!rule) return ''
    return rule.match_track_name || rule.match_artist_name || rule.id
  })

  const filteredRules = $derived.by(() => {
    const query = store.searchQuery.trim().toLowerCase()
    if (!query) return store.rules
    return store.rules.filter(r => {
      const fields = [
        r.match_track_name,
        r.match_artist_name,
        r.match_release_name,
        r.match_mbid,
        r.replace_track_name,
        r.replace_artist_name,
        r.replace_release_name,
        String(r.priority),
      ]
      return fields.some(f => f && f.toLowerCase().includes(query))
    })
  })

  const totalPages = $derived(Math.max(1, Math.ceil(filteredRules.length / store.pageSize)))

  const pageRules = $derived(
    filteredRules.slice(
      (store.page - 1) * store.pageSize,
      store.page * store.pageSize
    )
  )

  async function loadRules(): Promise<void> {
    store.loading = true
    store.error = null
    try {
      store.rules = await fetchRules()
    } catch (err) {
      const msg = err instanceof Error ? err.message : 'Failed to load rules'
      store.error = msg
      toastManager.show(msg, 'error')
    } finally {
      store.loading = false
    }
  }

  function openCreateDialog(): void {
    editingRule = null
    formDialogOpen = true
  }

  function openEditDialog(id: string): void {
    const rule = store.rules.find(r => r.id === id)
    if (rule) {
      editingRule = rule
      formDialogOpen = true
    }
  }

  function openDeleteDialog(id: string): void {
    deletingRuleId = id
    deleteDialogOpen = true
  }

  function buildRequestData(data: RuleFormData): Record<string, unknown> {
    const payload: Record<string, unknown> = {}

    if (data.match_track_name) payload.match_track_name = data.match_track_name
    if (data.match_artist_name) payload.match_artist_name = data.match_artist_name
    if (data.match_release_name) payload.match_release_name = data.match_release_name
    if (data.match_artist_names) {
      try { payload.match_artist_names = JSON.parse(data.match_artist_names) } catch { /* ignore */ }
    }
    if (data.match_duration_bucket) payload.match_duration_bucket = parseInt(data.match_duration_bucket, 10)
    if (data.match_mbid) payload.match_mbid = data.match_mbid
    if (data.replace_track_name) payload.replace_track_name = data.replace_track_name
    if (data.replace_artist_name) payload.replace_artist_name = data.replace_artist_name
    if (data.replace_release_name) payload.replace_release_name = data.replace_release_name
    if (data.replace_artist_names) {
      try { payload.replace_artist_names = JSON.parse(data.replace_artist_names) } catch { /* ignore */ }
    }
    payload.enabled = data.enabled

    return payload
  }

  function countFields(...fields: string[]): number {
    return fields.filter(f => f.trim().length > 0).length
  }

  async function handleSave(data: RuleFormData): Promise<void> {
    const matchCount = countFields(
      data.match_track_name,
      data.match_artist_name,
      data.match_release_name,
      data.match_artist_names,
      data.match_duration_bucket,
      data.match_mbid,
    )
    const replaceCount = countFields(
      data.replace_track_name,
      data.replace_artist_name,
      data.replace_release_name,
      data.replace_artist_names,
    )

    if (matchCount === 0 && replaceCount === 0) {
      toastManager.show('At least one match or replace field is required', 'error')
      return
    }

    if (replaceCount > 0 && (data.match_track_name.trim() || data.match_release_name.trim()) && matchCount < 2) {
      toastManager.show('At least two match criteria are required when using Track Name or Release Name match', 'error')
      return
    }

    try {
      const payload = buildRequestData(data)
      if (editingRule) {
        await updateRule(editingRule.id, payload)
        toastManager.show('Rule updated', 'success')
      } else {
        await createRule(payload)
        toastManager.show('Rule created', 'success')
      }
      formDialogOpen = false
      editingRule = null
      await loadRules()
    } catch (err) {
      const msg = err instanceof Error ? err.message : 'Failed to save rule'
      const friendly: Record<string, string> = {
        'no_fields': 'At least one match or replace field is required',
        'at_least_two_match_criteria': 'At least two match criteria are required when using Track Name or Release Name match',
        'at_least_one_match_criteria': 'At least one match field is required when replacing',
      }
      toastManager.show(friendly[msg] || msg, 'error')
    }
  }

  async function handleDelete(): Promise<void> {
    if (!deletingRuleId) return
    try {
      await deleteRule(deletingRuleId)
      toastManager.show('Rule deleted', 'success')
      deleteDialogOpen = false
      deletingRuleId = null
      store.page = 1
      await loadRules()
    } catch (err) {
      const msg = err instanceof Error ? err.message : 'Failed to delete rule'
      toastManager.show(msg, 'error')
    }
  }

  function closeAdmin(): void {
    const root = document.getElementById('koito-admin-root')
    if (root) root.style.display = 'none'
  }

  function onBackdropClick(e: MouseEvent): void {
    const target = e.target as HTMLElement
    if (target.classList.contains('koito-admin-backdrop')) {
      closeAdmin()
    }
  }

  function onBackdropKeydown(e: KeyboardEvent): void {
    if (e.key === 'Escape') {
      closeAdmin()
    }
  }

  onMount(() => {
    loadRules()
  })
</script>

<div data-theme="dark">
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
  role="dialog"
  tabindex="-1"
  class="koito-admin-backdrop fixed inset-0 z-[99999] flex items-center justify-center bg-black/90"
  onclick={onBackdropClick}
  onkeydown={onBackdropKeydown}
>
  <div
    id="koito-admin-content"
    class="flex flex-col overflow-hidden rounded-box border border-base-300 bg-base-200 text-base-content shadow-2xl"
    style="width: min(96vw, 1100px); height: min(92vh, 800px);"
  >
    <header class="flex items-center gap-3 border-b border-base-300 px-5 py-3.5">
      <Shield class="h-5 w-5 text-primary" />
      <h1 class="flex-1 text-lg font-semibold tracking-tight">Koito Proxy</h1>
      <button class="btn btn-ghost btn-square btn-sm" onclick={closeAdmin}>
        <X class="h-5 w-5" />
        <span class="sr-only">Close</span>
      </button>
    </header>

    <div class="flex flex-wrap items-center gap-3 border-b border-base-300 px-5 py-3">
      <span class="text-sm text-base-content/60">
        {store.rules.length} rule{store.rules.length !== 1 ? 's' : ''}
        {#if store.searchQuery}
          ({filteredRules.length} shown)
        {/if}
      </span>
      <div class="w-full sm:w-auto sm:max-w-xs">
        <SearchInput bind:value={store.searchQuery} />
      </div>
      <button class="btn btn-primary ml-auto" onclick={openCreateDialog}>
        <Plus class="h-4 w-4" /> New Rule
      </button>
    </div>

    <div class="min-h-0 flex-1 overflow-y-auto">
      {#if store.loading}
        <div class="flex items-center justify-center py-20 text-base-content/60">
          Loading rules...
        </div>
      {:else if store.error && store.rules.length === 0}
        <div class="flex items-center justify-center py-20 text-error">
          {store.error}
        </div>
      {:else if store.rules.length === 0}
        <div class="flex items-center justify-center py-20 text-base-content/60">
          No rules defined yet.
        </div>
      {:else}
        <RuleTable
          rules={pageRules}
          onEdit={openEditDialog}
          onDelete={openDeleteDialog}
        />
      {/if}
    </div>

    <Pagination
      bind:page={store.page}
      bind:pageSize={store.pageSize}
      total={store.rules.length}
      filteredTotal={filteredRules.length}
    />
  </div>
</div>

<RuleFormDialog
  bind:open={formDialogOpen}
  rule={editingRule}
  onSave={handleSave}
/>

<DeleteConfirmDialog
  bind:open={deleteDialogOpen}
  ruleName={deletingRuleName}
  onConfirm={handleDelete}
/>

<ToastContainer />
</div>
