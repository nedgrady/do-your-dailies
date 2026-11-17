import { Button, Typography } from '@mui/material'
import { format } from 'date-fns'
import ResponsiveDialog from '../../../components/ResponsiveDialog'

interface CompletedChoreDialogProps {
  choreName: string
  // Empty when the completion has no named completors.
  completorNames: string
  completedAt: Date
  onClose: () => void
  onMarkNotDone: () => void
  // Server state for the undo request, owned by the parent's mutation.
  isSaving: boolean
  saveFailed: boolean
}

export default function CompletedChoreDialog({
  choreName,
  completorNames,
  completedAt,
  onClose,
  onMarkNotDone,
  isSaving,
  saveFailed,
}: CompletedChoreDialogProps) {
  return (
    <ResponsiveDialog
      open
      onClose={onClose}
      title={choreName}
      actions={
        <>
          <Button size="large" onClick={onClose}>
            Cancel
          </Button>
          <Button
            size="large"
            variant="contained"
            color="error"
            disabled={isSaving}
            onClick={onMarkNotDone}
          >
            Mark not done
          </Button>
        </>
      }
    >
      <Typography>
        Done {format(completedAt, 'p')}
        {completorNames ? ` by ${completorNames}` : ''}
      </Typography>
      {saveFailed && (
        <Typography color="error" role="alert" sx={{ mt: 1 }}>
          Couldn&apos;t undo, try again.
        </Typography>
      )}
    </ResponsiveDialog>
  )
}
