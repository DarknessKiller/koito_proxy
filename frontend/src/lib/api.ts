import type { Rule, RuleCreateRequest, AuthCheckResponse } from './types'

const BASE = '/apis/admin'

class ApiError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

async function request<T>(
  path: string,
  options?: RequestInit
): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    credentials: 'same-origin',
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
    ...options,
  })

  if (res.status === 401) {
    throw new ApiError('Not authenticated', 401)
  }

  if (!res.ok) {
    const text = await res.text().catch(() => 'Unknown error')
    throw new ApiError(text || `HTTP ${res.status}`, res.status)
  }

  if (res.status === 204) {
    return undefined as T
  }

  return res.json() as Promise<T>
}

export async function checkAuth(): Promise<AuthCheckResponse> {
  return request<AuthCheckResponse>('/check')
}

export async function fetchRules(): Promise<Rule[]> {
  return request<Rule[]>('/rules')
}

export async function fetchRule(id: string): Promise<Rule> {
  return request<Rule>(`/rules/${encodeURIComponent(id)}`)
}

export async function createRule(data: RuleCreateRequest): Promise<Rule> {
  return request<Rule>('/rules', {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function updateRule(id: string, data: RuleCreateRequest): Promise<Rule> {
  return request<Rule>(`/rules/${encodeURIComponent(id)}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

export async function deleteRule(id: string): Promise<void> {
  return request<void>(`/rules/${encodeURIComponent(id)}`, {
    method: 'DELETE',
  })
}

export { ApiError }
