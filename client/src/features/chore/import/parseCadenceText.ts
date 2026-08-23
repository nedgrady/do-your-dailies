import type { DisplayUnit } from '../../../generated/api/api.schemas'
import { toDays } from '../../Cadence'

// Add a new raw-text spelling here to accept it (e.g. wk: 'WEEK').
const UNIT_ALIASES: Record<string, DisplayUnit | undefined> = {
  d: 'DAY',
  day: 'DAY',
  days: 'DAY',
  w: 'WEEK',
  week: 'WEEK',
  weeks: 'WEEK',
  m: 'MONTH',
  month: 'MONTH',
  months: 'MONTH',
  y: 'YEAR',
  year: 'YEAR',
  years: 'YEAR',
}

export interface ParsedCadence {
  cadenceInDays: number
  displayUnit: DisplayUnit
}

// Looks up a raw unit token (e.g. "d", "Week", "MONTHS") case-insensitively.
export function parseUnitToken(raw: string): DisplayUnit | null {
  return UNIT_ALIASES[raw.trim().toLowerCase()] ?? null
}

function toParsedCadence(
  amount: number,
  unitToken: string,
): ParsedCadence | null {
  if (!Number.isInteger(amount) || amount <= 0) return null

  const displayUnit = unitToken === '' ? 'DAY' : parseUnitToken(unitToken)
  if (!displayUnit) return null

  return { cadenceInDays: toDays(amount, displayUnit), displayUnit }
}

// "3", "3d", "3 days", "2 Weeks" -> amount + unit. No unit suffix defaults
// to DAY, matching the legacy bare-number behavior.
export function parseCadenceText(raw: string): ParsedCadence | null {
  const match = /^(\d+)\s*([a-zA-Z]*)$/.exec(raw.trim())
  if (!match) return null

  return toParsedCadence(Number(match[1]), match[2])
}

// Same as parseCadenceText, but for a spreadsheet paste with the amount and
// unit already split across two separate columns, e.g. "4" / "week".
export function parseCadenceParts(
  amountRaw: string,
  unitRaw: string,
): ParsedCadence | null {
  const amountMatch = /^\d+$/.exec(amountRaw.trim())
  if (!amountMatch) return null

  return toParsedCadence(Number(amountMatch[0]), unitRaw.trim())
}
