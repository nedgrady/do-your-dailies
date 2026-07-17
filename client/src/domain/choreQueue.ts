import { useSuspenseQuery } from '@tanstack/react-query'
import z from 'zod'
import axios from '../lib/axios'

export type ChoreInQueue = {
  cadenceInDays: number
  choreId: number
  choreName: string
  lastCompletedAt?: Date
  latestCompletionId?: number
  priority: number
}

const apiChoreInQueueSchema = z.object({
  choreId: z.number().int().nonnegative(),
  choreName: z.string(),
  cadenceInDays: z.number().int().positive(),
  priority: z.number(),
  latestCompletionId: z.number().int().nonnegative().optional(),
  lastCompletedAt: z.iso.datetime().optional(),
})
const apiChoreQueueSchema = z.array(apiChoreInQueueSchema)

function mapApiChoreInQueueToDomain(
  apiChoreInQueue: z.infer<typeof apiChoreInQueueSchema>,
): ChoreInQueue {
  return {
    cadenceInDays: apiChoreInQueue.cadenceInDays,
    choreId: apiChoreInQueue.choreId,
    choreName: apiChoreInQueue.choreName,
    lastCompletedAt: apiChoreInQueue.lastCompletedAt
      ? new Date(apiChoreInQueue.lastCompletedAt)
      : undefined,
    latestCompletionId: apiChoreInQueue.latestCompletionId,
    priority: apiChoreInQueue.priority,
  }
}
async function fetchCurrentDayChoreQueue(): Promise<ChoreInQueue[]> {
  const response = await axios.get('/api/chore-queue')
  const apiChores = apiChoreQueueSchema.parse(response.data)

  return apiChores.map(mapApiChoreInQueueToDomain)
}

export function useCurrentDayChoreQueueQuery() {
  return useSuspenseQuery({
    queryKey: ['chore-queue', 'current-day'],
    queryFn: fetchCurrentDayChoreQueue,
  })
}
