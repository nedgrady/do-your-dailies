import DeleteIcon from '@mui/icons-material/Delete'
import {
  Box,
  Button,
  IconButton,
  Stack,
  TextField,
  Typography,
} from '@mui/material'
import type {
  GridColDef,
  GridPreProcessEditCellProps,
  GridRowModel,
} from '@mui/x-data-grid'
import { useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useBulkCreateChoresMutation } from '../../../domain/chore'
import { EditableDataGrid } from '../../../integrations/mui/EditableDataGrid'
import { parseChoreRows } from './parseChoreRows'

interface DraftRow {
  id: number
  name: string
  cadenceInDays: number | null
}

export default function ImportChoresPage() {
  const [pasteText, setPasteText] = useState('')
  const [rows, setRows] = useState<DraftRow[]>([])
  const navigate = useNavigate()
  const bulkCreateMutation = useBulkCreateChoresMutation()

  const handleParse = () => {
    const parsed = parseChoreRows(pasteText)
    setRows(parsed.map((row, index) => ({ id: index, ...row })))
  }

  const processRowUpdate = (
    newRow: GridRowModel<DraftRow>,
  ): GridRowModel<DraftRow> => {
    const name = String(newRow.name).trim()
    const cadence = Number(newRow.cadenceInDays)
    const updated: DraftRow = {
      id: newRow.id,
      name,
      cadenceInDays: Number.isFinite(cadence) && cadence > 0 ? cadence : null,
    }
    setRows((prev) =>
      prev.map((row) => (row.id === updated.id ? updated : row)),
    )
    return updated
  }

  const handleDeleteRow = (id: number) => {
    setRows((prev) => prev.filter((row) => row.id !== id))
  }

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

  const columns: GridColDef<DraftRow>[] = [
    {
      field: 'name',
      headerName: 'Chore',
      flex: 1,
      editable: true,
      preProcessEditCellProps: validateName,
    },
    {
      field: 'cadenceInDays',
      headerName: 'Cadence (days)',
      width: 160,
      type: 'number',
      editable: true,
      preProcessEditCellProps: validateCadence,
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

  const allRowsValid =
    rows.length > 0 &&
    rows.every(
      (row) =>
        row.name.trim().length > 0 &&
        row.cadenceInDays !== null &&
        row.cadenceInDays > 0,
    )

  const handleSubmit = () => {
    bulkCreateMutation.mutate(
      {
        data: {
          chores: rows.map((row) => ({
            name: row.name.trim(),
            cadenceInDays: row.cadenceInDays as number,
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
      <Typography variant="h6" sx={{ mb: 1 }}>
        Bulk import chores
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
        Paste a list of chores, one per line. Copying two columns (name and
        cadence in days) out of a spreadsheet like Google Sheets or Excel works
        too.
      </Typography>
      <TextField
        multiline
        minRows={6}
        fullWidth
        placeholder={'Dishes\t1\nVacuum\t7\nMow lawn\t14'}
        value={pasteText}
        onChange={(e) => setPasteText(e.target.value)}
        sx={{ mb: 1 }}
      />
      <Stack direction="row" justifyContent="flex-end" sx={{ mb: 2 }}>
        <Button
          variant="outlined"
          onClick={handleParse}
          disabled={!pasteText.trim()}
        >
          Parse
        </Button>
      </Stack>

      {rows.length > 0 && (
        <>
          <EditableDataGrid
            rows={rows}
            columns={columns}
            processRowUpdate={processRowUpdate}
            disableRowSelectionOnClick
          />
          <Stack
            direction="row"
            justifyContent="flex-end"
            spacing={2}
            sx={{ mt: 2 }}
          >
            <Button
              variant="contained"
              onClick={handleSubmit}
              disabled={!allRowsValid || bulkCreateMutation.isPending}
            >
              Add {rows.length} chore{rows.length === 1 ? '' : 's'}
            </Button>
          </Stack>
          {bulkCreateMutation.isError && (
            <Typography color="error" sx={{ mt: 1 }}>
              Something went wrong — check the rows and try again.
            </Typography>
          )}
        </>
      )}
    </Box>
  )
}
