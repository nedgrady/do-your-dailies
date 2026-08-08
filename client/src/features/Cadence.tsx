interface CadenceProps {
  days: number
}

export function formatCadence(days: number): string {
  if (days <= 0) return 'immediately'
  if (days === 1) return 'every day'

  // Round to 1 decimal, but drop the decimal if it's a whole number
  const rounded = Math.round(days * 10) / 10
  const formatted = Number.isInteger(rounded)
    ? rounded.toString()
    : rounded.toFixed(1)

  // Nice-sounding common cases
  if (rounded === 7) return 'every week'
  if (rounded === 14) return 'every 2 weeks'
  if (rounded === 30 || rounded === 31) return 'every month'

  return `every ${formatted} days`
}

export function Cadence({ days }: CadenceProps) {
  return <span>{formatCadence(days)}</span>
}
