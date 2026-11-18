import AddIcon from '@mui/icons-material/Add'
import CloseIcon from '@mui/icons-material/Close'
import {
  Box,
  Button,
  Checkbox,
  IconButton,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Stack,
  TextField,
  Typography,
} from '@mui/material'
import { useState } from 'react'
import type { ChoreCompletor } from '../../../domain/choreCompletor'

// Tall enough for the name field row, so swapping "Add person" for the field
// doesn't resize the dialog.
const ADD_SLOT_HEIGHT = 64

interface PeopleStepProps {
  completors: ChoreCompletor[]
  selectedIds: number[]
  onChange: (choreCompletorIds: number[]) => void
  // Resolves with the new completor's id. Rejects on failure.
  onAddPerson: (name: string) => Promise<number>
  isAddingPerson: boolean
  addPersonFailed: boolean
}

// Who did it. Ticked people are the selection; adding a person ticks them too.
export default function PeopleStep({
  completors,
  selectedIds,
  onChange,
  onAddPerson,
  isAddingPerson,
  addPersonFailed,
}: PeopleStepProps) {
  const [isAddRowOpen, setIsAddRowOpen] = useState(false)
  const [newName, setNewName] = useState('')

  const toggle = (id: number) => {
    onChange(
      selectedIds.includes(id)
        ? selectedIds.filter((selected) => selected !== id)
        : [...selectedIds, id],
    )
  }

  const collapseAddPerson = () => {
    setIsAddRowOpen(false)
    setNewName('')
  }

  const addPerson = async () => {
    const name = newName.trim()
    if (!name) return
    const id = await onAddPerson(name).catch(() => null)
    // A failure is shown from addPersonFailed; keep the row open to retry.
    if (id === null) return
    onChange([...selectedIds, id])
    collapseAddPerson()
  }

  return (
    <>
      <List disablePadding aria-label="People who did it">
        {completors.map((completor) => {
          const selected = selectedIds.includes(completor.id)
          return (
            <ListItem key={completor.id} disablePadding>
              <ListItemButton
                onClick={() => toggle(completor.id)}
                aria-pressed={selected}
              >
                <ListItemIcon>
                  <Checkbox
                    checked={selected}
                    tabIndex={-1}
                    disableRipple
                    sx={{ pointerEvents: 'none' }}
                  />
                </ListItemIcon>
                <ListItemText primary={completor.name} />
              </ListItemButton>
            </ListItem>
          )
        })}
      </List>

      {/* One fixed slot: either the "Add person" button or the name field. */}
      <Box
        sx={{
          minHeight: ADD_SLOT_HEIGHT,
          display: 'flex',
          alignItems: 'center',
        }}
      >
        {isAddRowOpen ? (
          <Stack
            direction="row"
            spacing={1}
            sx={{ width: '100%', alignItems: 'center' }}
          >
            <TextField
              autoFocus
              fullWidth
              label="New person's name"
              value={newName}
              onChange={(event) => setNewName(event.target.value)}
              onKeyDown={(event) => {
                if (event.key === 'Enter') void addPerson()
                if (event.key === 'Escape') collapseAddPerson()
              }}
            />
            <Button
              size="large"
              onClick={() => void addPerson()}
              disabled={!newName.trim() || isAddingPerson}
            >
              Add
            </Button>
            <IconButton
              aria-label="Cancel adding person"
              onClick={collapseAddPerson}
            >
              <CloseIcon />
            </IconButton>
          </Stack>
        ) : (
          <ListItemButton
            onClick={() => setIsAddRowOpen(true)}
            sx={{ width: '100%', borderRadius: 1 }}
          >
            <ListItemIcon>
              <AddIcon />
            </ListItemIcon>
            <ListItemText primary="Add person" />
          </ListItemButton>
        )}
      </Box>

      {addPersonFailed && (
        <Typography color="error" role="alert" sx={{ mt: 1 }}>
          Couldn&apos;t add that name, try again.
        </Typography>
      )}
    </>
  )
}
