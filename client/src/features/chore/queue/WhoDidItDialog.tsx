import AddIcon from '@mui/icons-material/Add'
import CloseIcon from '@mui/icons-material/Close'
import {
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
import ResponsiveDialog from '../../../components/ResponsiveDialog'
import type { ChoreCompletor } from '../../../domain/choreCompletor'

interface WhoDidItDialogProps {
  choreName: string
  completors: ChoreCompletor[]
  onClose: () => void
  onConfirm: (choreCompletorIds: number[]) => void
  // Resolves with the new completor's id. Rejects on failure.
  onAddPerson: (name: string) => Promise<number>
  // Server state for the two saves, owned by the parent's mutations.
  isSaving: boolean
  saveFailed: boolean
  isAddingPerson: boolean
  addPersonFailed: boolean
}

// Local state here is UI only: which people are ticked and whether the add
// row is open. Saving and failure come from the parent's mutations.
export default function WhoDidItDialog({
  choreName,
  completors,
  onClose,
  onConfirm,
  onAddPerson,
  isSaving,
  saveFailed,
  isAddingPerson,
  addPersonFailed,
}: WhoDidItDialogProps) {
  const [selectedIds, setSelectedIds] = useState<number[]>([])
  const [isAddRowOpen, setIsAddRowOpen] = useState(false)
  const [newName, setNewName] = useState('')

  const toggle = (id: number) => {
    setSelectedIds((current) =>
      current.includes(id)
        ? current.filter((selected) => selected !== id)
        : [...current, id],
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
    setSelectedIds((current) => [...current, id])
    collapseAddPerson()
  }

  return (
    <ResponsiveDialog
      open
      onClose={onClose}
      title={`Who did "${choreName}"?`}
      actions={
        <>
          <Button size="large" onClick={onClose}>
            Cancel
          </Button>
          <Button
            size="large"
            variant="contained"
            disabled={selectedIds.length === 0 || isSaving}
            onClick={() => onConfirm(selectedIds)}
          >
            Done
          </Button>
        </>
      }
    >
      {saveFailed && (
        <Typography color="error" role="alert" sx={{ mb: 1 }}>
          Couldn&apos;t save, try again.
        </Typography>
      )}
      <List disablePadding>
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

        {!isAddRowOpen && (
          <ListItem disablePadding>
            <ListItemButton onClick={() => setIsAddRowOpen(true)}>
              <ListItemIcon>
                <AddIcon />
              </ListItemIcon>
              <ListItemText primary="Add person" />
            </ListItemButton>
          </ListItem>
        )}
      </List>

      {isAddRowOpen && (
        <>
          <Stack direction="row" spacing={1} sx={{ mt: 1, alignItems: 'center' }}>
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
            <IconButton aria-label="Cancel adding person" onClick={collapseAddPerson}>
              <CloseIcon />
            </IconButton>
          </Stack>
          {addPersonFailed && (
            <Typography color="error" role="alert" sx={{ mt: 1 }}>
              Couldn&apos;t add that name, try again.
            </Typography>
          )}
        </>
      )}
    </ResponsiveDialog>
  )
}
