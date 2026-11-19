import { z } from 'zod'

// The one definition of "a single emoji": exactly one grapheme that starts with
// a pictograph, a flag (two regional indicators), or a keycap (e.g. 1️⃣).
// Plain digits, letters and spaces are rejected.
const segmenter = new Intl.Segmenter('en', { granularity: 'grapheme' })
const EMOJI_START =
  /^(?:\p{Extended_Pictographic}|\p{Regional_Indicator}{2}|[0-9#*]️?⃣)/u

export function isEmoji(value: string): boolean {
  const graphemes = [...segmenter.segment(value)]
  return graphemes.length === 1 && EMOJI_START.test(graphemes[0].segment)
}

// Branded so a plain string can't be passed where an emoji is required.
export const emojiSchema = z
  .string()
  .refine(isEmoji, { message: 'Must be a single emoji' })
  .brand<'Emoji'>()

export type Emoji = z.infer<typeof emojiSchema>

// The API uses "" for "no emoji". Everywhere else in the client it is null.
export const emojiOrNoneSchema = z
  .union([z.literal(''), emojiSchema])
  .transform((value): Emoji | null => (value === '' ? null : value))

export function toApiEmoji(emoji: Emoji | null): string {
  return emoji ?? ''
}
