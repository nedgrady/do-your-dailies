import { z } from 'zod'
import { useListChoreQueueSuspense } from '../generated/api/default/default'
import { ListChoreQueueResponseItem } from '../generated/zod/default/default'
import { useSimulatedUser } from '../lib/simulatedUser'
import { isoDateTime } from './dateCoercion'
import { emojiOrNoneSchema } from './emoji'
import { queryKeys } from './queryKeys'

const choreInQueueSchema = ListChoreQueueResponseItem.extend({
  lastCompletedAt: isoDateTime.optional(),
  dueAt: isoDateTime,
  emoji: emojiOrNoneSchema,
})

export type ChoreInQueue = z.infer<typeof choreInQueueSchema>

// Polling keeps "time left" labels fresh; each fetch updates dataUpdatedAt,
// which callers use as `now`.
const QUEUE_POLL_INTERVAL_MS = 5 * 60 * 1000

export function useCurrentDayChoreQueueQuery() {
  const { userId } = useSimulatedUser()
  return useListChoreQueueSuspense({
    query: {
      queryKey: queryKeys.chores.queue(userId),
      select: (data) => z.array(choreInQueueSchema).parse(data),
      refetchInterval: QUEUE_POLL_INTERVAL_MS,
    },
  })
}
