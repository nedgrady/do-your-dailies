import { Box, Button, Stack, Typography } from '@mui/material'
import type { Emoji } from '../../../../domain/emoji'
import ChoreEmojiPicker from './ChoreEmojiPicker'

interface EmojiStepProps {
  value: Emoji | null
  onChange: (emoji: Emoji | null) => void
}

export function EmojiStep({ value, onChange }: EmojiStepProps) {
  return (
    <Stack spacing={2}>
      <Stack
        direction="row"
        sx={{ alignItems: 'center', justifyContent: 'space-between' }}
      >
        <Stack direction="row" spacing={2} sx={{ alignItems: 'center' }}>
          <Box
            aria-hidden="true"
            sx={{
              width: 56,
              height: 56,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontSize: 36,
              borderRadius: 2,
              bgcolor: 'action.hover',
            }}
          >
            {value ?? ''}
          </Box>
          <Typography color="text.secondary">
            {value ? 'Current emoji' : 'Optional. Pick one, or save without.'}
          </Typography>
        </Stack>
        {value && (
          <Button color="inherit" onClick={() => onChange(null)}>
            Clear
          </Button>
        )}
      </Stack>
      <ChoreEmojiPicker value={value} onPick={onChange} />
    </Stack>
  )
}
