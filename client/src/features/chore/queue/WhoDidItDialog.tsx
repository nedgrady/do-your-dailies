import AddIcon from '@mui/icons-material/Add'
import CloseIcon from '@mui/icons-material/Close'
import {
  Button,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  Stack,
  TextField,
} from '@mui/material'
import { useState } from 'react'
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
    <Dialog open onClose={onClose} fullWidth maxWidth="xs">
      <DialogTitle>Who did "{choreName}"?</DialogTitle>
      <DialogContent>
        <Stack
          direction="row"
          spacing={1}
          useFlexGap
          sx={{ flexWrap: 'wrap', alignItems: 'center', mb: 1 }}
        >
          {completors.map((completor) => {
            const selected = selectedIds.includes(completor.id)
            return (
              <Chip
                key={completor.id}
                label={completor.name}
                color={selected ? 'primary' : 'default'}
                variant={selected ? 'filled' : 'outlined'}
                onClick={() => toggle(completor.id)}
                aria-pressed={selected}
              />
            )
          })}

          {isAddingPerson ? (
            <>
              <TextField
                size="small"
                autoFocus
                placeholder="Name"
                slotProps={{ htmlInput: { 'aria-label': 'New person name' } }}
                value={newName}
                onChange={(event) => setNewName(event.target.value)}
                onKeyDown={(event) => {
                  if (event.key === 'Enter') void addPerson()
                  if (event.key === 'Escape') collapseAddPerson()
                }}
                sx={{ width: 160 }}
              />
              <Button
                size="small"
                onClick={() => void addPerson()}
                disabled={!newName.trim() || isAdding}
              >
                Add
              </Button>
              <IconButton
                size="small"
                aria-label="Cancel adding person"
                onClick={collapseAddPerson}
              >
                <CloseIcon fontSize="small" />
              </IconButton>
            </>
          ) : (
            <IconButton
              size="small"
              aria-label="Add person"
              onClick={() => setIsAddingPerson(true)}
            >
              <AddIcon />
            </IconButton>
          )}
        </Stack>
      </DialogContent>
      <DialogActions sx={{ px: 3, pb: 2 }}>
        <Button onClick={onClose}>Cancel</Button>
        <Button
          variant="contained"
          disabled={selectedIds.length === 0}
          onClick={() => onConfirm(selectedIds)}
        >
          Done
        </Button>
      </DialogActions>
    </Dialog>
  )
}
