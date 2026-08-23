import type { DisplayUnit } from '../generated/api/api.schemas'

export const DISPLAY_UNITS: DisplayUnit[] = ['DAY', 'WEEK', 'MONTH', 'YEAR']

const DAYS_PER_UNIT: Record<DisplayUnit, number> = {
  DAY: 1,
  WEEK: 7,
  MONTH: 30,
  YEAR: 365,
}

const UNIT_WORDS: Record<
  DisplayUnit,
  { label: string; singular: string; plural: string }
> = {
  DAY: { label: 'Day', singular: 'day', plural: 'days' },
  WEEK: { label: 'Week', singular: 'week', plural: 'weeks' },
  MONTH: { label: 'Month', singular: 'month', plural: 'months' },
  YEAR: { label: 'Year', singular: 'year', plural: 'years' },
}

export function daysPerUnit(unit: DisplayUnit): number {
  return DAYS_PER_UNIT[unit]
}

export function unitLabel(unit: DisplayUnit): string {
  return UNIT_WORDS[unit].label
}

// Exact by construction: every real chore's cadenceInDays is
// amount * daysPerUnit(displayUnit) for some integer amount, as long as
// every edit path goes through toDays(). Rounded defensively regardless.
export function toAmount(days: number, unit: DisplayUnit): number {
  return Math.round(days / daysPerUnit(unit))
}

export function toDays(amount: number, unit: DisplayUnit): number {
  return Math.round(amount) * daysPerUnit(unit)
}

interface CadenceProps {
  days: number
  unit: DisplayUnit
}

// Handles both exact chore cadences and fractional/projected values (e.g.
// insights' projectedCadence) by rounding the unit-amount to 1 decimal.
export function formatCadence(days: number, unit: DisplayUnit): string {
  if (days <= 0) return 'immediately'

  const rawAmount = days / daysPerUnit(unit)
  const rounded = Math.round(rawAmount * 10) / 10
  const words = UNIT_WORDS[unit]

  if (rounded === 1) return `every ${words.singular}`

  const formatted = Number.isInteger(rounded)
    ? rounded.toString()
    : rounded.toFixed(1)

  return `every ${formatted} ${words.plural}`
}

export function Cadence({ days, unit }: CadenceProps) {
  return <span>{formatCadence(days, unit)}</span>
}
