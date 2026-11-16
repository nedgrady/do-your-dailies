import type { z } from 'zod'
import { useGetChorePreviewSuspense } from '../generated/api/default/default'
import { GetChorePreviewResponse } from '../generated/zod/default/default'
import { useSimulatedUser } from '../lib/simulatedUser'
import { queryKeys } from './queryKeys'

export type ChorePreview = z.infer<typeof GetChorePreviewResponse>

export function useChorePreviewQuery(days: number) {
  const { userId } = useSimulatedUser()
  return useGetChorePreviewSuspense(
    { days },
    {
      query: {
        queryKey: queryKeys.chores.preview(userId, days),
        select: (data) => GetChorePreviewResponse.parse(data),
      },
    },
  )
}
