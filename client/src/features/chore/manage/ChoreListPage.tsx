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
import { getRouteApi, Link } from '@tanstack/react-router'
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
import type { DisplayUnit } from '../../../generated/api/api.schemas'
import { EditableDataGrid } from '../../../integrations/mui/EditableDataGrid'
import { DISPLAY_UNITS, toAmount, toDays, unitLabel } from '../../Cadence'
import { ChoreInsights } from './ChoreInsights'

const NEW_CHORE_ID = -2

const routeApi = getRouteApi('/chores/manage')

function resolveSelectedChore(
  choreParam: string | undefined,
  chores: Chore[],
): Chore {
  if (choreParam === undefined) return nullChore
  if (choreParam === 'new') return draftChore
  return chores.find((c) => c.id === Number(choreParam)) ?? nullChore
}

const draftChore: Chore = {
  id: NEW_CHORE_ID,
  name: '',
  cadenceInDays: toDays(1, 'DAY'),
  displayUnit: 'DAY',
  createdAt: new Date(0),
  updatedAt: new Date(0),
}

const compareStrings = (a: unknown, b: unknown) =>
  String(a).localeCompare(String(b))

const compareNumbers = (a: unknown, b: unknown) => Number(a) - Number(b)

// Column sort direction doesn't auto-negate our comparator when
// getSortComparator is used (unlike a plain sortComparator), so we handle
// 'desc' ourselves — which lets us keep the draft row pinned first
// regardless of the chosen direction, not just flipped with everything else.
const withDraftFirst =
  (
    compare: (a: unknown, b: unknown) => number,
  ): NonNullable<GridColDef<Chore>['getSortComparator']> =>
  (sortDirection) =>
  (v1, v2, params1, params2) => {
    const isDraft1 = params1.id === NEW_CHORE_ID
    const isDraft2 = params2.id === NEW_CHORE_ID
    if (isDraft1 || isDraft2) {
      return isDraft1 === isDraft2 ? 0 : isDraft1 ? -1 : 1
    }
    const result = compare(v1, v2)
    return sortDirection === 'desc' ? -result : result
  }

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
    <Grid container spacing={2} sx={{ flexGrow: 1, alignItems: 'center' }}>
      <Grid size={{ xs: 12, sm: 'auto' }}>
        <Stack direction="row" spacing={1} sx={{ alignItems: 'center' }}>
          <Typography>every</Typography>
          <TextField
            label="Amount"
            type="number"
            sx={{ minWidth: 90 }}
            value={amount}
            onChange={(e) => onAmountChange(e.target.value)}
            slotProps={{ htmlInput: { min: 1, step: 1 } }}
          />
        </Stack>
      </Grid>
      <Grid size={{ xs: 12, sm: 'grow' }}>
        <ToggleButtonGroup
          value={unit}
          exclusive
          onChange={(_e, value: DisplayUnit | null) => {
            if (value) onUnitChange(value)
          }}
          aria-label="cadence unit"
          fullWidth
        >
          {DISPLAY_UNITS.map((u) => (
            <ToggleButton key={u} value={u} aria-label={unitLabel(u)}>
              {unitLabel(u)}
            </ToggleButton>
          ))}
        </ToggleButtonGroup>
      </Grid>
    </Grid>
  )
}

