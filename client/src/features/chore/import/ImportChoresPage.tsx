import ArrowBackIcon from '@mui/icons-material/ArrowBack'
import DeleteIcon from '@mui/icons-material/Delete'
import {
  Alert,
  AlertTitle,
  Box,
  Button,
  FormControl,
  FormControlLabel,
  FormLabel,
  IconButton,
  Radio,
  RadioGroup,
  Stack,
  TextField,
  ToggleButton,
  ToggleButtonGroup,
  Typography,
} from '@mui/material'
import { alpha } from '@mui/material/styles'
import type { GridCellParams, GridColDef, GridRowModel } from '@mui/x-data-grid'
import { useNavigate } from '@tanstack/react-router'
import type { KeyboardEvent, MouseEvent } from 'react'
import { useState } from 'react'
import { useBulkCreateChoresMutation } from '../../../domain/chore'
import type { DisplayUnit } from '../../../generated/api/api.schemas'
import { EditableDataGrid } from '../../../integrations/mui/EditableDataGrid'
import { DISPLAY_UNITS, toAmount, toDays, unitLabel } from '../../Cadence'
import { parseChoreRows } from './parseChoreRows'

type FieldDelimiterChoice = 'auto' | 'tab' | 'comma' | 'custom'
type RowDelimiterChoice = 'newline' | 'semicolon' | 'custom'

const FIELD_DELIMITER_VALUES: Record<
  Exclude<FieldDelimiterChoice, 'custom'>,
  string
> = {
  auto: '',
  tab: '\t',
  comma: ',',
}

const ROW_DELIMITER_VALUES: Record<
  Exclude<RowDelimiterChoice, 'custom'>,
  string
> = {
  newline: '\n',
  semicolon: ';',
}

interface DraftRow {
  id: number
  name: string
  cadenceInDays: number | null
  displayUnit: DisplayUnit
  cadenceRaw: string
  parseError: string | null
}

// Explains why a row can't be imported yet, or null if it's valid.
function describeRowError(row: DraftRow): string | null {
  if (row.parseError) return row.parseError
  if (row.name.trim().length === 0) return 'Chore name is empty'
  if (row.cadenceInDays === null || row.cadenceInDays <= 0) {
    return row.cadenceRaw
      ? `Couldn't understand cadence "${row.cadenceRaw}"`
      : 'Cadence is missing'
  }
  return null
}

