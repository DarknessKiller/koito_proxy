<script lang="ts">
  import type { Rule, RuleFormData } from '$lib/types'

  let {
    open = $bindable(false),
    rule = null as Rule | null,
    onSave = (_data: RuleFormData) => {},
  }: {
    open?: boolean
    rule?: Rule | null
    onSave?: (data: RuleFormData) => void
  } = $props()

  let dialogEl: HTMLDialogElement | undefined = $state()

  let formData = $state<RuleFormData>({
    match_track_name: '',
    match_artist_name: '',
    match_release_name: '',
    match_artist_names: '',
    match_duration_bucket: '',
    match_mbid: '',
    replace_track_name: '',
    replace_artist_name: '',
    replace_release_name: '',
    replace_artist_names: '',
    enabled: true,
  })

  $effect(() => {
    if (open && dialogEl && !dialogEl.open) {
      resetForm()
      dialogEl.showModal()
    } else if (!open && dialogEl?.open) {
      dialogEl.close()
    }
  })

  function resetForm(): void {
    if (rule) {
      formData = {
        match_track_name: rule.match_track_name ?? '',
        match_artist_name: rule.match_artist_name ?? '',
        match_release_name: rule.match_release_name ?? '',
        match_artist_names: rule.match_artist_names ? JSON.stringify(rule.match_artist_names) : '',
        match_duration_bucket: rule.match_duration_bucket?.toString() ?? '',
        match_mbid: rule.match_mbid ?? '',
        replace_track_name: rule.replace_track_name ?? '',
        replace_artist_name: rule.replace_artist_name ?? '',
        replace_release_name: rule.replace_release_name ?? '',
        replace_artist_names: rule.replace_artist_names ? JSON.stringify(rule.replace_artist_names) : '',
        enabled: rule.enabled,
      }
    } else {
      formData = {
        match_track_name: '',
        match_artist_name: '',
        match_release_name: '',
        match_artist_names: '',
        match_duration_bucket: '',
        match_mbid: '',
        replace_track_name: '',
        replace_artist_name: '',
        replace_release_name: '',
        replace_artist_names: '',
        enabled: true,
      }
    }
  }

  function handleSubmit(e: Event): void {
    e.preventDefault()
    onSave(formData)
  }

  function closeDialog(): void {
    open = false
  }
</script>

<dialog class="modal" bind:this={dialogEl} onclose={closeDialog}>
  <div class="modal-box max-w-xl">
    <div class="mb-6">
      <h3 class="text-lg font-bold">{rule ? 'Edit Rule' : 'New Rule'}</h3>
      <p class="text-sm text-base-content/60">
        {rule ? 'Update the rule criteria and replacements.' : 'Create a new metadata correction rule.'}
      </p>
    </div>

    <form onsubmit={handleSubmit} class="space-y-4">
      <div class="space-y-3">
        <h4 class="text-xs font-semibold uppercase tracking-wider text-base-content/60">Match Criteria</h4>
        <div class="grid grid-cols-2 gap-3">
          <fieldset class="fieldset">
            <legend class="fieldset-label">Track Name</legend>
            <input class="input w-full" bind:value={formData.match_track_name} placeholder="e.g. Bohemian Rhapsody" />
          </fieldset>
          <fieldset class="fieldset">
            <legend class="fieldset-label">Artist Name</legend>
            <input class="input w-full" bind:value={formData.match_artist_name} placeholder="e.g. Queen" />
          </fieldset>
          <fieldset class="fieldset">
            <legend class="fieldset-label">Release Name</legend>
            <input class="input w-full" bind:value={formData.match_release_name} placeholder="e.g. A Night at the Opera" />
          </fieldset>
          <fieldset class="fieldset">
            <legend class="fieldset-label">Artist Names (JSON)</legend>
            <input class="input w-full" bind:value={formData.match_artist_names} placeholder='["Queen"]' />
          </fieldset>
          <fieldset class="fieldset">
            <legend class="fieldset-label">MBID</legend>
            <input class="input w-full" bind:value={formData.match_mbid} placeholder="e.g. 0603a798-..." />
          </fieldset>
          <fieldset class="fieldset">
            <legend class="fieldset-label">Duration Bucket</legend>
            <input class="input w-full" type="number" bind:value={formData.match_duration_bucket} placeholder="e.g. 107" />
          </fieldset>
        </div>
      </div>

      <div class="space-y-3">
        <h4 class="text-xs font-semibold uppercase tracking-wider text-base-content/60">Replacement Values</h4>
        <div class="grid grid-cols-2 gap-3">
          <fieldset class="fieldset">
            <legend class="fieldset-label">Track Name</legend>
            <input class="input w-full" bind:value={formData.replace_track_name} placeholder="Corrected title" />
          </fieldset>
          <fieldset class="fieldset">
            <legend class="fieldset-label">Artist Name</legend>
            <input class="input w-full" bind:value={formData.replace_artist_name} placeholder="Corrected artist" />
          </fieldset>
          <fieldset class="fieldset">
            <legend class="fieldset-label">Release Name</legend>
            <input class="input w-full" bind:value={formData.replace_release_name} placeholder="Corrected release" />
          </fieldset>
          <fieldset class="fieldset">
            <legend class="fieldset-label">Artist Names (JSON)</legend>
            <input class="input w-full" bind:value={formData.replace_artist_names} placeholder='["Corrected Artist"]' />
          </fieldset>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <input type="checkbox" class="toggle" bind:checked={formData.enabled} id="enabled" />
        <label for="enabled" class="cursor-pointer text-sm">Enabled</label>
      </div>

      <div class="flex justify-end gap-3 border-t border-base-300 pt-4">
        <button type="button" class="btn btn-ghost" onclick={closeDialog}>Cancel</button>
        <button type="submit" class="btn btn-primary">Save</button>
      </div>
    </form>
  </div>

  <form method="dialog" class="modal-backdrop">
    <button onclick={closeDialog}>close</button>
  </form>
</dialog>
