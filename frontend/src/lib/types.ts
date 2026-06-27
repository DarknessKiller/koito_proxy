export interface Rule {
  id: string
  match_track_name: string | null
  match_artist_name: string | null
  match_release_name: string | null
  match_artist_names: string[] | null
  match_duration_bucket: number | null
  match_mbid: string | null
  replace_track_name: string | null
  replace_artist_name: string | null
  replace_release_name: string | null
  replace_artist_names: string[] | null
  enabled: boolean
  priority: number
  valid: boolean
  created_at: string
  updated_at: string
}

export interface RuleCreateRequest {
  match_track_name?: string | null
  match_artist_name?: string | null
  match_release_name?: string | null
  match_artist_names?: string[] | null
  match_duration_bucket?: number | null
  match_mbid?: string | null
  replace_track_name?: string | null
  replace_artist_name?: string | null
  replace_release_name?: string | null
  replace_artist_names?: string[] | null
  enabled?: boolean
}

export interface RuleFormData {
  match_track_name: string
  match_artist_name: string
  match_release_name: string
  match_artist_names: string
  match_duration_bucket: string
  match_mbid: string
  replace_track_name: string
  replace_artist_name: string
  replace_release_name: string
  replace_artist_names: string
  enabled: boolean
}

export interface AuthCheckResponse {
  ok: boolean
}