function EditChoreDialog({
  chore,
  onClose,
  isDesktop,
}: {
  chore: Chore
  onClose: () => void
  isDesktop: boolean
}) {
  const [name, setName] = useState('')
  const [amount, setAmount] = useState('')
  const [unit, setUnit] = useState<DisplayUnit>('DAY')
  const isOpen = chore.id !== nullChore.id
  const isCreating = chore.id === NEW_CHORE_ID

  // sync local state whenever a new chore is opened
  if (isOpen && name === '' && amount === '') {
    setName(chore.name)
    setAmount(String(toAmount(chore.cadenceInDays, chore.displayUnit)))
    setUnit(chore.displayUnit)
  }

  const createMutation = useCreateChoreMutation()
  const updateMutation = useUpdateChoreMutation()
  const deleteMutation = useDeleteChoreMutation()
  const saveMutation = isCreating ? createMutation : updateMutation

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
    const data = {
      name: name.trim(),
      cadenceInDays: toDays(parsedAmount, unit),
      displayUnit: unit,
    }
    if (isCreating) {
      createMutation.mutate({ data }, { onSuccess: handleClose })
    } else {
      updateMutation.mutate({ id: chore.id, data }, { onSuccess: handleClose })
    }
  }

  const handleDelete = () => {
    deleteMutation.mutate({ id: chore.id }, { onSuccess: handleClose })
  }

  return (
    <Dialog
      open={isOpen}
      onClose={handleClose}
      fullWidth
      maxWidth="sm"
      fullScreen={!isDesktop}
    >
      <DialogTitle>{isCreating ? 'Add chore' : 'Edit chore'}</DialogTitle>
      <DialogContent>
        <Grid container spacing={4} sx={{ mt: 1 }}>
          <Grid size={12}>
            <TextField
              label="Chore name"
              fullWidth
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </Grid>
          <Grid size={12}>
            <AmountUnitInput
              amount={amount}
              onAmountChange={setAmount}
              unit={unit}
              onUnitChange={setUnit}
            />
          </Grid>
        </Grid>
      </DialogContent>
      <DialogActions sx={{ px: 3, pb: 2 }}>
        {!isCreating && (
          <IconButton
            aria-label="delete chore"
            color="error"
            onClick={handleDelete}
            size="large"
          >
            <DeleteIcon />
          </IconButton>
        )}
        <Box sx={{ flex: 1 }} />
        <Button onClick={handleClose} size="large">
          Cancel
        </Button>
        <Button
          variant="contained"
          onClick={handleSave}
          disabled={saveMutation.isPending}
          size="large"
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
  const { chore: choreParam } = routeApi.useSearch()
  const navigate = routeApi.useNavigate()
  const selectedChore = resolveSelectedChore(choreParam, chores)
  const createMutation = useCreateChoreMutation()
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

    if (oldRow.id === NEW_CHORE_ID) {
      // Let partial edits accumulate across separate single-cell commits
      // (e.g. cadence set before name) instead of resetting on each one.
      if (!name || !cadenceInDays || cadenceInDays <= 0) {
        return newRow
      }
      await createMutation.mutateAsync({
        data: { name, cadenceInDays, displayUnit },
      })
      return draftChore
    }

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
      getSortComparator: withDraftFirst(compareStrings),
      renderCell: (params) =>
        params.id === NEW_CHORE_ID && !params.value ? (
          <Typography
            component="span"
            sx={{ color: 'text.disabled', fontStyle: 'italic' }}
          >
            + Add a chore…
          </Typography>
        ) : (
          params.value
        ),
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
      getSortComparator: withDraftFirst(compareNumbers),
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
      getSortComparator: withDraftFirst(compareStrings),
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
      <Stack
        direction="row"
        sx={{ justifyContent: 'space-between', alignItems: 'center', mb: 2 }}
      >
        <Typography variant="h6" component="h2">
          Manage chores
        </Typography>
        <Button component={Link} to="/chores/import" size="small">
          Bulk import
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
        rows={[draftChore, ...chores]}
        columns={columns}
        showToolbar
        editMode="row"
        rowHeight={isDesktop ? 40 : 80}
        rowModesModel={rowModesModel}
        onRowModesModelChange={setRowModesModel}
        getRowClassName={(params) =>
          params.id === NEW_CHORE_ID ? 'draft-row' : ''
        }
        sx={{
          '& .draft-row .MuiDataGrid-cell:not(.MuiDataGrid-cell--editing)': {
            color: 'text.disabled',
          },
        }}
        onRowClick={
          isDesktop
            ? undefined
            : (params) => {
                const row = params.row as Chore
                void navigate({
                  search: {
                    chore: row.id === NEW_CHORE_ID ? 'new' : String(row.id),
                  },
                })
              }
        }
        onCellClick={isDesktop ? handleCellClick : undefined}
        onRowEditStop={() => setEditing((x) => x - 1)}
        disableRowSelectionOnClick
        processRowUpdate={processRowUpdate}
      />

      <EditChoreDialog
        chore={selectedChore}
        onClose={() => void navigate({ search: {}, replace: true })}
        isDesktop={isDesktop}
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
