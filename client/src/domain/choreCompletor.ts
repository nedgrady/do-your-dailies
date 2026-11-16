import { useQueryClient } from '@tanstack/react-query'
import { z } from 'zod'
import {
  useCreateChoreCompletor,
  useListChoreCompletorsSuspense,
} from '../generated/api/default/default'
import { ListChoreCompletorsResponseItem } from '../generated/zod/default/default'
import { useSimulatedUser } from '../lib/simulatedUser'
import { isoDateTime } from './dateCoercion'
import { queryKeys } from './queryKeys'

const choreCompletorSchema = ListChoreCompletorsResponseItem.extend({
  createdAt: isoDateTime,
  updatedAt: isoDateTime,
})

export type ChoreCompletor = z.infer<typeof choreCompletorSchema>

export function useChoreCompletorsQuery() {
  const { userId } = useSimulatedUser()
  return useListChoreCompletorsSuspense({
    query: {
      queryKey: queryKeys.chores.completors(userId),
      select: (data) => z.array(choreCompletorSchema).parse(data),
    },
  })
}

export function useCreateChoreCompletorMutation() {
  const queryClient = useQueryClient()
  const { userId } = useSimulatedUser()

  return useCreateChoreCompletor({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: queryKeys.chores.completors(userId),
        })
      },
    },
  })
}
