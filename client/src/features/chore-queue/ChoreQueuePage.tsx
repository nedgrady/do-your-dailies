import { QueryErrorResetBoundary } from '@tanstack/react-query'
import { sub } from 'date-fns'
import { Suspense } from 'react'
import { ErrorBoundary } from 'react-error-boundary'
import { nullChore, useAllChoresQuery, type Chore } from '../../domain/chore'
import { useChoreCompletionsBetweenQuery } from '../../domain/choreCompletion'
import { useCurrentDayChoreQueueQuery } from '../../domain/choreQueue'

function ChoreQueueContent() {
  // TODO: Bit lame to get all the chores just for the name but oh well
  const { data: allChoresLookup } = useAllChoresQuery()

  const { data: choresInQueue } = useCurrentDayChoreQueueQuery()
  const { data: completions } = useChoreCompletionsBetweenQuery(
    sub(new Date(), { days: 1 }),
    new Date(),
  )

  const toDoChores = choresInQueue.filter(
    (choreInQueue) =>
      !completions.some(
        (completion) => completion.choreId === choreInQueue.choreId,
      ),
  )

  const completedChores: Chore[] = completions.map(
    (completion) =>
      allChoresLookup.find((chore) => chore.id === completion.choreId) ??
      nullChore,
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
