import { Box, Stack, Typography } from '@mui/material'
import type { ReactNode } from 'react'

type ChoreQueueColumnProps = {
  title: string
  count: number
  footer?: ReactNode
  emptyMessage: string
  children: ReactNode
}

// A plain column: heading, count, the list, and an optional footer. The cards
// carry the visual weight, so the column itself has no fill or border.
export default function ChoreQueueColumn({
  title,
  count,
  footer,
  emptyMessage,
  children,
}: ChoreQueueColumnProps) {
  return (
    <Box sx={{ height: '100%' }}>
      <Stack direction="row" spacing={1} sx={{ alignItems: 'baseline', mb: 1 }}>
        <Typography variant="h6" component="h2">
          {title}
        </Typography>
        <Typography component="span" color="text.secondary">
          {count}
        </Typography>
      </Stack>
      {count === 0 ? (
        <Typography color="text.secondary">{emptyMessage}</Typography>
      ) : (
        children
      )}
      {footer && <Box sx={{ mt: 1 }}>{footer}</Box>}
    </Box>
  )
}