export default function ImportChoresPage() {
  const [pasteText, setPasteText] = useState('')
  const [rows, setRows] = useState<DraftRow[]>([])
  const navigate = useNavigate()
  const bulkCreateMutation = useBulkCreateChoresMutation()

  // Auto-detect handles the common case (pasting straight from a
  // spreadsheet) with zero input, so the picker stays tucked away until
  // someone actively switches into custom mode.
  const [importMode, setImportMode] = useState<'auto' | 'custom'>('auto')

  const [fieldDelimiterChoice, setFieldDelimiterChoice] =
    useState<FieldDelimiterChoice>('auto')
  const [customFieldDelimiter, setCustomFieldDelimiter] = useState('')
  const [rowDelimiterChoice, setRowDelimiterChoice] =
    useState<RowDelimiterChoice>('newline')
  const [customRowDelimiter, setCustomRowDelimiter] = useState('')

  const handleImportModeChange = (
    _e: MouseEvent<HTMLElement>,
    mode: 'auto' | 'custom' | null,
  ) => {
    if (mode === null) return // exclusive group: re-clicking the active button
    setImportMode(mode)
    if (mode === 'auto') {
      // A genuine reset, not just hiding the panel — going back to Auto
      // means the underlying choices are back to where they started too.
      setFieldDelimiterChoice('auto')
      setCustomFieldDelimiter('')
      setRowDelimiterChoice('newline')
      setCustomRowDelimiter('')
    }
  }

  const resolvedFieldDelimiter =
    fieldDelimiterChoice === 'custom'
      ? customFieldDelimiter
      : FIELD_DELIMITER_VALUES[fieldDelimiterChoice]
  // A blank custom row delimiter would otherwise split on every character —
  // fall back to the default rather than let that happen.
  const resolvedRowDelimiter =
    rowDelimiterChoice === 'custom'
      ? customRowDelimiter || '\n'
      : ROW_DELIMITER_VALUES[rowDelimiterChoice]

  // Keeps the placeholder example in sync with whatever's actually
  // selected, so it never demonstrates a format the box won't accept.
  const placeholderFieldSep = resolvedFieldDelimiter || '\t'
  const placeholder = [
    ['Dishes', '1'].join(placeholderFieldSep),
    ['Vacuum', '1', 'week'].join(placeholderFieldSep),
    ['Mow lawn', '2', 'weeks'].join(placeholderFieldSep),
  ].join(resolvedRowDelimiter)

  // Lets Tab insert a literal tab instead of moving focus, since that's
  // otherwise impossible to type — but that alone would be a keyboard trap
  // (WCAG 2.1.2), so Escape arms a one-time exception: the next Tab press
  // moves focus normally, then trapping resumes once the user types again.
  const [tabEscaped, setTabEscaped] = useState(false)

  const handlePasteBoxKeyDown = (e: KeyboardEvent<HTMLDivElement>) => {
    if (e.key === 'Escape') {
      setTabEscaped(true)
      return
    }
    if (e.key === 'Tab') {
      if (tabEscaped) {
        setTabEscaped(false)
        return
      }
      // Environments without execCommand (some non-browser test runners,
      // and it's a deprecated API generally) fall back to normal Tab
      // behavior rather than silently swallowing the keystroke.
      if (typeof document.execCommand !== 'function') return
      e.preventDefault()
      document.execCommand('insertText', false, '\t')
      return
    }
    if (tabEscaped) {
      setTabEscaped(false)
    }
  }

  // A cell mid-edit hasn't been committed to `rows` yet, so its validity
  // isn't reflected there — block Submit while any cell is still being
  // edited rather than let it use stale, already-committed data underneath.
  const [editingCount, setEditingCount] = useState(0)

  const handleCellClick = (params: GridCellParams<DraftRow>) => {
    if (params.isEditable && params.cellMode !== 'edit') {
      setEditingCount((x) => x + 1)
    }
  }

  // Whether the box has been edited since the last Parse — not derivable
  // from pasteText/rows alone (that's an ordering-of-events question), so
  // this is the minimal one-bit note needed: cleared on Parse, set on edit.
  const [isStale, setIsStale] = useState(false)

  const handleParse = () => {
    const parsed = parseChoreRows(pasteText, {
      field: resolvedFieldDelimiter,
      row: resolvedRowDelimiter,
    })
    setRows(parsed.map((row, index) => ({ id: index, ...row })))
    setIsStale(false)
  }

  const processRowUpdate = (
    newRow: GridRowModel<DraftRow>,
  ): GridRowModel<DraftRow> => {
    const name = String(newRow.name).trim()
    const cadenceInDays = Number(newRow.cadenceInDays)
    const updated: DraftRow = {
      id: newRow.id,
      name,
      cadenceInDays:
        Number.isFinite(cadenceInDays) && cadenceInDays > 0
          ? cadenceInDays
          : null,
      displayUnit: newRow.displayUnit,
      // No longer "as pasted" once the row has been edited directly — keeps
      // describeRowError from quoting stale original paste text after a
      // previously-valid row gets edited into a different invalid state.
      cadenceRaw: '',
      // A structural parse error is about the original raw paste — once the
      // user has edited this row's fields directly, it no longer applies.
      parseError: null,
    }
    setRows((prev) =>
      prev.map((row) => (row.id === updated.id ? updated : row)),
    )
    return updated
  }

  const handleDeleteRow = (id: number) => {
    setRows((prev) => prev.filter((row) => row.id !== id))
  }

  const columns: GridColDef<DraftRow>[] = [
    {
      field: 'name',
      headerName: 'Chore',
      flex: 1,
      editable: true,
    },
    {
      field: 'cadenceAmount',
      headerName: 'Cadence',
      width: 110,
      type: 'number',
      editable: true,
      valueGetter: (_value, row) =>
        toAmount(row.cadenceInDays ?? 0, row.displayUnit),
      valueSetter: (value, row) => ({
        ...row,
        cadenceInDays: toDays(Number(value), row.displayUnit),
      }),
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
      editable: true,
      valueSetter: (value, row) => {
        const amount = toAmount(row.cadenceInDays ?? 0, row.displayUnit)
        const newUnit = value as DisplayUnit
        return {
          ...row,
          displayUnit: newUnit,
          cadenceInDays: toDays(amount, newUnit),
        }
      },
    },
    {
      field: 'remove',
      headerName: '',
      width: 56,
      sortable: false,
      filterable: false,
      disableColumnMenu: true,
      renderCell: (params) => (
        <IconButton
          aria-label="remove row"
          size="small"
          onClick={() => handleDeleteRow(params.row.id)}
        >
          <DeleteIcon fontSize="small" />
        </IconButton>
      ),
    },
  ]

  const rowErrors = rows
    .map((row) => ({ row, error: describeRowError(row) }))
    .filter(
      (entry): entry is { row: DraftRow; error: string } =>
        entry.error !== null,
    )

  const allRowsValid = rows.length > 0 && rowErrors.length === 0

  const handleSubmit = () => {
    bulkCreateMutation.mutate(
      {
        data: {
          chores: rows.map((row) => ({
            name: row.name.trim(),
            cadenceInDays: row.cadenceInDays as number,
            displayUnit: row.displayUnit,
          })),
        },
      },
      {
        onSuccess: () => navigate({ to: '/chores/manage' }),
      },
    )
  }

  return (
    <Box>
      <Button
        size="small"
        startIcon={<ArrowBackIcon />}
        onClick={() => navigate({ to: '/chores/manage' })}
        sx={{ mb: 1 }}
      >
        Back to chores
      </Button>
      <Typography variant="h6" sx={{ mb: 1 }}>
        Bulk import chores
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
        Paste a list of chores, one per line. Copying columns out of a
        spreadsheet like Google Sheets or Excel works too.
      </Typography>
      <TextField
        multiline
        minRows={6}
        fullWidth
        placeholder={placeholder}
        value={pasteText}
        onChange={(e) => {
          setPasteText(e.target.value)
          setIsStale(true)
        }}
        onKeyDown={handlePasteBoxKeyDown}
        sx={{ mb: 0.5 }}
        slotProps={{ htmlInput: { 'aria-label': 'Chores to import' } }}
      />
      <Stack
        direction="row"
        sx={{ justifyContent: 'space-between', alignItems: 'center', mb: 2 }}
      >
        <ToggleButtonGroup
          value={importMode}
          exclusive
          size="small"
          onChange={handleImportModeChange}
          aria-label="import format"
        >
          <ToggleButton value="auto">Auto-detect</ToggleButton>
          <ToggleButton value="custom">Custom settings</ToggleButton>
        </ToggleButtonGroup>
        <Button
          // Once rows exist, "Add N chores" below is the primary action —
          // only one contained button should be on screen at a time. If the
          // box is edited afterward without re-parsing, Parse is relevant
          // again, so it reclaims primary until the grid is back in sync.
          variant={rows.length === 0 || isStale ? 'contained' : 'outlined'}
          onClick={handleParse}
          disabled={!pasteText.trim()}
        >
          Parse
        </Button>
      </Stack>

      {importMode === 'custom' && (
        <Stack
          direction={{ xs: 'column', sm: 'row' }}
          spacing={4}
          sx={{ mb: 2 }}
        >
          <FormControl>
            <FormLabel id="field-delimiter-label">Between columns</FormLabel>
            <RadioGroup
              aria-labelledby="field-delimiter-label"
              value={fieldDelimiterChoice}
              onChange={(e) =>
                setFieldDelimiterChoice(e.target.value as FieldDelimiterChoice)
              }
            >
              <FormControlLabel
                value="auto"
                control={<Radio size="small" />}
                label="Auto-detect"
              />
              <FormControlLabel
                value="tab"
                control={<Radio size="small" />}
                label="Tab"
              />
              <FormControlLabel
                value="comma"
                control={<Radio size="small" />}
                label="Comma"
              />
              <FormControlLabel
                value="custom"
                control={<Radio size="small" />}
                label={
                  <TextField
                    size="small"
                    placeholder="Custom"
                    value={customFieldDelimiter}
                    onFocus={() => setFieldDelimiterChoice('custom')}
                    onChange={(e) => {
                      setCustomFieldDelimiter(e.target.value)
                      setFieldDelimiterChoice('custom')
                    }}
                  />
                }
              />
            </RadioGroup>
          </FormControl>
          <FormControl>
            <FormLabel id="row-delimiter-label">Between rows</FormLabel>
            <RadioGroup
              aria-labelledby="row-delimiter-label"
              value={rowDelimiterChoice}
              onChange={(e) =>
                setRowDelimiterChoice(e.target.value as RowDelimiterChoice)
              }
            >
              <FormControlLabel
                value="newline"
                control={<Radio size="small" />}
                label="New line"
              />
              <FormControlLabel
                value="semicolon"
                control={<Radio size="small" />}
                label="Semicolon"
              />
              <FormControlLabel
                value="custom"
                control={<Radio size="small" />}
                label={
                  <TextField
                    size="small"
                    placeholder="Custom"
                    value={customRowDelimiter}
                    onFocus={() => setRowDelimiterChoice('custom')}
                    onChange={(e) => {
                      setCustomRowDelimiter(e.target.value)
                      setRowDelimiterChoice('custom')
                    }}
                  />
                }
              />
            </RadioGroup>
          </FormControl>
        </Stack>
      )}

      {rows.length > 0 && (
        <>
          {rowErrors.length > 0 && (
            <Alert severity="error" role="alert" sx={{ mb: 2 }}>
              <AlertTitle>
                Fix {rowErrors.length} row{rowErrors.length === 1 ? '' : 's'}{' '}
                before importing
              </AlertTitle>
              <Box component="ul" sx={{ m: 0, pl: 2 }}>
                {rowErrors.map(({ row, error }) => (
                  <Typography component="li" variant="body2" key={row.id}>
                    {row.name.trim() || `Row ${row.id + 1}`}: {error}
                  </Typography>
                ))}
              </Box>
            </Alert>
          )}
          <EditableDataGrid
            rows={rows}
            columns={columns}
            processRowUpdate={processRowUpdate}
            disableRowSelectionOnClick
            onCellClick={handleCellClick}
            onCellEditStop={() => setEditingCount((x) => x - 1)}
            getRowClassName={(params) =>
              describeRowError(params.row) !== null ? 'invalid-row' : ''
            }
            sx={{
              '& .invalid-row': {
                bgcolor: (theme) => alpha(theme.palette.error.main, 0.12),
              },
              '& .invalid-row:hover': {
                bgcolor: (theme) => alpha(theme.palette.error.main, 0.2),
              },
            }}
          />
          <Stack direction="row" spacing={2} sx={{ mt: 2 }}>
            <Button
              // Mirrors the Parse button's logic in reverse — only one of
              // the two should be primary at a time.
              variant={isStale ? 'outlined' : 'contained'}
              onClick={handleSubmit}
              disabled={
                !allRowsValid ||
                editingCount > 0 ||
                bulkCreateMutation.isPending
              }
            >
              Add {rows.length} chore{rows.length === 1 ? '' : 's'}
            </Button>
          </Stack>
          {bulkCreateMutation.isError && (
            <Typography color="error" role="alert" sx={{ mt: 1 }}>
              Something went wrong — check the rows and try again.
            </Typography>
          )}
        </>
      )}
    </Box>
  )
}
