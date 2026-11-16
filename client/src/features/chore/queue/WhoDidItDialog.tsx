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
} from '@mui/material'
import { useState } from 'react'
import ResponsiveDialog from '../../../components/ResponsiveDialog'
import type { ChoreCompletor } from '../../../domain/choreCompletor'

interface WhoDidItDialogProps {
  choreName: string
  completors: ChoreCompletor[]
  onClose: () => void
  onConfirm: (choreCompletorIds: number[]) => void
  onAddPerson: (name: string) => Promise<number>
}

// Presentational: the parent mounts a fresh instance per chore, so selection
// always starts empty and nothing is preselected.
export default function WhoDidItDialog({
  choreName,
  completors,
  onClose,
  onConfirm,
  onAddPerson,
}: WhoDidItDialogProps) {
  const [selectedIds, setSelectedIds] = useState<number[]>([])
  const [isAddingPerson, setIsAddingPerson] = useState(false)
  const [newName, setNewName] = useState('')
  const [isAdding, setIsAdding] = useState(false)

  const toggle = (id: number) => {
    setSelectedIds((current) =>
      current.includes(id)
        ? current.filter((selected) => selected !== id)
        : [...current, id],
    )
  }

  const collapseAddPerson = () => {
    setIsAddingPerson(false)
    setNewName('')
  }

  const addPerson = async () => {
    const name = newName.trim()
    if (!name) return
    setIsAdding(true)
    try {
      const id = await onAddPerson(name)
      setSelectedIds((current) => [...current, id])
      collapseAddPerson()
    } finally {
      setIsAdding(false)
    }
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
            disabled={selectedIds.length === 0}
            onClick={() => onConfirm(selectedIds)}
          >
            Done
          </Button>
        </>
      }
    >
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

        {!isAddingPerson && (
          <ListItem disablePadding>
            <ListItemButton onClick={() => setIsAddingPerson(true)}>
              <ListItemIcon>
                <AddIcon />
              </ListItemIcon>
              <ListItemText primary="Add person" />
            </ListItemButton>
          </ListItem>
        )}
      </List>

      {isAddingPerson && (
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
            disabled={!newName.trim() || isAdding}
          >
            Add
          </Button>
          <IconButton aria-label="Cancel adding person" onClick={collapseAddPerson}>
            <CloseIcon />
          </IconButton>
        </Stack>
      )}
    </ResponsiveDialog>
  )
}
