import AddIcon from '@mui/icons-material/Add'
import {
  Box,
  Button,
  Grid,
  Stack,
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
import ChoreFlow from './ChoreFlow'
import { ChoreInsights } from './ChoreInsights'

const NEW_CHORE_ID = -2

const routeApi = getRouteApi('/chores/manage')

function resolveSelectedChore(
  choreParam: number | 'new' | undefined,
  chores: Chore[],
): Chore {
  if (choreParam === undefined) return nullChore
  if (choreParam === 'new') return draftChore
  return chores.find((c) => c.id === choreParam) ?? nullChore
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
          <Stack
            direction="row"
            spacing={0.5}
            sx={{ alignItems: 'center', height: '100%' }}
          >
            <AddIcon fontSize="small" />
            <Typography component="span" variant="body2">
              Add a chore…
            </Typography>
          </Stack>
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
        <Stack direction="row" spacing={1}>
          <Button component={Link} to="/chores/preview" size="small">
            Preview upcoming days
          </Button>
          <Button component={Link} to="/chores/import" size="small">
            Bulk import
          </Button>
        </Stack>
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
            color: 'primary.main',
          },
        }}
        onRowClick={
          isDesktop
            ? undefined
            : (params) => {
                const row = params.row as Chore
                void navigate({
                  search: { chore: row.id === NEW_CHORE_ID ? 'new' : row.id },
                })
              }
        }
        onCellClick={isDesktop ? handleCellClick : undefined}
        onRowEditStop={() => setEditing((x) => x - 1)}
        disableRowSelectionOnClick
        processRowUpdate={processRowUpdate}
      />

      {selectedChore.id !== nullChore.id && (
        <ChoreFlow
          key={selectedChore.id}
          chore={selectedChore}
          isCreating={choreParam === 'new'}
          onClose={() => void navigate({ search: {}, replace: true })}
        />
      )}
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
