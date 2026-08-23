import { describe, expect, it } from 'vitest'
import { parseCadenceParts, parseCadenceText } from './parseCadenceText'

describe('parseCadenceText', () => {
  it('defaults a bare number to days', () => {
    expect(parseCadenceText('3')).toEqual({
      cadenceInDays: 3,
      displayUnit: 'DAY',
    })
  })

  it.each([
    ['3d', 3],
    ['3 day', 3],
    ['3 days', 3],
  ])('parses "%s" as %d day(s)', (raw, amount) => {
    expect(parseCadenceText(raw)).toEqual({
      cadenceInDays: amount,
      displayUnit: 'DAY',
    })
  })

  it.each([
    ['2w', 14],
    ['2 week', 14],
    ['2 weeks', 14],
  ])('parses "%s" as weeks', (raw, cadenceInDays) => {
    expect(parseCadenceText(raw)).toEqual({
      cadenceInDays,
      displayUnit: 'WEEK',
    })
  })

  it.each([
    ['1m', 30],
    ['1 month', 30],
    ['2 months', 60],
  ])('parses "%s" as months', (raw, cadenceInDays) => {
    expect(parseCadenceText(raw)).toEqual({
      cadenceInDays,
      displayUnit: 'MONTH',
    })
  })

  it.each([
    ['1y', 365],
    ['1 year', 365],
    ['2 years', 730],
  ])('parses "%s" as years', (raw, cadenceInDays) => {
    expect(parseCadenceText(raw)).toEqual({
      cadenceInDays,
      displayUnit: 'YEAR',
    })
  })

  it('is case-insensitive', () => {
    expect(parseCadenceText('3D')).toEqual({
      cadenceInDays: 3,
      displayUnit: 'DAY',
    })
    expect(parseCadenceText('2 WEEKS')).toEqual({
      cadenceInDays: 14,
      displayUnit: 'WEEK',
    })
  })

  it('tolerates extra whitespace between amount and unit', () => {
    expect(parseCadenceText('3   days')).toEqual({
      cadenceInDays: 3,
      displayUnit: 'DAY',
    })
  })

  it.each([
    [''],
    ['   '],
    ['d'],
    ['3q'],
    ['3xyz'],
    ['0d'],
    ['-3d'],
    ['3.5'],
    ['3 d s'],
  ])('rejects invalid input "%s"', (raw) => {
    expect(parseCadenceText(raw)).toBeNull()
  })
})

describe('parseCadenceParts', () => {
  it('parses a plain amount with a separate unit column', () => {
    expect(parseCadenceParts('4', 'week')).toEqual({
      cadenceInDays: 28,
      displayUnit: 'WEEK',
    })
    expect(parseCadenceParts('2', 'day')).toEqual({
      cadenceInDays: 2,
      displayUnit: 'DAY',
    })
  })

  it('is case-insensitive and tolerates surrounding whitespace', () => {
    expect(parseCadenceParts(' 3 ', ' Months ')).toEqual({
      cadenceInDays: 90,
      displayUnit: 'MONTH',
    })
  })

  it('defaults to DAY when the unit column is blank', () => {
    expect(parseCadenceParts('5', '')).toEqual({
      cadenceInDays: 5,
      displayUnit: 'DAY',
    })
  })

  it.each([
    ['3d', 'week'], // amount cell must be a plain number, not "3d"
    ['3', 'fortnight'],
    ['0', 'week'],
    ['-3', 'week'],
    ['3.5', 'week'],
    ['', 'week'],
  ])('rejects invalid amount/unit "%s"/"%s"', (amount, unit) => {
    expect(parseCadenceParts(amount, unit)).toBeNull()
  })
})
