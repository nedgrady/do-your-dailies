import type { DisplayUnit } from '../../../generated/api/api.schemas'
import type { ParsedCadence } from './parseCadenceText'
import { parseCadenceParts, parseCadenceText } from './parseCadenceText'

export interface ParsedChoreRow {
  name: string
  cadenceInDays: number | null
  displayUnit: DisplayUnit
}

function splitLine(line: string): string[] {
  if (line.includes('\t')) return line.split('\t')
  return line.split(',')
}

// A row's cadence can either be one combined cell ("3d", "3 days") or split
// across two cells (amount, unit) — e.g. a 3-column spreadsheet paste of
// name/amount/unit. A non-blank third cell means the split form is in use.
function parseRowCadence(cells: string[]): ParsedCadence | null {
  const amountCell = cells[1]?.trim() ?? ''
  const unitCell = cells[2]?.trim() ?? ''
  if (unitCell !== '') return parseCadenceParts(amountCell, unitCell)
  return amountCell ? parseCadenceText(amountCell) : null
}

function isHeaderRow(cells: string[]): boolean {
  if (cells.length < 2) return false
  const hasCadenceContent =
    (cells[1]?.trim() ?? '') !== '' || (cells[2]?.trim() ?? '') !== ''
  return hasCadenceContent && parseRowCadence(cells) === null
}

function toRow(cells: string[]): ParsedChoreRow {
  const name = (cells[0] ?? '').trim()
  const parsed = parseRowCadence(cells)
  return {
    name,
    cadenceInDays: parsed?.cadenceInDays ?? null,
    displayUnit: parsed?.displayUnit ?? 'DAY',
  }
}

// Tab-delimited first since pasting cells from Sheets/Excel produces tabs, not commas.
export function parseChoreRows(text: string): ParsedChoreRow[] {
  const lines = text
    .split(/\r?\n/)
    .map((line) => line.trim())
    .filter((line) => line.length > 0)

  if (lines.length === 0) return []

  let rowLines = lines.map(splitLine)
  if (isHeaderRow(rowLines[0])) {
    rowLines = rowLines.slice(1)
  }

  return rowLines.map(toRow)
}
