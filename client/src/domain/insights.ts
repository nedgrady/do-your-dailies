import { useGetChoreInsights } from '../generated/api/default/default'
import { useSimulatedUser } from '../lib/simulatedUser'
import { queryKeys } from './queryKeys'

// const choreInsightsSchema = GetChoreInsightsResponse.extend({})

// export type ChoreInsights = z.infer<typeof choreInsightsSchema>

export function useChoreInsightsQuery() {
  const { userId } = useSimulatedUser()
  return useGetChoreInsights({
    query: {
      queryKey: queryKeys.chores.insights(userId),
    },
  })
}
