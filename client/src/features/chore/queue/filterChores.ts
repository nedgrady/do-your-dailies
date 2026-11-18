import type { Chore } from '../../../domain/chore'

// Case-insensitive "contains" match on chore names. An empty query keeps all.
export function filterChores(chores: Chore[], query: string): Chore[] {
  const needle = query.trim().toLowerCase()
  if (!needle) return chores
  return chores.filter((chore) => chore.name.toLowerCase().includes(needle))
}
