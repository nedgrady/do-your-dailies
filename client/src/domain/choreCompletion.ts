import { useQueryClient } from '@tanstack/react-query'
import { z } from 'zod'
import {
  useCreateChoreCompletion,
  useListChoreCompletionsSuspense,
} from '../generated/api/default/default'
import { ListChoreCompletionsResponseItem } from '../generated/zod/default/default'
import { isoDateTime } from './dateCoercion'
import { queryKeys } from './queryKeys'

const choreCompletionSchema = ListChoreCompletionsResponseItem.extend({
  createdAt: isoDateTime,
  updatedAt: isoDateTime,
})

export type ChoreCompletion = z.infer<typeof choreCompletionSchema>

export function useChoreCompletionsBetweenQuery(start: Date, end: Date) {
  return useListChoreCompletionsSuspense(
    { start: start.toISOString(), end: end.toISOString() },
    {
      query: {
        queryKey: queryKeys.chores.completions(start, end),
        select: (data) => z.array(choreCompletionSchema).parse(data),
      },
    },
  )
}

export function useCreateChoreCompletionMutation() {
  const queryClient = useQueryClient()

  return useCreateChoreCompletion({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: queryKeys.chores.all() })
      },
    },
  })
}
