import Papa from 'papaparse'
import type { DisplayUnit } from '../../../generated/api/api.schemas'
import type { ParsedCadence } from './parseCadenceText'
import { parseCadenceParts, parseCadenceText } from './parseCadenceText'

export interface ParsedChoreRow {
  name: string
  cadenceInDays: number | null
  displayUnit: DisplayUnit
  // The raw cadence text as pasted, so an invalid row can quote back what
  // the user actually typed instead of just showing a bare "0".
  cadenceRaw: string
  // A structural problem PapaParse itself flagged for this row (e.g.
  // mismatched quotes), distinct from a cadence that simply doesn't parse.
  parseError: string | null
}

export interface ParseDelimiters {
  // '' or omitted means auto-detect (tab, comma, semicolon, pipe...).
  field?: string
  // Omitted means '\n' (the only row separator PapaParse itself supports).
  row?: string
}

// A row's cadence can either be one combined cell ("3d", "3 days") or split
// across two cells (amount, unit) — e.g. a 3-column spreadsheet paste of
// name/amount/unit. A non-blank third cell means the split form is in use.
function parseRowCadence(cells: string[]): ParsedCadence | null {
  const amountCell = cells[1] ?? ''
  const unitCell = cells[2] ?? ''
  if (unitCell !== '') return parseCadenceParts(amountCell, unitCell)
  return amountCell ? parseCadenceText(amountCell) : null
}

// A header's cadence-ish cells ("Cadence", "Amount", "Unit") never start with
// a digit — even a malformed amount like "4," still looks like data, so a
// typo in the very first row doesn't get misclassified as a header and
// silently dropped instead of surfacing as an invalid row.
function looksLikeCadenceData(cell: string): boolean {
  return /^\d/.test(cell)
}

function isHeaderRow(cells: string[]): boolean {
  if (cells.length < 2) return false
  const amountCell = cells[1] ?? ''
  const unitCell = cells[2] ?? ''
  if (amountCell === '' && unitCell === '') return false
  return !looksLikeCadenceData(amountCell) && !looksLikeCadenceData(unitCell)
}

function describeCadenceCells(cells: string[]): string {
  const amountCell = cells[1] ?? ''
  const unitCell = cells[2] ?? ''
  return [amountCell, unitCell].filter((cell) => cell !== '').join(' ')
}

function toRow(cells: string[], parseError: string | null): ParsedChoreRow {
  const name = cells[0] ?? ''
  const parsed = parseRowCadence(cells)
  return {
    name,
    cadenceInDays: parsed?.cadenceInDays ?? null,
    displayUnit: parsed?.displayUnit ?? 'DAY',
    cadenceRaw: describeCadenceCells(cells),
    parseError,
  }
}

// Column delimiter is auto-detected by default (tab, comma, semicolon,
// pipe...) so pasting from Sheets/Excel (tabs) or a plain CSV (commas) both
// just work, and a quoted field containing the delimiter itself is handled
// correctly. `delimiters` lets the caller override either choice explicitly
// (e.g. a user-facing picker) instead of relying on auto-detection.
export function parseChoreRows(
  text: string,
  delimiters: ParseDelimiters = {},
): ParsedChoreRow[] {
  const rowDelimiter = delimiters.row ?? '\n'
  // PapaParse's own `newline` option only accepts '\r'/'\n'/'\r\n' — it can't
  // be told to treat e.g. a semicolon as a row separator. So a non-default
  // row delimiter is normalized to '\n' before Papa ever sees the text; the
  // field delimiter is passed straight through to Papa's own `delimiter`
  // option, which already accepts an arbitrary string (or '' for auto-detect).
  const normalized =
    rowDelimiter === '\n' ? text : text.split(rowDelimiter).join('\n')

  const { data, errors } = Papa.parse<string[]>(normalized, {
    delimiter: delimiters.field ?? '',
    skipEmptyLines: 'greedy',
    transform: (value) => value.trim(),
  })

  if (data.length === 0) return []

  // Keep the first error PapaParse reported for each row (indexed into
  // `data`, before any header row is stripped below).
  const errorsByRow = new Map<number, string>()
  for (const error of errors) {
    if (error.row !== undefined && !errorsByRow.has(error.row)) {
      errorsByRow.set(error.row, error.message)
    }
  }

  const hasHeader = isHeaderRow(data[0])
  const rowLines = hasHeader ? data.slice(1) : data
  const rowOffset = hasHeader ? 1 : 0

  return rowLines.map((cells, index) =>
    toRow(cells, errorsByRow.get(index + rowOffset) ?? null),
  )
}
