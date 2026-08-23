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
  ToggleButton,
  ToggleButtonGroup,
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
import { DISPLAY_UNITS, toAmount, toDays, unitLabel } from '../../Cadence'
import {
  mutationStatus,
  SaveStatusIndicator,
} from '../../../forms/SaveStatusIndicator'
import { UnsavedChangesGuard } from '../../../forms/UnsavedChangesGuard'
import { EditableDataGrid } from '../../../integrations/mui/EditableDataGrid'
import type { DisplayUnit } from '../../../generated/api/api.schemas'
import { ChoreInsights } from './ChoreInsights'

const validateName = ({ props }: GridPreProcessEditCellProps) => {
  const value = String(props.value ?? '').trim()
  return {
    ...props,
    error: value.length === 0 ? 'Name is required' : undefined,
  }
}

const validateAmount = ({ props }: GridPreProcessEditCellProps) => {
  const value = Number(props.value)
  return {
    ...props,
    error:
      !Number.isInteger(value) || value <= 0
        ? 'Must be a positive whole number'
        : undefined,
  }
}

function AmountUnitInput({
  amount,
  onAmountChange,
  unit,
  onUnitChange,
}: {
  amount: string
  onAmountChange: (value: string) => void
  unit: DisplayUnit
  onUnitChange: (unit: DisplayUnit) => void
}) {
  return (
    <>
      <Typography>every</Typography>
      <TextField
        label="Amount"
        type="number"
        sx={{ minWidth: 90 }}
        value={amount}
        onChange={(e) => onAmountChange(e.target.value)}
        slotProps={{ htmlInput: { min: 1, step: 1 } }}
      />
      <ToggleButtonGroup
        value={unit}
        exclusive
        onChange={(_e, value: DisplayUnit | null) => {
          if (value) onUnitChange(value)
        }}
        aria-label="cadence unit"
        size="small"
      >
        {DISPLAY_UNITS.map((u) => (
          <ToggleButton key={u} value={u} aria-label={unitLabel(u)}>
            {unitLabel(u)}
          </ToggleButton>
        ))}
      </ToggleButtonGroup>
    </>
  )
}

function AddChoreForm() {
  const [name, setName] = useState('')
  const [amount, setAmount] = useState('1')
  const [unit, setUnit] = useState<DisplayUnit>('DAY')
  const createMutation = useCreateChoreMutation()

  const parsedAmount = Number(amount)
  const isAmountValid = Number.isInteger(parsedAmount) && parsedAmount > 0

  const handleSubmit = (e: React.SubmitEvent) => {
    e.preventDefault()
    if (!name.trim() || !isAmountValid) return
    createMutation.mutate(
      {
        data: {
          name: name.trim(),
          cadenceInDays: toDays(parsedAmount, unit),
          displayUnit: unit,
        },
      },
      {
        onSuccess: () => {
          setName('')
          setAmount('1')
          setUnit('DAY')
        },
      },
    )
  }

  return (
    <Box component="form" onSubmit={handleSubmit} sx={{ mb: 2 }}>
      <Stack
        direction={{ xs: 'column', sm: 'row' }}
        spacing={2}
        sx={{ alignItems: { sm: 'center' } }}
      >
        <TextField
          label="Chore name"
          fullWidth
          value={name}
          onChange={(e) => setName(e.target.value)}
        />

        <AmountUnitInput
          amount={amount}
          onAmountChange={setAmount}
          unit={unit}
          onUnitChange={setUnit}
        />

        <Button
          type="submit"
          variant="contained"
          size="large"
          disabled={createMutation.isPending || !name.trim() || !isAmountValid}
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
  const [amount, setAmount] = useState('')
  const [unit, setUnit] = useState<DisplayUnit>('DAY')
  const isOpen = chore.id !== nullChore.id

  // sync local state whenever a new chore is opened
  if (isOpen && name === '' && amount === '') {
    setName(chore.name)
    setAmount(String(toAmount(chore.cadenceInDays, chore.displayUnit)))
    setUnit(chore.displayUnit)
  }

  const updateMutation = useUpdateChoreMutation()
  const deleteMutation = useDeleteChoreMutation()

  const handleClose = () => {
    setName('')
    setAmount('')
    setUnit('DAY')
    onClose()
  }

  const handleSave = () => {
    const parsedAmount = Number(amount)
    if (!name.trim() || !Number.isInteger(parsedAmount) || parsedAmount <= 0) {
      return
    }
    updateMutation.mutate(
      {
        id: chore.id,
        data: {
          name: name.trim(),
          cadenceInDays: toDays(parsedAmount, unit),
          displayUnit: unit,
        },
      },
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
          <Stack direction="row" spacing={2} sx={{ alignItems: 'center' }}>
            <AmountUnitInput
              amount={amount}
              onAmountChange={setAmount}
              unit={unit}
              onUnitChange={setUnit}
            />
          </Stack>
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

  const handleCellClick = (params: GridCellParams<Chore>) => {
    if (params.isEditable && params.cellMode !== 'edit') {
      setEditing((x) => x + 1)
    }
  }

  const processRowUpdate = async (
    newRow: GridRowModel<Chore>,
    oldRow: GridRowModel<Chore>,
  ): Promise<GridRowModel<Chore>> => {
    const name = String(newRow.name ?? '').trim()
    const cadenceInDays = Number(newRow.cadenceInDays)
    const displayUnit = newRow.displayUnit

    if (!name || !cadenceInDays || cadenceInDays <= 0) {
      // Shouldn't normally hit this — preProcessEditCellProps below stops
      // Save from running while a field is invalid — but guard anyway.
      throw new Error('Chore name and a positive cadence are required')
    }

    if (
      name === oldRow.name &&
      cadenceInDays === oldRow.cadenceInDays &&
      displayUnit === oldRow.displayUnit
    ) {
      return oldRow
    }

    await updateMutation.mutateAsync({
      id: oldRow.id,
      data: { name, cadenceInDays, displayUnit },
    })

    // Only name/cadenceInDays/displayUnit are editable. id/createdAt/updatedAt
    // always come from oldRow, never from newRow or the mutation response.
    return { ...oldRow, name, cadenceInDays, displayUnit }
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
      field: 'cadenceAmount',
      headerName: 'Cadence',
      width: 110,
      type: 'number',
      editable: isDesktop,
      valueGetter: (_value, row) =>
        toAmount(row.cadenceInDays, row.displayUnit),
      valueSetter: (value, row) => ({
        ...row,
        cadenceInDays: toDays(Number(value), row.displayUnit),
      }),
      preProcessEditCellProps: validateAmount,
    },
    {
      field: 'displayUnit',
      headerName: 'Unit',
      width: 130,
      type: 'singleSelect',
      valueOptions: DISPLAY_UNITS.map((u) => ({
        value: u,
        label: unitLabel(u),
      })),
      editable: isDesktop,
      valueSetter: (value, row) => {
        const amount = toAmount(row.cadenceInDays, row.displayUnit)
        const newUnit = value as DisplayUnit
        return {
          ...row,
          displayUnit: newUnit,
          cadenceInDays: toDays(amount, newUnit),
        }
      },
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
      <EditableDataGrid
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
