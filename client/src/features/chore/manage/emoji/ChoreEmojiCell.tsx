import { ButtonBase, Popover } from '@mui/material'
import { useState } from 'react'
import type { Emoji } from '../../../../domain/emoji'
import ChoreEmojiPicker from './ChoreEmojiPicker'

interface ChoreEmojiCellProps {
  name: string
  emoji: Emoji | null
  onChange: (emoji: Emoji) => void
}

// A grid cell that shows the chore's emoji (or its initial). Clicking it opens
// the search-and-pick list in a popover anchored to the cell.
export default function ChoreEmojiCell({
  name,
  emoji,
  onChange,
}: ChoreEmojiCellProps) {
  const [anchor, setAnchor] = useState<HTMLElement | null>(null)
  const initial = name.trim().charAt(0).toUpperCase() || '?'

  return (
    <>
      <ButtonBase
        aria-label={`Change emoji for ${name}`}
        onClick={(event) => setAnchor(event.currentTarget)}
        sx={{
          width: 40,
          height: 36,
          borderRadius: 1,
          fontSize: 24,
          color: emoji ? 'text.primary' : 'text.secondary',
          bgcolor: emoji ? undefined : 'action.hover',
          '&:hover': { bgcolor: 'action.selected' },
        }}
      >
        {emoji ?? initial}
      </ButtonBase>
      <Popover
        open={anchor !== null}
        anchorEl={anchor}
        onClose={() => setAnchor(null)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'left' }}
        transformOrigin={{ vertical: 'top', horizontal: 'left' }}
        slotProps={{ paper: { sx: { p: 2 } } }}
      >
        <ChoreEmojiPicker
          value={emoji}
          autoFocus
          onPick={(picked) => {
            onChange(picked)
            setAnchor(null)
          }}
        />
      </Popover>
    </>
  )
}
