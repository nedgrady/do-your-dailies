import type { Theme } from '@mui/material'
import { DataGrid, useGridApiRef } from '@mui/x-data-grid'
import type {
  DataGridProps,
  GridEventListener,
  GridValidRowModel,
} from '@mui/x-data-grid'

// Styling that makes editable cells look clickable, and highlights whichever
// cell is currently being edited.
const editableCellSx = {
  '& .MuiDataGrid-cell--editable': {
    cursor: 'text',
  },
  '& .MuiDataGrid-cell--editable:hover': {
    backgroundColor: 'action.hover',
  },
  '& .MuiDataGrid-cell--editing': {
    backgroundColor: 'action.selected',
    outline: (t: Theme) => `2px solid ${t.palette.primary.main}`,
    outlineOffset: -2,
  },
  '& .MuiDataGrid-cell--editing:focus-within': {
    outline: (t: Theme) => `2px solid ${t.palette.primary.main}`,
  },
} as const

/**
 * A DataGrid that starts editing a cell on a single click, instead of
 * requiring a double click, and that visually distinguishes editable and
 * actively-edited cells. Pass an `onCellClick` as usual to also react to
 * clicks yourself — it runs before the edit-mode logic here.
 */
export function EditableDataGrid<TRow extends GridValidRowModel>({
  onCellClick,
  sx,
  apiRef: apiRefProp,
  editMode,
  ...props
}: DataGridProps<TRow>) {
  const ownApiRef = useGridApiRef()
  const apiRef = apiRefProp ?? ownApiRef

  const handleCellClick: GridEventListener<'cellClick'> = (
    params,
    event,
    details,
  ) => {
    onCellClick?.(params, event, details)

    if (!params.isEditable || params.cellMode === 'edit') return

    // Ignore clicks on portaled elements (e.g. a select menu) that bubble
    // up from outside the cell's own DOM subtree.
    if (
      event.target instanceof Element &&
      event.currentTarget instanceof Element &&
      !event.currentTarget.contains(event.target)
    ) {
      return
    }

    // In "row" edit mode all of a row's cells edit together (and only
    // commit once you leave the row), so start the whole row rather than
    // just this cell.
    if (editMode === 'row') {
      apiRef.current?.startRowEditMode({
        id: params.id,
        fieldToFocus: params.field,
      })
    } else {
      apiRef.current?.startCellEditMode({ id: params.id, field: params.field })
    }
  }

  return (
    <DataGrid
      apiRef={apiRef}
      editMode={editMode}
      onCellClick={handleCellClick}
      sx={[editableCellSx, ...(Array.isArray(sx) ? sx : sx ? [sx] : [])]}
      {...props}
    />
  )
}
