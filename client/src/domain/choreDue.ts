const MINUTE_MS = 60_000
const HOUR_MS = 60 * MINUTE_MS
const DAY_MS = 24 * HOUR_MS

function formatAmount(absMs: number): string {
  if (absMs < HOUR_MS) return `${Math.floor(absMs / MINUTE_MS)}m`
  if (absMs < DAY_MS) return `${Math.floor(absMs / HOUR_MS)}h`
  return `${Math.floor(absMs / DAY_MS)}d`
}

// Pure: the caller supplies `now`, so the result depends only on its inputs.
export function describeDue(dueAt: Date, now: Date) {
  const diffMs = dueAt.getTime() - now.getTime()
  const isOverdue = diffMs < 0
  const amount = formatAmount(Math.abs(diffMs))
  return {
    isOverdue,
    label: isOverdue ? `${amount} over` : `${amount} left`,
  }
}
