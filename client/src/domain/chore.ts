import { useSuspenseQuery } from '@tanstack/react-query'
import z from 'zod'
import axios from '../lib/axios'

export const nullChore: Chore = {
  id: -1,
  name: '?',
  cadenceInDays: 0,
  createdAt: new Date(),
  updatedAt: new Date(),
}

export type Chore = {
  id: number
  name: string
  cadenceInDays: number
  createdAt: Date
  updatedAt: Date
}

const apiChoreSchema = z.object({
  id: z.number().int().nonnegative(),
  name: z.string(),
  cadenceInDays: z.number().int().positive(),
  createdAt: z.iso.datetime(),
  updatedAt: z.iso.datetime(),
})

const apiChoresSchema = z.array(apiChoreSchema)

function mapApiChoreToDomain(apiChore: z.infer<typeof apiChoreSchema>): Chore {
  return {
    id: apiChore.id,
    name: apiChore.name,
    cadenceInDays: apiChore.cadenceInDays,
    createdAt: new Date(apiChore.createdAt),
    updatedAt: new Date(apiChore.updatedAt),
  }
}

async function fetchAllChores(): Promise<Chore[]> {
  const response = await axios.get('/api/chores')
  const apiChores = apiChoresSchema.parse(response.data)

  return apiChores.map(mapApiChoreToDomain)
}

export function useAllChoresQuery() {
  return useSuspenseQuery({
    queryKey: ['chores', 'all'],
    queryFn: fetchAllChores,
  })
}
