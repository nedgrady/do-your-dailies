import {
  QueryErrorResetBoundary,
  useSuspenseQuery,
} from '@tanstack/react-query'
import axios from 'axios'
import { Suspense } from 'react'
import { ErrorBoundary } from 'react-error-boundary'
import { z } from 'zod'

const apiChoreSchema = z.object({
  id: z.number().int().nonnegative(),
  name: z.string(),
  cadenceInDays: z.number().int().positive(),
  createdAt: z.iso.datetime(),
  updatedAt: z.iso.datetime(),
})

const apiChoreQueueSchema = z.array(apiChoreSchema)

type Chore = {
  id: number
  name: string
  cadenceInDays: number
  createdAt: Date
  updatedAt: Date
}

function mapApiChoreToDomain(apiChore: z.infer<typeof apiChoreSchema>): Chore {
  return {
    id: apiChore.id,
    name: apiChore.name,
    cadenceInDays: apiChore.cadenceInDays,
    createdAt: new Date(apiChore.createdAt),
    updatedAt: new Date(apiChore.updatedAt),
  }
}

async function fetchCurrentDayChoreQueue(): Promise<Chore[]> {
  const response = await axios.get('/api/chore-queue')
  const apiChores = apiChoreQueueSchema.parse(response.data)

  return apiChores.map(mapApiChoreToDomain)
}

function useCurrentDayChoreQueueQuery() {
  return useSuspenseQuery({
    queryKey: ['chore-queue', 'current-day'],
    queryFn: fetchCurrentDayChoreQueue,
  })
}

function ChoreQueueContent() {
  const { data: chores } = useCurrentDayChoreQueueQuery()

  return (
    <>
      <ul aria-label="Current day chore queue">
        {chores.map((chore) => (
          <li key={chore.id}>{chore.name}</li>
        ))}
      </ul>
      {chores.length === 0 ? <p>No chores due today.</p> : null}
    </>
  )
}

export default function ChoreQueuePage() {
  return (
    <>
      <h1>Required Chores</h1>
      <QueryErrorResetBoundary>
        {({ reset }) => (
          <ErrorBoundary
            onReset={reset}
            fallbackRender={({ resetErrorBoundary }) => (
              <>
                <p>Could not load chore queue.</p>
                <button type="button" onClick={resetErrorBoundary}>
                  Retry
                </button>
              </>
            )}
          >
            <Suspense fallback={<p>Loading chore queue...</p>}>
              <ChoreQueueContent />
            </Suspense>
          </ErrorBoundary>
        )}
      </QueryErrorResetBoundary>
    </>
  )
}
