import { parseCadenceText } from './parseCadenceText'

// "3 weeks ago", "10d ago", "2 months" -> a Date that many amount-units in
// the past. Reuses parseCadenceText's amount+unit grammar (its "cadenceInDays"
// is exactly the day-count we want here too, just measuring time already
// elapsed instead of a recurring interval) after stripping an optional
// trailing "ago".
export function parseLastCompletedAgo(raw: string): Date | null {
  const trimmed = raw.trim()
  if (trimmed === '') return null

  const withoutAgo = trimmed.replace(/\s*ago\s*$/i, '')
  const parsed = parseCadenceText(withoutAgo)
  if (!parsed) return null

  return new Date(Date.now() - parsed.cadenceInDays * 24 * 60 * 60 * 1000)
}
