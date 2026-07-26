import DeleteIcon from '@mui/icons-material/Delete'
import {
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  Stack,
  TextField,
  useMediaQuery,
  useTheme,
} from '@mui/material'
import { DataGrid, type GridColDef } from '@mui/x-data-grid'
import {
  QueryErrorResetBoundary,
  useMutation,
  useQueryClient,
} from '@tanstack/react-query'
import { Suspense, useState } from 'react'
import { ErrorBoundary } from 'react-error-boundary'
import { nullChore, useAllChoresQuery } from '../../../domain/chore'

type Chore = {
  id: number
  name: string
  cadenceInDays: number
}

const CHORES_QUERY_KEY = ['chores']

async function updateChore(
  id: number,
  data: { name?: string; cadenceInDays?: number },
): Promise<Chore> {
  const res = await fetch(`/api/chores/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!res.ok) throw new Error('Failed to update chore')
  return res.json()
}

async function deleteChore(id: number): Promise<void> {
  const res = await fetch(`/api/chores/${id}`, { method: 'DELETE' })
  if (!res.ok) throw new Error('Failed to delete chore')
}

async function createChore(data: {
  name: string
  cadenceInDays: number
}): Promise<Chore> {
  const res = await fetch('/api/chores', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!res.ok) throw new Error('Failed to create chore')
  return res.json()
}

function AddChoreForm() {
  const [name, setName] = useState('')
  const [cadenceInDays, setCadenceInDays] = useState('')
  const queryClient = useQueryClient()

  const createMutation = useMutation({
    mutationFn: createChore,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CHORES_QUERY_KEY })
      setName('')
      setCadenceInDays('')
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const cadence = Number(cadenceInDays)
    if (!name.trim() || !cadence || cadence <= 0) return
    createMutation.mutate({ name: name.trim(), cadenceInDays: cadence })
  }

  return (
    <Box component="form" onSubmit={handleSubmit} sx={{ mb: 2 }}>
      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
        <TextField
          label="Chore name"
          fullWidth
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <TextField
          label="Cadence (days)"
          type="number"
          sx={{ minWidth: 140 }}
          value={cadenceInDays}
          onChange={(e) => setCadenceInDays(e.target.value)}
        />
        <Button
          type="submit"
          variant="contained"
          size="large"
          disabled={createMutation.isPending}
        >
          Add
        </Button>
      </Stack>
    </Box>
  )
}

function EditChoreDialog({
  chore,
  onClose,
}: {
  chore: Chore | null
  onClose: () => void
}) {
  const [name, setName] = useState('')
  const [cadenceInDays, setCadenceInDays] = useState('')
  const queryClient = useQueryClient()

  // sync local state whenever a new chore is opened
  if (chore && name === '' && cadenceInDays === '') {
    setName(chore.name)
    setCadenceInDays(String(chore.cadenceInDays))
  }

  const updateMutation = useMutation({
    mutationFn: (data: { name?: string; cadenceInDays?: number }) =>
      updateChore(chore!.id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CHORES_QUERY_KEY })
      handleClose()
    },
  })

  const deleteMutation = useMutation({
    mutationFn: () => deleteChore(chore!.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CHORES_QUERY_KEY })
      handleClose()
    },
  })

  const handleClose = () => {
    setName('')
    setCadenceInDays('')
    onClose()
  }

  const handleSave = () => {
    const cadence = Number(cadenceInDays)
    if (!name.trim() || !cadence || cadence <= 0) return
    updateMutation.mutate({ name: name.trim(), cadenceInDays: cadence })
  }

  return (
    <Dialog
      open={chore !== nullChore}
      onClose={handleClose}
      fullWidth
      maxWidth="xs"
    >
      <DialogTitle>Edit chore</DialogTitle>
      <DialogContent>
        <Stack spacing={2} sx={{ mt: 1 }}>
          <TextField
            label="Chore name"
            fullWidth
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          <TextField
            label="Cadence (days)"
            type="number"
            value={cadenceInDays}
            onChange={(e) => setCadenceInDays(e.target.value)}
          />
        </Stack>
      </DialogContent>
      <DialogActions sx={{ px: 3, pb: 2 }}>
        <IconButton
          aria-label="delete chore"
          color="error"
          onClick={() => deleteMutation.mutate()}
        >
          <DeleteIcon />
        </IconButton>
        <Box sx={{ flex: 1 }} />
        <Button onClick={handleClose}>Cancel</Button>
        <Button
          variant="contained"
          onClick={handleSave}
          disabled={updateMutation.isPending}
        >
          Save
        </Button>
      </DialogActions>
    </Dialog>
  )
}

function ChoreList() {
  const theme = useTheme()
  const isDesktop = useMediaQuery(theme.breakpoints.up('md'))
  const { data: chores } = useAllChoresQuery()
  const [selectedChore, setSelectedChore] = useState<Chore>(nullChore)

  const columns: GridColDef[] = [
    { field: 'name', headerName: 'Chore', flex: 1, editable: isDesktop },
    {
      field: 'cadenceInDays',
      headerName: 'Cadence (days)',
      width: 160,
      editable: isDesktop,
    },
  ]

  return (
    <Box sx={{ mt: 2, mx: 'auto' }}>
      <AddChoreForm />
      <p>{isDesktop}</p>
      <DataGrid
        rows={chores}
        columns={columns}
        showToolbar
        slotProps={{ toolbar: { showQuickFilter: true } }}
        rowHeight={isDesktop ? 40 : 80}
        onRowClick={
          isDesktop
            ? undefined
            : (params) => setSelectedChore(params.row as Chore)
        }
        disableRowSelectionOnClick
      />
      <EditChoreDialog
        chore={selectedChore}
        onClose={() => setSelectedChore(nullChore)}
      />
    </Box>
  )
}

export default function ChoreListPage() {
  return (
    <>
      <QueryErrorResetBoundary>
        {({ reset }) => (
          <ErrorBoundary
            onReset={reset}
            fallbackRender={({ resetErrorBoundary }) => (
              <>
                <p>Could not load chore queue.</p>
                <button type="button" onClick={resetErrorBoundary}>
                  Retry
                </button>
              </>
            )}
          >
            <Suspense fallback={<p>Loading chores...</p>}>
              <ChoreList />
            </Suspense>
          </ErrorBoundary>
        )}
      </QueryErrorResetBoundary>
    </>
  )
}
