import DeleteIcon from '@mui/icons-material/Delete'
import {
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Grid,
  IconButton,
  Stack,
  TextField,
  Typography,
  useMediaQuery,
  useTheme,
} from '@mui/material'
import type {
  GridCellParams,
  GridColDef,
  GridPreProcessEditCellProps,
  GridRowModel,
  GridRowModesModel,
} from '@mui/x-data-grid'
import { DataGrid, useGridApiRef } from '@mui/x-data-grid'
import { QueryErrorResetBoundary } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'
import { Suspense, useState } from 'react'
import { ErrorBoundary } from 'react-error-boundary'
import type { Chore } from '../../../domain/chore'
import {
  nullChore,
  useAllChoresQuery,
  useCreateChoreMutation,
  useDeleteChoreMutation,
  useUpdateChoreMutation,
} from '../../../domain/chore'
import {
  mutationStatus,
  SaveStatusIndicator,
} from '../../../forms/SaveStatusIndicator'
import { UnsavedChangesGuard } from '../../../forms/UnsavedChangesGuard'
import { ChoreInsights } from './ChoreInsights'

const validateName = ({ props }: GridPreProcessEditCellProps) => {
  const value = String(props.value ?? '').trim()
  return {
    ...props,
    error: value.length === 0 ? 'Name is required' : undefined,
  }
}

const validateCadence = ({ props }: GridPreProcessEditCellProps) => {
  const value = Number(props.value)
  return {
    ...props,
    error: !value || value <= 0 ? 'Must be a positive number' : undefined,
  }
}

function AddChoreForm() {
  const [name, setName] = useState('')
  const [cadenceInDays, setCadenceInDays] = useState('')
  const createMutation = useCreateChoreMutation()

  const handleSubmit = (e: React.SubmitEvent) => {
    // e.preventDefault()
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

function ChoreList() {
  const theme = useTheme()
  const isDesktop = useMediaQuery(theme.breakpoints.up('md'))
  const { data: chores } = useAllChoresQuery()
  const [selectedChore, setSelectedChore] = useState<Chore>(nullChore)
  const updateMutation = useUpdateChoreMutation()

  const [rowModesModel, setRowModesModel] = useState<GridRowModesModel>({})
  const [editing, setEditing] = useState(0)
  const apiRef = useGridApiRef()

  const handleCellClick = (
    params: GridCellParams<Chore>,
    event: React.MouseEvent,
  ) => {
    if (!params.isEditable || params.cellMode === 'edit') return

    // Ignore clicks on portaled elements (e.g. a select menu) that bubble
    // up from outside the cell's own DOM subtree.
    if (
      event.target instanceof Element &&
      !event.currentTarget.contains(event.target)
    ) {
      return
    }

    setEditing((x) => x + 1)
    apiRef.current?.startCellEditMode({ id: params.id, field: params.field })
  }

  const processRowUpdate = async (
    newRow: GridRowModel<Chore>,
    oldRow: GridRowModel<Chore>,
  ): Promise<GridRowModel<Chore>> => {
    const name = String(newRow.name ?? '').trim()
    const cadenceInDays = Number(newRow.cadenceInDays)

    if (!name || !cadenceInDays || cadenceInDays <= 0) {
      // Shouldn't normally hit this — preProcessEditCellProps below stops
      // Save from running while a field is invalid — but guard anyway.
      throw new Error('Chore name and a positive cadence are required')
    }

    if (name === oldRow.name && cadenceInDays === oldRow.cadenceInDays) {
      return oldRow
    }

    await updateMutation.mutateAsync({
      id: oldRow.id,
      data: { name, cadenceInDays },
    })

    // Only name/cadenceInDays are editable. id/createdAt/updatedAt always
    // come from oldRow, never from newRow or the mutation response.
    return { ...oldRow, name, cadenceInDays }
  }

  const columns: GridColDef<Chore>[] = [
    {
      field: 'name',
      headerName: 'Chore',
      flex: 1,
      editable: isDesktop,
      preProcessEditCellProps: validateName,
    },
    {
      field: 'cadenceInDays',
      headerName: 'Cadence (days)',
      width: 160,
      type: 'number',
      editable: isDesktop,
      preProcessEditCellProps: validateCadence,
    },
  ]

  return (
    <Box>
      <UnsavedChangesGuard hasUnsavedChanges={editing > 0} />
      <AddChoreForm />
      <Stack direction="row" sx={{ mb: 1 }}>
        <Button component={Link} to="/chores/import" size="small">
          Bulk import (hello beauty)
        </Button>
      </Stack>
      <Grid container sx={{ alignItems: 'stretch', mb: 0.5 }}>
        <Grid>
          <ChoreInsights />
        </Grid>

        <Grid
          size="grow"
          sx={{
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'flex-end',
            alignItems: 'flex-end',
            pb: 0.5,
          }}
        >
          <SaveStatusIndicator status={mutationStatus(updateMutation)} />
        </Grid>
      </Grid>
      {isDesktop && (
        <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
          Click a cell to edit it &mdash; press Enter to save, Esc to cancel.
        </Typography>
      )}
      <DataGrid
        apiRef={apiRef}
        rows={chores}
        columns={columns}
        showToolbar
        rowHeight={isDesktop ? 40 : 80}
        rowModesModel={rowModesModel}
        onRowModesModelChange={setRowModesModel}
        onRowClick={
          isDesktop
            ? undefined
            : (params) => setSelectedChore(params.row as Chore)
        }
        onCellClick={isDesktop ? handleCellClick : undefined}
        onCellEditStop={() => setEditing((x) => x - 1)}
        disableRowSelectionOnClick
        processRowUpdate={processRowUpdate}
        sx={{
          '& .MuiDataGrid-cell--editable': {
            cursor: 'text',
          },
          '& .MuiDataGrid-cell--editable:hover': {
            backgroundColor: 'action.hover',
          },
          '& .MuiDataGrid-cell--editing': {
            backgroundColor: 'action.selected',
            outline: (t) => `2px solid ${t.palette.primary.main}`,
            outlineOffset: -2,
          },
          '& .MuiDataGrid-cell--editing:focus-within': {
            outline: (t) => `2px solid ${t.palette.primary.main}`,
          },
        }}
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
