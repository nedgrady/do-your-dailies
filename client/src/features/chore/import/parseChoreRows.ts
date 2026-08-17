export interface ParsedChoreRow {
  name: string
  cadenceInDays: number | null
}

function splitLine(line: string): string[] {
  if (line.includes('\t')) return line.split('\t')
  return line.split(',')
}

function isHeaderRow(cells: string[]): boolean {
  if (cells.length < 2) return false
  const cadenceCell = cells[1].trim()
  return cadenceCell !== '' && Number.isNaN(Number(cadenceCell))
}

function toRow(cells: string[]): ParsedChoreRow {
  const name = (cells[0] ?? '').trim()
  const cadenceRaw = cells[1]?.trim()
  const cadence = cadenceRaw ? Number(cadenceRaw) : NaN
  return {
    name,
    cadenceInDays: Number.isFinite(cadence) && cadence > 0 ? cadence : null,
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
