import { Box } from '@mui/material'
import type { Emoji } from '../domain/emoji'

type ChoreAvatarProps = {
  name: string
  emoji: Emoji | null
}

// A chore's visual identity: its emoji when one is set, otherwise its initial on
// a soft tile. Decorative, so hidden from assistive tech; the row text already
// names the chore.
export default function ChoreAvatar({ name, emoji }: ChoreAvatarProps) {
  if (emoji) {
    return (
      <Box
        aria-hidden="true"
        sx={{
          width: 48,
          height: 48,
          mr: 2,
          flexShrink: 0,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          fontSize: 32,
        }}
      >
        {emoji}
      </Box>
    )
  }

  const initial = name.trim().charAt(0).toUpperCase() || '?'
  return (
    <Box
      aria-hidden="true"
      sx={{
        width: 48,
        height: 48,
        mr: 2,
        flexShrink: 0,
        borderRadius: 2,
        bgcolor: 'action.hover',
        color: 'text.secondary',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        fontSize: 20,
        fontWeight: 600,
      }}
    >
      {initial}
    </Box>
  )
}
