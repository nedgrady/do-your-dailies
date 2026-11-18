import { Typography } from '@mui/material'
import { z } from 'zod'
import type { FlowDefinition } from '../../../components/flow/flowTypes'
import type { Chore } from '../../../domain/chore'
import { useDeleteChoreMutation } from '../../../domain/chore'

// No values are collected: this flow only confirms a delete.
const deleteSchema = z.object({})
export type DeleteValues = z.infer<typeof deleteSchema>

interface UseDeleteChoreFlowOptions {
  chore: Chore
  onDone: () => void
}

// Describes "delete this chore?" as a flow. A single confirm step with a
// destructive submit.
export function useDeleteChoreFlow({
  chore,
  onDone,
}: UseDeleteChoreFlowOptions): FlowDefinition<DeleteValues> {
  const deleteChore = useDeleteChoreMutation()

  return {
    title: `Delete "${chore.name}"?`,
    schema: deleteSchema,
    defaultValues: {},
    steps: [
      {
        id: 'confirm',
        title: 'Confirm',
        fields: [],
        render: () => (
          <Typography>This removes the chore from your list.</Typography>
        ),
      },
    ],
    submitLabel: 'Delete chore',
    submitColor: 'error',
    submitFailedMessage: "Couldn't delete, try again.",
    onSubmit: () => {
      deleteChore.mutate({ id: chore.id }, { onSuccess: onDone })
    },
    isSubmitting: deleteChore.isPending,
    submitFailed: deleteChore.isError,
  }
}
