import { describe, expect, it } from 'vitest'
import { parseLastCompletedAgo } from './parseLastCompletedAgo'

const DAY_MS = 24 * 60 * 60 * 1000

describe('parseLastCompletedAgo', () => {
  it('returns null for blank input', () => {
    expect(parseLastCompletedAgo('')).toBeNull()
    expect(parseLastCompletedAgo('   ')).toBeNull()
  })

  it.each([
    ['3d', 3],
    ['3 days', 3],
    ['3 days ago', 3],
    ['3d ago', 3],
  ])('parses "%s" as %d day(s) ago', (raw, days) => {
    const result = parseLastCompletedAgo(raw)
    expect(result).not.toBeNull()
    expect(Math.round((Date.now() - (result as Date).getTime()) / DAY_MS)).toBe(
      days,
    )
  })

  it('supports weeks/months/years, with or without "ago"', () => {
    const threeWeeks = parseLastCompletedAgo('3 weeks ago')
    expect(threeWeeks).not.toBeNull()
    expect(
      Math.round((Date.now() - (threeWeeks as Date).getTime()) / DAY_MS),
    ).toBe(21)

    const oneMonth = parseLastCompletedAgo('1 month')
    expect(oneMonth).not.toBeNull()
    expect(
      Math.round((Date.now() - (oneMonth as Date).getTime()) / DAY_MS),
    ).toBe(30)
  })

  it('is case-insensitive on the "ago" suffix', () => {
    expect(parseLastCompletedAgo('3 days AGO')).not.toBeNull()
  })

  it.each([['garbage'], ['3 fortnights ago'], ['ago']])(
    'rejects unparseable input "%s"',
    (raw) => {
      expect(parseLastCompletedAgo(raw)).toBeNull()
    },
  )
})
