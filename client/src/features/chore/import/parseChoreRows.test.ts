import { describe, expect, it } from 'vitest'
import { parseChoreRows } from './parseChoreRows'

describe('parseChoreRows', () => {
  it('does not drop the first row as a false header when it uses a unit suffix', () => {
    const rows = parseChoreRows('Dishes\t3d\nVacuum\t1w')

    expect(rows).toEqual([
      { name: 'Dishes', cadenceInDays: 3, displayUnit: 'DAY' },
      { name: 'Vacuum', cadenceInDays: 7, displayUnit: 'WEEK' },
    ])
  })

  it('still strips a real header row', () => {
    const rows = parseChoreRows(
      'Chore\tCadence\nDishes\t3\nVacuum\t2 weeks\nMow\t1m\nTaxes\t1y',
    )

    expect(rows).toEqual([
      { name: 'Dishes', cadenceInDays: 3, displayUnit: 'DAY' },
      { name: 'Vacuum', cadenceInDays: 14, displayUnit: 'WEEK' },
      { name: 'Mow', cadenceInDays: 30, displayUnit: 'MONTH' },
      { name: 'Taxes', cadenceInDays: 365, displayUnit: 'YEAR' },
    ])
  })

  it('supports comma-delimited rows', () => {
    const rows = parseChoreRows('Dishes,3d\nVacuum,1w')

    expect(rows).toEqual([
      { name: 'Dishes', cadenceInDays: 3, displayUnit: 'DAY' },
      { name: 'Vacuum', cadenceInDays: 7, displayUnit: 'WEEK' },
    ])
  })

  it('flags an unparseable cadence instead of dropping the row', () => {
    // A garbage cadence on row 0 alone is indistinguishable from a header
    // (true both before and after this change), so use a second row here.
    const rows = parseChoreRows('Dishes\t3d\nWeird\t3 x y')

    expect(rows).toEqual([
      { name: 'Dishes', cadenceInDays: 3, displayUnit: 'DAY' },
      { name: 'Weird', cadenceInDays: null, displayUnit: 'DAY' },
    ])
  })

  it('returns an empty list for blank input', () => {
    expect(parseChoreRows('   \n  \n')).toEqual([])
  })

  it('supports a separate amount and unit column', () => {
    const rows = parseChoreRows(
      'test from sheets\t4\tweek\ntest from sheets 2\t2\tday',
    )

    expect(rows).toEqual([
      { name: 'test from sheets', cadenceInDays: 28, displayUnit: 'WEEK' },
      { name: 'test from sheets 2', cadenceInDays: 2, displayUnit: 'DAY' },
    ])
  })

  it('strips a header row for the three-column format too', () => {
    const rows = parseChoreRows(
      'Chore\tAmount\tUnit\nDishes\t3\tday\nVacuum\t1\tweek',
    )

    expect(rows).toEqual([
      { name: 'Dishes', cadenceInDays: 3, displayUnit: 'DAY' },
      { name: 'Vacuum', cadenceInDays: 7, displayUnit: 'WEEK' },
    ])
  })

  it('flags an unrecognized unit column instead of dropping the row', () => {
    const rows = parseChoreRows('Dishes\t3d\nWeird\t3\tfortnight')

    expect(rows).toEqual([
      { name: 'Dishes', cadenceInDays: 3, displayUnit: 'DAY' },
      { name: 'Weird', cadenceInDays: null, displayUnit: 'DAY' },
    ])
  })
})
