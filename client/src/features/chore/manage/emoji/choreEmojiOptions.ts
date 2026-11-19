import { type Emoji, emojiSchema } from '../../../../domain/emoji'

export type ChoreEmojiOption = { emoji: Emoji; name: string }

// Each entry is an emoji and the words people search for it by. Entries go
// through emojiSchema at load, so a bad one fails loudly.
const RAW_CHORE_EMOJI: ReadonlyArray<readonly [string, string]> = [
  ['🧹', 'sweep broom floor'],
  ['🧽', 'sponge scrub'],
  ['🧼', 'soap wash clean'],
  ['🧺', 'basket laundry washing'],
  ['🪣', 'bucket mop'],
  ['🚿', 'shower'],
  ['🛁', 'bath bathroom bathtub'],
  ['🚽', 'toilet loo'],
  ['🪟', 'window windows glass'],
  ['🧴', 'lotion bottle cleaner'],
  ['🧻', 'toilet paper tissue'],
  ['🫧', 'bubbles'],
  ['🍳', 'cook pan frying kitchen'],
  ['🍽️', 'dinner plate dishes'],
  ['🔪', 'knife cut chop'],
  ['🧊', 'ice fridge freezer'],
  ['🥣', 'bowl cereal'],
  ['☕', 'coffee tea'],
  ['🫖', 'teapot tea'],
  ['🍞', 'bread bake'],
  ['🥕', 'carrot vegetables veg'],
  ['🥛', 'milk'],
  ['🧃', 'juice drink'],
  ['🍲', 'pot stew soup cooking'],
  ['🌱', 'seedling plant water'],
  ['🪴', 'houseplant pot plants water'],
  ['🌿', 'herb leaves plant'],
  ['🍂', 'leaves autumn rake'],
  ['🌸', 'blossom flower'],
  ['🌳', 'tree'],
  ['🚗', 'car'],
  ['🚲', 'bike bicycle'],
  ['🗑️', 'bin trash rubbish bins'],
  ['♻️', 'recycling recycle'],
  ['🪵', 'wood logs firewood'],
  ['🏡', 'house home'],
  ['🐶', 'dog walk'],
  ['🐱', 'cat litter feed'],
  ['🐟', 'fish tank'],
  ['🐦', 'bird'],
  ['🐹', 'hamster'],
  ['🐰', 'rabbit'],
  ['🛏️', 'bed sheets bedroom'],
  ['👕', 'shirt clothes laundry'],
  ['🧦', 'socks'],
  ['🧥', 'coat jacket'],
  ['📦', 'parcel box package'],
  ['📬', 'post mail letters'],
  ['💳', 'card payment bills'],
  ['🧾', 'receipt bills invoice'],
  ['🔧', 'wrench fix repair'],
  ['🪛', 'screwdriver'],
  ['💡', 'bulb light'],
  ['🔑', 'key'],
  ['📅', 'calendar'],
  ['🛒', 'shopping cart groceries'],
  ['🧯', 'extinguisher fire'],
  ['🪜', 'ladder'],
]

export const CHORE_EMOJI_OPTIONS: readonly ChoreEmojiOption[] = RAW_CHORE_EMOJI.map(
  ([emoji, name]) => ({ emoji: emojiSchema.parse(emoji), name }),
)

// Every word typed must appear in the entry's words, so "bath toilet" narrows.
export function searchChoreEmoji(query: string): readonly ChoreEmojiOption[] {
  const terms = query.trim().toLowerCase().split(/\s+/).filter(Boolean)
  if (terms.length === 0) return CHORE_EMOJI_OPTIONS
  return CHORE_EMOJI_OPTIONS.filter(({ name }) =>
    terms.every((term) => name.includes(term)),
  )
}
