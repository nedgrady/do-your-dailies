import { describe, expect, it } from 'vitest'
import { parseChoreRows } from './parseChoreRows'

describe('parseChoreRows', () => {
  it('does not drop the first row as a false header when it uses a unit suffix', () => {
    const rows = parseChoreRows('Dishes\t3d\nVacuum\t1w')

    expect(rows).toEqual([
      {
        name: 'Dishes',
        cadenceInDays: 3,
        displayUnit: 'DAY',
        cadenceRaw: '3d',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
      {
        name: 'Vacuum',
        cadenceInDays: 7,
        displayUnit: 'WEEK',
        cadenceRaw: '1w',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
    ])
  })

  it('still strips a real header row', () => {
    const rows = parseChoreRows(
      'Chore\tCadence\nDishes\t3\nVacuum\t2 weeks\nMow\t1m\nTaxes\t1y',
    )

    expect(rows).toEqual([
      {
        name: 'Dishes',
        cadenceInDays: 3,
        displayUnit: 'DAY',
        cadenceRaw: '3',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
      {
        name: 'Vacuum',
        cadenceInDays: 14,
        displayUnit: 'WEEK',
        cadenceRaw: '2 weeks',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
      {
        name: 'Mow',
        cadenceInDays: 30,
        displayUnit: 'MONTH',
        cadenceRaw: '1m',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
      {
        name: 'Taxes',
        cadenceInDays: 365,
        displayUnit: 'YEAR',
        cadenceRaw: '1y',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
    ])
  })

  it('supports comma-delimited rows', () => {
    const rows = parseChoreRows('Dishes,3d\nVacuum,1w')

    expect(rows).toEqual([
      {
        name: 'Dishes',
        cadenceInDays: 3,
        displayUnit: 'DAY',
        cadenceRaw: '3d',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
      {
        name: 'Vacuum',
        cadenceInDays: 7,
        displayUnit: 'WEEK',
        cadenceRaw: '1w',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
    ])
  })

  it('flags an unparseable cadence instead of dropping the row, and keeps the raw text', () => {
    // A garbage cadence on row 0 alone is indistinguishable from a header
    // (true both before and after this change), so use a second row here.
    const rows = parseChoreRows('Dishes\t3d\nWeird\t3 x y')

    expect(rows).toEqual([
      {
        name: 'Dishes',
        cadenceInDays: 3,
        displayUnit: 'DAY',
        cadenceRaw: '3d',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
      {
        name: 'Weird',
        cadenceInDays: null,
        displayUnit: 'DAY',
        cadenceRaw: '3 x y',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
    ])
  })

  it('does not misclassify a malformed-but-numeric-looking first row as a header', () => {
    // Regression test: a typo like "4," used to fail parseRowCadence, and the
    // old isHeaderRow heuristic ("cadence cell doesn't parse") treated any
    // unparseable first row as a header and silently dropped it — even
    // though "4," clearly looks like data, not a header label.
    const rows = parseChoreRows(
      'test from sheets\t4,\tweek\ntest from sheets\t2\tday',
    )

    expect(rows).toEqual([
      {
        name: 'test from sheets',
        cadenceInDays: null,
        displayUnit: 'DAY',
        cadenceRaw: '4, week',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
      {
        name: 'test from sheets',
        cadenceInDays: 2,
        displayUnit: 'DAY',
        cadenceRaw: '2 day',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
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
      {
        name: 'test from sheets',
        cadenceInDays: 28,
        displayUnit: 'WEEK',
        cadenceRaw: '4 week',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
      {
        name: 'test from sheets 2',
        cadenceInDays: 2,
        displayUnit: 'DAY',
        cadenceRaw: '2 day',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
    ])
  })

  it('strips a header row for the three-column format too', () => {
    const rows = parseChoreRows(
      'Chore\tAmount\tUnit\nDishes\t3\tday\nVacuum\t1\tweek',
    )

    expect(rows).toEqual([
      {
        name: 'Dishes',
        cadenceInDays: 3,
        displayUnit: 'DAY',
        cadenceRaw: '3 day',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
      {
        name: 'Vacuum',
        cadenceInDays: 7,
        displayUnit: 'WEEK',
        cadenceRaw: '1 week',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
    ])
  })

  it('flags an unrecognized unit column instead of dropping the row', () => {
    const rows = parseChoreRows('Dishes\t3d\nWeird\t3\tfortnight')

    expect(rows).toEqual([
      {
        name: 'Dishes',
        cadenceInDays: 3,
        displayUnit: 'DAY',
        cadenceRaw: '3d',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
      {
        name: 'Weird',
        cadenceInDays: null,
        displayUnit: 'DAY',
        cadenceRaw: '3 fortnight',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
    ])
  })

  it('honors an explicit field delimiter over auto-detection', () => {
    // A comma inside the name would confuse comma auto-detection, but an
    // explicit tab delimiter sidesteps that entirely.
    const rows = parseChoreRows('Dishes, indoor\t3d', { field: '\t' })

    expect(rows).toEqual([
      {
        name: 'Dishes, indoor',
        cadenceInDays: 3,
        displayUnit: 'DAY',
        cadenceRaw: '3d',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
    ])
  })

  it('honors an explicit row delimiter', () => {
    const rows = parseChoreRows('Dishes\t3d;Vacuum\t1w', { row: ';' })

    expect(rows).toEqual([
      {
        name: 'Dishes',
        cadenceInDays: 3,
        displayUnit: 'DAY',
        cadenceRaw: '3d',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
      {
        name: 'Vacuum',
        cadenceInDays: 7,
        displayUnit: 'WEEK',
        cadenceRaw: '1w',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
    ])
  })

  it('carries an optional 4th "last done" cell through, split-cadence form', () => {
    const rows = parseChoreRows('Mow lawn\t2\tweek\t3 weeks ago')

    expect(rows).toEqual([
      {
        name: 'Mow lawn',
        cadenceInDays: 14,
        displayUnit: 'WEEK',
        cadenceRaw: '2 week',
        parseError: null,
        lastCompletedAgoRaw: '3 weeks ago',
      },
    ])
  })

  it('carries an optional "last done" cell through, combined-cadence form with a blank unit cell', () => {
    const rows = parseChoreRows('Mow lawn\t2w\t\t3 weeks ago')

    expect(rows).toEqual([
      {
        name: 'Mow lawn',
        cadenceInDays: 14,
        displayUnit: 'WEEK',
        cadenceRaw: '2w',
        parseError: null,
        lastCompletedAgoRaw: '3 weeks ago',
      },
    ])
  })

  it('defaults "last done" to blank when the 4th cell is absent', () => {
    const rows = parseChoreRows('Dishes\t3d')

    expect(rows[0].lastCompletedAgoRaw).toBe('')
  })

  it('honors both an explicit field and row delimiter together', () => {
    const rows = parseChoreRows('Dishes,3d;Vacuum,1w', {
      field: ',',
      row: ';',
    })

    expect(rows).toEqual([
      {
        name: 'Dishes',
        cadenceInDays: 3,
        displayUnit: 'DAY',
        cadenceRaw: '3d',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
      {
        name: 'Vacuum',
        cadenceInDays: 7,
        displayUnit: 'WEEK',
        cadenceRaw: '1w',
        parseError: null,
        lastCompletedAgoRaw: '',
      },
    ])
  })
})
