import { Box, ButtonBase, TextField, Typography } from '@mui/material'
import { useState } from 'react'
import type { Emoji } from '../../../../domain/emoji'
import { searchChoreEmoji } from './choreEmojiOptions'

// Fixed height so the surface doesn't resize as the search narrows the list.
const GRID_HEIGHT = 288

interface ChoreEmojiPickerProps {
  value: Emoji | null
  onPick: (emoji: Emoji) => void
  // Focus search on open. Only for popovers, so mobile doesn't pop the keyboard.
  autoFocus?: boolean
}

// A search box over the curated chore emoji. Instant, because it filters a
// short local list rather than a full emoji library.
export default function ChoreEmojiPicker({
  value,
  onPick,
  autoFocus = false,
}: ChoreEmojiPickerProps) {
  const [query, setQuery] = useState('')
  const options = searchChoreEmoji(query)

  return (
    <Box sx={{ width: 320, maxWidth: '100%' }}>
      <TextField
        fullWidth
        size="small"
        placeholder="Search"
        value={query}
        onChange={(event) => setQuery(event.target.value)}
        autoFocus={autoFocus}
        slotProps={{ htmlInput: { 'aria-label': 'Search emoji' } }}
        sx={{ mb: 1 }}
      />
      <Box
        role="group"
        aria-label="Chore emoji"
        sx={{
          height: GRID_HEIGHT,
          overflowY: 'auto',
          display: 'grid',
          alignContent: 'start',
          gridTemplateColumns: 'repeat(auto-fill, minmax(44px, 1fr))',
          gap: 0.5,
        }}
      >
        {options.map(({ emoji, name }) => (
          <ButtonBase
            key={emoji}
            aria-label={name}
            aria-pressed={emoji === value}
            onClick={() => onPick(emoji)}
            sx={{
              width: 44,
              height: 44,
              fontSize: 26,
              borderRadius: 2,
              bgcolor: emoji === value ? 'action.selected' : undefined,
              '&:hover': { bgcolor: 'action.hover' },
            }}
          >
            {emoji}
          </ButtonBase>
        ))}
        {options.length === 0 && (
          <Typography
            color="text.secondary"
            sx={{ gridColumn: '1 / -1', py: 2, textAlign: 'center' }}
          >
            No match
          </Typography>
        )}
      </Box>
    </Box>
  )
}
