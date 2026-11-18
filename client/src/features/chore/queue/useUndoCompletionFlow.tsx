import { Typography } from '@mui/material'
import { format } from 'date-fns'
import { z } from 'zod'
import type { FlowDefinition } from '../../../components/flow/flowTypes'
import { useDeleteChoreCompletionMutation } from '../../../domain/choreCompletion'

// No values are collected: this flow only confirms an undo.
const undoSchema = z.object({})
export type UndoValues = z.infer<typeof undoSchema>

interface UseUndoCompletionFlowOptions {
  completionId: number
  choreName: string
  // Empty when the completion has no named completors.
  completorNames: string
  completedAt: Date
  onDone: () => void
}

// Describes "mark this completion not done" as a flow. A single summary step
// with a destructive submit.
export function useUndoCompletionFlow({
  completionId,
  choreName,
  completorNames,
  completedAt,
  onDone,
}: UseUndoCompletionFlowOptions): FlowDefinition<UndoValues> {
  const deleteCompletion = useDeleteChoreCompletionMutation()

  return {
    title: choreName,
    schema: undoSchema,
    defaultValues: {},
    steps: [
      {
        id: 'summary',
        title: 'Completed',
        fields: [],
        render: () => (
          <Typography>
            Done {format(completedAt, 'p')}
            {completorNames ? ` by ${completorNames}` : ''}
          </Typography>
        ),
      },
    ],
    submitLabel: 'Mark not done',
    submitColor: 'error',
    submitFailedMessage: "Couldn't undo, try again.",
    onSubmit: () => {
      deleteCompletion.mutate({ id: completionId }, { onSuccess: onDone })
    },
    isSubmitting: deleteCompletion.isPending,
    submitFailed: deleteCompletion.isError,
  }
}
