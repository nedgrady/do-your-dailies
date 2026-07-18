import {
  useMutation,
  useQueryClient,
  useSuspenseQuery,
} from '@tanstack/react-query'
import z from 'zod'
import axios from '../lib/axios'

export type ChoreCompletion = {
  id: number
  choreId: number
  createdAt: Date
  updatedAt: Date
}

const apiChoreCompletionSchema = z.object({
  id: z.number().int().nonnegative(),
  choreId: z.number().int().nonnegative(),
  createdAt: z.iso.datetime(),
  updatedAt: z.iso.datetime(),
})

const apiChoreCompletionsSchema = z.array(apiChoreCompletionSchema)

function mapApiChoreCompletionToDomain(
  apiChoreCompletion: z.infer<typeof apiChoreCompletionSchema>,
): ChoreCompletion {
  return {
    id: apiChoreCompletion.id,
    choreId: apiChoreCompletion.choreId,
    createdAt: new Date(apiChoreCompletion.createdAt),
    updatedAt: new Date(apiChoreCompletion.updatedAt),
  }
}
async function fetchChoreCompletionsBetween(
  start: Date,
  end: Date,
): Promise<ChoreCompletion[]> {
  const response = await axios.get('/api/chore-completions', {
    params: { start, end },
  })
  const apiChoreCompletions = apiChoreCompletionsSchema.parse(response.data)

  return apiChoreCompletions.map(mapApiChoreCompletionToDomain)
}
export async function createChoreCompletion(choreId: number): Promise<void> {
  await axios.post('/api/chore-completions', { choreId })
}

export function useChoreCompletionsBetweenQuery(start: Date, end: Date) {
  return useSuspenseQuery({
    queryKey: ['chores', 'completions', 'by-date', start, end],
    queryFn: () => fetchChoreCompletionsBetween(start, end),
  })
}

export function useCreateChoreCompletionMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: createChoreCompletion,
    onSuccess: async () => {
      // TODO - invalidate completed query
      await queryClient.invalidateQueries({
        queryKey: ['chores'],
      })
    },
  })
}
