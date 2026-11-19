package models

import (
	"unicode"
	"unicode/utf8"
)

// Emoji is a chore's emoji, or "" for none. The server has no emoji dataset, so
// Valid is a coarse safety check (short, valid UTF-8, no letters or whitespace).
// The client enforces that the value is exactly one emoji.
type Emoji string

const maxEmojiBytes = 32

func (e Emoji) Valid() bool {
	if e == "" {
		return true
	}
	if len(e) > maxEmojiBytes || !utf8.ValidString(string(e)) {
		return false
	}
	for _, r := range string(e) {
		if unicode.IsSpace(r) || (r < utf8.RuneSelf && unicode.IsLetter(r)) {
			return false
		}
	}
	return true
}
