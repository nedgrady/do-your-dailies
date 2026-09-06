import { z } from 'zod'
import { useListChoreQueueSuspense } from '../generated/api/default/default'
import { ListChoreQueueResponseItem } from '../generated/zod/default/default'
import { useSimulatedUser } from '../lib/simulatedUser'
import { isoDateTime } from './dateCoercion'
import { queryKeys } from './queryKeys'

const choreInQueueSchema = ListChoreQueueResponseItem.extend({
  lastCompletedAt: isoDateTime.optional(),
})

export type ChoreInQueue = z.infer<typeof choreInQueueSchema>

export function useCurrentDayChoreQueueQuery() {
  const { userId } = useSimulatedUser()
  return useListChoreQueueSuspense({
    query: {
      queryKey: queryKeys.chores.queue(userId),
      select: (data) => z.array(choreInQueueSchema).parse(data),
    },
  })
}
