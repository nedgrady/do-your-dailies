import { useQueryClient } from '@tanstack/react-query'
import { z } from 'zod'
import {
  useCreateChoreCompletion,
  useDeleteChoreCompletion,
  useListChoreCompletionsSuspense,
} from '../generated/api/default/default'
import { ListChoreCompletionsResponseItem } from '../generated/zod/default/default'
import { useSimulatedUser } from '../lib/simulatedUser'
import { isoDateTime } from './dateCoercion'
import { queryKeys } from './queryKeys'

const choreCompletionSchema = ListChoreCompletionsResponseItem.extend({
  createdAt: isoDateTime,
  updatedAt: isoDateTime,
})

export type ChoreCompletion = z.infer<typeof choreCompletionSchema>

export function useChoreCompletionsBetweenQuery(start: Date, end: Date) {
  const { userId } = useSimulatedUser()
  return useListChoreCompletionsSuspense(
    { start: start.toISOString(), end: end.toISOString() },
    {
      query: {
        queryKey: queryKeys.chores.completions(userId, start, end),
        select: (data) => z.array(choreCompletionSchema).parse(data),
      },
    },
  )
}

export function useCreateChoreCompletionMutation() {
  const queryClient = useQueryClient()
  const { userId } = useSimulatedUser()

  return useCreateChoreCompletion({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: queryKeys.chores.all(userId),
        })
      },
    },
  })
}

// Undoing a completion changes the queue, the done list and the insights, so
// everything under the chores key is refetched.
export function useDeleteChoreCompletionMutation() {
  const queryClient = useQueryClient()
  const { userId } = useSimulatedUser()

  return useDeleteChoreCompletion({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: queryKeys.chores.all(userId),
        })
      },
    },
  })
}
