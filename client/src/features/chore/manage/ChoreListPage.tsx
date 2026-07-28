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
  Typography,
  useMediaQuery,
  useTheme,
} from '@mui/material'
import { DataGrid, type GridColDef } from '@mui/x-data-grid'
import { QueryErrorResetBoundary } from '@tanstack/react-query'
import { Suspense, useState } from 'react'
import { ErrorBoundary } from 'react-error-boundary'
import {
  type Chore,
  nullChore,
  useAllChoresQuery,
  useCreateChoreMutation,
  useDeleteChoreMutation,
  useUpdateChoreMutation,
} from '../../../domain/chore'
import { useChoreInsightsQuery } from '../../../domain/insights'

function AddChoreForm() {
  const [name, setName] = useState('')
  const [cadenceInDays, setCadenceInDays] = useState('')
  const createMutation = useCreateChoreMutation()

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const cadence = Number(cadenceInDays)
    if (!name.trim() || !cadence || cadence <= 0) return
    createMutation.mutate(
      { data: { name: name.trim(), cadenceInDays: cadence } },
      {
        onSuccess: () => {
          setName('')
          setCadenceInDays('')
        },
      },
    )
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
  chore: Chore
  onClose: () => void
}) {
  const [name, setName] = useState('')
  const [cadenceInDays, setCadenceInDays] = useState('')
  const isOpen = chore.id !== nullChore.id

  // sync local state whenever a new chore is opened
  if (isOpen && name === '' && cadenceInDays === '') {
    setName(chore.name)
    setCadenceInDays(String(chore.cadenceInDays))
  }

  const updateMutation = useUpdateChoreMutation()
  const deleteMutation = useDeleteChoreMutation()

  const handleClose = () => {
    setName('')
    setCadenceInDays('')
    onClose()
  }

  const handleSave = () => {
    const cadence = Number(cadenceInDays)
    if (!name.trim() || !cadence || cadence <= 0) return
    updateMutation.mutate(
      { id: chore.id, data: { name: name.trim(), cadenceInDays: cadence } },
      { onSuccess: handleClose },
    )
  }

  const handleDelete = () => {
    deleteMutation.mutate({ id: chore.id }, { onSuccess: handleClose })
  }

  return (
    <Dialog open={isOpen} onClose={handleClose} fullWidth maxWidth="xs">
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
          onClick={handleDelete}
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

function ChoreInsights() {
  const { data } = useChoreInsightsQuery()

  if (!data) return <>Loading...</>

  const {
    minCapacityToKeepUtilizationRatioGreaterThanOne,
    userDesiredCapacity,
    utilizationRatio,
  } = data

  const willBeSippage = utilizationRatio > 1

  return (
    <Box>
      <Typography>
        You want to do {userDesiredCapacity} chores per day
      </Typography>
      <Typography>
        Your chore list would require{' '}
        {minCapacityToKeepUtilizationRatioGreaterThanOne.toFixed(1)} chores per
        day
      </Typography>
      {willBeSippage && (
        <>
          <Typography>This means there will be slippage.</Typography>
          <Typography>Example:</Typography>
        </>
      )}
      {!willBeSippage && (
        <Typography>
          This means you will complete all your chores in good time.
        </Typography>
      )}
    </Box>
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
      <ChoreInsights />
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
