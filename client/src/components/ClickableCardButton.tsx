import ChevronRightIcon from '@mui/icons-material/ChevronRight'
import { ListItem, ListItemButton, Paper } from '@mui/material'
import type { ReactNode } from 'react'

type ClickableCardButtonProps = {
  onClick: () => void
  children: ReactNode
}

// A raised, tappable row card with a trailing chevron. Children fill the row.
export default function ClickableCardButton({
  onClick,
  children,
}: ClickableCardButtonProps) {
  return (
    <ListItem disablePadding sx={{ mb: 1 }}>
      <Paper
        elevation={2}
        sx={{ width: '100%', borderRadius: 3, overflow: 'hidden' }}
      >
        <ListItemButton onClick={onClick} sx={{ p: 2 }}>
          {children}
          <ChevronRightIcon color="action" />
        </ListItemButton>
      </Paper>
    </ListItem>
  )
}
