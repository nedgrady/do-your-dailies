import axios from '#/lib/axios'
import {
  QueryErrorResetBoundary,
  useMutation,
  useQueryClient,
  useSuspenseQuery,
} from '@tanstack/react-query'
import { Suspense, useMemo, useState } from 'react'
import { ErrorBoundary } from 'react-error-boundary'
import { z } from 'zod'

const apiChoreSchema = z.object({
  id: z.number().int().nonnegative(),
  name: z.string(),
  cadenceInDays: z.number().int().positive(),
  createdAt: z.iso.datetime(),
  updatedAt: z.iso.datetime(),
})

const apiChoreInQueueSchema = z.object({
  choreId: z.number().int().nonnegative(),
  choreName: z.string(),
  cadenceInDays: z.number().int().positive(),
  priority: z.number(),
  latestCompletionId: z.number().int().nonnegative().optional(),
  lastCompletedAt: z.iso.datetime().optional(),
})

const apiChoresSchema = z.array(apiChoreSchema)
const apiChoreQueueSchema = z.array(apiChoreInQueueSchema)

const apiChoreCompletionSchema = z.object({
  id: z.number().int().nonnegative(),
  choreId: z.number().int().nonnegative(),
  createdAt: z.iso.datetime(),
  updatedAt: z.iso.datetime(),
})

const apiChoreCompletionsSchema = z.array(apiChoreCompletionSchema)

type Chore = {
  id: number
  name: string
  cadenceInDays: number
  createdAt: Date
  updatedAt: Date
}

type ChoreCompletion = {
  id: number
  choreId: number
  createdAt: Date
  updatedAt: Date
}

type ChoreInQueue = {
  cadenceInDays: number
  choreId: number
  choreName: string
  lastCompletedAt?: Date
  latestCompletionId?: number
  priority: number
}

const utcDateFormatter = new Intl.DateTimeFormat('sv-SE', {
  timeZone: 'UTC',
})

function mapApiChoreToDomain(apiChore: z.infer<typeof apiChoreSchema>): Chore {
  return {
    id: apiChore.id,
    name: apiChore.name,
    cadenceInDays: apiChore.cadenceInDays,
    createdAt: new Date(apiChore.createdAt),
    updatedAt: new Date(apiChore.updatedAt),
  }
}

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

function formatUTCDateOnly(input: Date): string {
  return utcDateFormatter.format(input)
}

async function fetchCurrentDayChoreQueue(): Promise<ChoreInQueue[]> {
  const response = await axios.get('/api/chore-queue')
  const apiChores = apiChoreQueueSchema.parse(response.data)

  return apiChores.map(mapApiChoreInQueueToDomain)
}

async function fetchAllChores(): Promise<Chore[]> {
  const response = await axios.get('/api/chores')
  const apiChores = apiChoresSchema.parse(response.data)

  return apiChores.map(mapApiChoreToDomain)
}

async function fetchChoreCompletionsForDate(
  date: string,
): Promise<ChoreCompletion[]> {
  const response = await axios.get('/api/chore-completions', {
    params: { date },
  })
  const apiChoreCompletions = apiChoreCompletionsSchema.parse(response.data)

  return apiChoreCompletions.map(mapApiChoreCompletionToDomain)
}

async function createChoreCompletion(choreId: number): Promise<void> {
  await axios.post('/api/chore-completions', { choreId })
}

function useCurrentDayChoreQueueQuery() {
  return useSuspenseQuery({
    queryKey: ['chore-queue', 'current-day'],
    queryFn: fetchCurrentDayChoreQueue,
  })
}

function useAllChoresQuery() {
  return useSuspenseQuery({
    queryKey: ['chores', 'all'],
    queryFn: fetchAllChores,
  })
}

function useChoreCompletionsForDateQuery(date: string) {
  return useSuspenseQuery({
    queryKey: ['chore-completions', 'by-date', date],
    queryFn: () => fetchChoreCompletionsForDate(date),
  })
}

function ChoreQueueContent() {
  const today = formatUTCDateOnly(new Date())
  const queryClient = useQueryClient()
  const { data: chores } = useCurrentDayChoreQueueQuery()
  const { data: allChores } = useAllChoresQuery()
  const { data: completions } = useChoreCompletionsForDateQuery(today)
  const [completionErrorsByChoreID, setCompletionErrorsByChoreID] = useState<
    Record<number, string>
  >({})

  const createMutation = useMutation({
    mutationFn: createChoreCompletion,
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: ['chore-completions', 'by-date', today],
      })
    },
  })

  const completedChoreIDs = useMemo(() => {
    const ids = new Set<number>()
    completions.forEach((completion) => {
      ids.add(completion.choreId)
    })

    return ids
  }, [completions])

  const displayedChores = useMemo(() => {
    const choresByID = new Map<number, Chore>()
    allChores.forEach((chore) => {
      choresByID.set(chore.id, chore)
    })

    const queueChoreIDs = new Set<number>(chores.map((chore) => chore.choreId))
    const completedNotInQueue = completions
      .filter((completion) => !queueChoreIDs.has(completion.choreId))
      .map((completion) => choresByID.get(completion.choreId))
      .filter((chore): chore is Chore => chore !== undefined)

    return [...chores, ...completedNotInQueue]
  }, [allChores, chores, completions])

  const isDayOff = displayedChores.length === 0
  const isAllDone =
    chores.length > 0 &&
    chores.every((chore) => completedChoreIDs.has(chore.choreId))

  async function handleMarkDone(choreID: number) {
    setCompletionErrorsByChoreID((current) => ({ ...current, [choreID]: '' }))

    try {
      await createMutation.mutateAsync(choreID)
    } catch {
      setCompletionErrorsByChoreID((current) => ({
        ...current,
        [choreID]: 'Could not mark this chore as done.',
      }))
    }
  }

  return (
    <main>
      <section aria-labelledby="chore-queue-heading">
        <h2 id="chore-queue-heading">Today's chore queue</h2>
        {isDayOff ? (
          <p role="status" aria-label="Day off today">
            You&apos;re all caught up. Enjoy your day off.
          </p>
        ) : null}
        {isAllDone ? (
          <p role="status" aria-label="All chores completed today">
            Great work. All chores for today are done.
          </p>
        ) : null}
        <ul aria-label="Current day chore queue">
          <p>TODO</p>
          {/* {displayedChores.map((chore) => (
            <li key={chore.id} className="chore-queue-item">
              <button
                type="button"
                className={`chore-queue-button ${completedChoreIDs.has(chore.id) ? 'is-complete' : ''}`}
                aria-pressed={completedChoreIDs.has(chore.id)}
                aria-label={
                  completedChoreIDs.has(chore.id)
                    ? `${chore.name} completed`
                    : `Mark ${chore.name} as done`
                }
                onClick={() => {
                  if (!completedChoreIDs.has(chore.id)) {
                    void handleMarkDone(chore.id)
                  }
                }}
              >
                <span>{chore.name}</span>
              </button>
              {completionErrorsByChoreID[chore.id] ? (
                <p role="alert">{completionErrorsByChoreID[chore.id]}</p>
              ) : null}
            </li>
          ))} */}
        </ul>
      </section>
    </main>
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
