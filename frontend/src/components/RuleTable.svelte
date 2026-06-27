<script lang="ts">
  import type { Rule } from '$lib/types'
  import { Pencil, Trash2 } from '@lucide/svelte'

  let {
    rules = [] as Rule[],
    onEdit = (_id: string) => {},
    onDelete = (_id: string) => {},
  }: {
    rules?: Rule[]
    onEdit?: (id: string) => void
    onDelete?: (id: string) => void
  } = $props()

  function statusBadge(rule: Rule) {
    if (!rule.enabled) return 'badge-error'
    if (!rule.valid) return 'badge-ghost'
    return 'badge-success'
  }

  function statusLabel(rule: Rule): string {
    if (!rule.enabled) return 'off'
    if (!rule.valid) return 'low'
    return 'on'
  }

  function cellValue(value: string | null): string {
    return value ?? '-'
  }
</script>

<div class="hidden md:block">
  <table class="table">
    <thead>
      <tr class="text-xs uppercase text-base-content/60">
        <th colspan="3" class="text-center">Match</th>
        <th colspan="3" class="text-center">Replace</th>
        <th class="text-center">MBID</th>
        <th class="text-center">Pri</th>
        <th class="text-center">Status</th>
        <th class="text-center">Actions</th>
      </tr>
      <tr class="text-xs uppercase text-base-content/60">
        <th class="text-center font-normal">Track</th>
        <th class="text-center font-normal">Artist</th>
        <th class="text-center font-normal">Release</th>
        <th class="text-center font-normal">Track</th>
        <th class="text-center font-normal">Artist</th>
        <th class="text-center font-normal">Release</th>
        <th class="text-center"></th>
        <th class="text-center"></th>
        <th class="text-center"></th>
        <th class="text-center"></th>
      </tr>
    </thead>
    <tbody>
      {#each rules as rule (rule.id)}
        <tr class="hover:bg-base-300/50">
          <td class="text-center text-sm">{cellValue(rule.match_track_name)}</td>
          <td class="text-center text-sm">{cellValue(rule.match_artist_name)}</td>
          <td class="text-center text-sm">{cellValue(rule.match_release_name)}</td>
          <td class="text-center text-sm">{cellValue(rule.replace_track_name)}</td>
          <td class="text-center text-sm">{cellValue(rule.replace_artist_name)}</td>
          <td class="text-center text-sm">{cellValue(rule.replace_release_name)}</td>
          <td class="text-center font-mono text-xs text-base-content/60">
            {cellValue(rule.match_mbid)}
          </td>
          <td class="text-center font-mono text-xs text-base-content/60">
            {rule.priority}
          </td>
          <td class="text-center">
            <span class="badge {statusBadge(rule)}">{statusLabel(rule)}</span>
          </td>
          <td class="text-center">
            <div class="flex items-center justify-center gap-1">
              <button class="btn btn-ghost btn-square btn-sm" onclick={() => onEdit(rule.id)}>
                <Pencil class="h-4 w-4" />
                <span class="sr-only">Edit</span>
              </button>
              <button class="btn btn-ghost btn-square btn-sm text-error" onclick={() => onDelete(rule.id)}>
                <Trash2 class="h-4 w-4" />
                <span class="sr-only">Delete</span>
              </button>
            </div>
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>

<div class="md:hidden">
  <div class="flex flex-col gap-3 p-3">
    {#each rules as rule (rule.id)}
      <div class="rounded-box border border-base-300 bg-base-100 p-4">
        <div class="mb-2 flex items-start justify-between gap-2">
          <div class="text-sm font-semibold">
            {rule.match_track_name || rule.match_artist_name || rule.replace_track_name || 'Untitled'}
          </div>
          <span class="badge {statusBadge(rule)} shrink-0">{statusLabel(rule)}</span>
        </div>
        <div class="space-y-1 text-xs text-base-content/60">
          {#if rule.match_artist_name}
            <div><span class="text-[10px] font-semibold uppercase tracking-wider">Match Artist</span> {rule.match_artist_name}</div>
          {/if}
          {#if rule.match_release_name}
            <div><span class="text-[10px] font-semibold uppercase tracking-wider">Match Release</span> {rule.match_release_name}</div>
          {/if}
          {#if rule.replace_track_name}
            <div><span class="text-[10px] font-semibold uppercase tracking-wider">Replace Track</span> {rule.replace_track_name}</div>
          {/if}
          {#if rule.replace_artist_name}
            <div><span class="text-[10px] font-semibold uppercase tracking-wider">Replace Artist</span> {rule.replace_artist_name}</div>
          {/if}
          {#if rule.replace_release_name}
            <div><span class="text-[10px] font-semibold uppercase tracking-wider">Replace Release</span> {rule.replace_release_name}</div>
          {/if}
          {#if rule.match_mbid}
            <div><span class="text-[10px] font-semibold uppercase tracking-wider">MBID</span> <span class="font-mono">{rule.match_mbid}</span></div>
          {/if}
          <div><span class="text-[10px] font-semibold uppercase tracking-wider">Pri</span> {rule.priority}</div>
        </div>
        <div class="mt-3 flex gap-2 border-t border-base-300 pt-3">
          <button class="btn btn-outline btn-sm flex-1" onclick={() => onEdit(rule.id)}>
            <Pencil class="mr-1 h-3.5 w-3.5" /> Edit
          </button>
          <button class="btn btn-outline btn-sm flex-1 text-error" onclick={() => onDelete(rule.id)}>
            <Trash2 class="mr-1 h-3.5 w-3.5" /> Delete
          </button>
        </div>
      </div>
    {/each}
  </div>
</div>
