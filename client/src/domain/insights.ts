import { useGetChoreInsights } from '../generated/api/default/default'
import { queryKeys } from './queryKeys'

// const choreInsightsSchema = GetChoreInsightsResponse.extend({})

// export type ChoreInsights = z.infer<typeof choreInsightsSchema>

export function useChoreInsightsQuery() {
  return useGetChoreInsights({
    query: {
      queryKey: queryKeys.chores.insights(),
    },
  })
}
