import { useQueryClient } from '@tanstack/react-query'
import { z } from 'zod'
import {
  useBulkCreateChores,
  useCreateChore,
  useDeleteChore,
  useListChoresSuspense,
  useUpdateChore,
} from '../generated/api/default/default'
import { ListChoresResponseItem } from '../generated/zod/default/default'
import { useSimulatedUser } from '../lib/simulatedUser'
import { isoDateTime } from './dateCoercion'
import { queryKeys } from './queryKeys'

const choreSchema = ListChoresResponseItem.extend({
  createdAt: isoDateTime,
  updatedAt: isoDateTime,
})

export type Chore = z.infer<typeof choreSchema>

export const nullChore: Chore = {
  id: -1,
  name: '?',
  cadenceInDays: 0,
  displayUnit: 'DAY',
  createdAt: new Date(0),
  updatedAt: new Date(0),
}

export function useAllChoresQuery() {
  const { userId } = useSimulatedUser()
  return useListChoresSuspense({
    query: {
      queryKey: queryKeys.chores.list(userId),
      select: (data) => z.array(choreSchema).parse(data),
    },
  })
}

export function useCreateChoreMutation() {
  const queryClient = useQueryClient()
  const { userId } = useSimulatedUser()
  return useCreateChore({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: queryKeys.chores.all(userId),
        })
      },
    },
  })
}

export function useBulkCreateChoresMutation() {
  const queryClient = useQueryClient()
  const { userId } = useSimulatedUser()
  return useBulkCreateChores({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: queryKeys.chores.all(userId),
        })
      },
    },
  })
}

export function useUpdateChoreMutation() {
  const queryClient = useQueryClient()
  const { userId } = useSimulatedUser()
  return useUpdateChore({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: queryKeys.chores.all(userId),
        })
      },
    },
  })
}

export function useDeleteChoreMutation() {
  const queryClient = useQueryClient()
  const { userId } = useSimulatedUser()
  return useDeleteChore({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: queryKeys.chores.all(userId),
        })
      },
    },
  })
}
