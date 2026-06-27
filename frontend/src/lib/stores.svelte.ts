import type { Rule } from './types'

let rules = $state<Rule[]>([])
let searchQuery = $state('')
let page = $state(1)
let pageSize = $state(10)
let loading = $state(false)
let error = $state<string | null>(null)

function setRules(value: Rule[]): void {
  rules = value
}

function setSearchQuery(value: string): void {
  searchQuery = value
  page = 1
}

function setPage(value: number): void {
  page = value
}

function setPageSize(value: number): void {
  pageSize = value
  page = 1
}

function setLoading(value: boolean): void {
  loading = value
}

function setError(value: string | null): void {
  error = value
}

export function getStores() {
  return {
    get rules() { return rules },
    set rules(value: Rule[]) { setRules(value) },
    get searchQuery() { return searchQuery },
    set searchQuery(value: string) { setSearchQuery(value) },
    get page() { return page },
    set page(value: number) { setPage(value) },
    get pageSize() { return pageSize },
    set pageSize(value: number) { setPageSize(value) },
    get loading() { return loading },
    set loading(value: boolean) { setLoading(value) },
    get error() { return error },
    set error(value: string | null) { setError(value) },
  }
}
