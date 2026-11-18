import {
  Box,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  TextField,
  Typography,
} from '@mui/material'
import { useState } from 'react'
import type { Chore } from '../../../domain/chore'
import { filterChores } from './filterChores'

interface ChoreStepProps {
  chores: Chore[]
  selectedId: number
  onSelect: (choreId: number) => void
}

// Rows shown before asking the user to type. Keeps the dialog from needing
// its own scroll bar on desktop.
const MAX_VISIBLE = 6

// Fixed row height, so the reserved space below matches the rows exactly.
const ROW_HEIGHT = 48

// Lets the user find the chore they just did. Only chores that aren't in the
// queue and haven't been completed today are offered.
export default function ChoreStep({
  chores,
  selectedId,
  onSelect,
}: ChoreStepProps) {
  const [query, setQuery] = useState('')

  if (chores.length === 0) {
    return (
      <Typography color="text.secondary">
        Every other chore is already in today&apos;s queue or done.
      </Typography>
    )
  }

  const matches = filterChores(chores, query)
  const visible = matches.slice(0, MAX_VISIBLE)
  const hiddenCount = matches.length - visible.length
  // Reserve the list and "N more" space up front, so typing a filter never
  // resizes the dialog.
  const reservedRows = Math.min(MAX_VISIBLE, chores.length)
  const hasMoreLine = chores.length > MAX_VISIBLE

  return (
    <>
      {/* Sticky so the filter stays reachable while a long list scrolls. */}
      <Box
        sx={{
          position: 'sticky',
          top: 0,
          zIndex: 1,
          bgcolor: 'background.paper',
          pb: 1,
        }}
      >
        <TextField
          fullWidth
          label="Find a chore"
          value={query}
          onChange={(event) => setQuery(event.target.value)}
        />
      </Box>
      <Box sx={{ minHeight: reservedRows * ROW_HEIGHT }}>
        {matches.length === 0 ? (
          <Typography color="text.secondary" sx={{ mt: 1 }}>
            No chores match &quot;{query.trim()}&quot;.
          </Typography>
        ) : (
          <List disablePadding aria-label="Chores you can mark done">
            {visible.map((chore) => (
              <ListItem key={chore.id} disablePadding>
                <ListItemButton
                  sx={{ minHeight: ROW_HEIGHT }}
                  selected={chore.id === selectedId}
                  aria-pressed={chore.id === selectedId}
                  onClick={() => onSelect(chore.id)}
                >
                  <ListItemText primary={chore.name} />
                </ListItemButton>
              </ListItem>
            ))}
          </List>
        )}
      </Box>
      <Box sx={{ mt: hasMoreLine ? 1 : 0, minHeight: hasMoreLine ? 24 : 0 }}>
        {hiddenCount > 0 && (
          <Typography color="text.secondary" variant="body2">
            {hiddenCount} more
          </Typography>
        )}
      </Box>
    </>
  )
}
