import TaskAltIcon from '@mui/icons-material/TaskAlt'
import { Box, Button, Grid, List } from '@mui/material'
import { QueryErrorResetBoundary } from '@tanstack/react-query'
import { Suspense } from 'react'
import { ErrorBoundary } from 'react-error-boundary'
import { nullChore, useAllChoresQuery } from '../../../domain/chore'
import { useChoreCompletionsBetweenQuery } from '../../../domain/choreCompletion'
import { useChoreCompletorsQuery } from '../../../domain/choreCompletor'
import { useCurrentDayChoreQueueQuery } from '../../../domain/choreQueue'
import type { ChoreInQueue } from '../../../domain/choreQueue'
import {
  useCompletingChore,
  usePickingChore,
  useViewingCompletion,
} from '../../../lib/useCompletingChore'
import ChoreQueueColumn from './ChoreQueueColumn'
import ChoreTodoListItem from './ChoreTodoListItem'
import CompleteChoreFlow from './CompleteChoreFlow'
import CompletionListItem from './CompletionListItem'
import UndoCompletionFlow from './UndoCompletionFlow'

function ChoreQueueContent() {
  const { data: allChoresLookup = [] } = useAllChoresQuery()
  const { data: choresInQueue = [], dataUpdatedAt } =
    useCurrentDayChoreQueueQuery()
  // dataUpdatedAt changes on every poll, so the due labels refresh with it.
  const now = new Date(dataUpdatedAt)
  const { data: completions = [] } = useChoreCompletionsBetweenQuery(
    new Date(new Date().setHours(0, 0, 0, 0)),
    new Date(new Date().setHours(24, 0, 0, 0)),
  )

  const { data: completors = [] } = useChoreCompletorsQuery()
  const { completingChoreId, open, close } = useCompletingChore()
  const picking = usePickingChore()
  const viewing = useViewingCompletion()

  const toDoChores = choresInQueue.filter(
    (choreInQueue) =>
      !completions.some(
        (completion) => completion.choreId === choreInQueue.choreId,
      ),
  )

  // A stale, unknown or already-completed id in the URL just shows no sheet.
  const presetChoreId =
    completingChoreId !== undefined &&
    allChoresLookup.some((chore) => chore.id === completingChoreId) &&
    !completions.some((completion) => completion.choreId === completingChoreId)
      ? completingChoreId
      : undefined

  // The same flow serves both entry points: a tapped to-do chore (preset) or
  // the "Mark something else done" button (no preset).
  const isFlowOpen = presetChoreId !== undefined || picking.isPicking
  const closeFlow = presetChoreId !== undefined ? close : picking.close

  const completedEntries = completions.map((completion) => ({
    completion,
    chore:
      allChoresLookup.find((chore) => chore.id === completion.choreId) ??
      nullChore,
    names: completion.choreCompletorIds
      .map((id) => completors.find((completor) => completor.id === id)?.name)
      .filter((name): name is string => name !== undefined)
      .join(', '),
  }))

  // A stale, invalid or already-undone id in the URL just shows no sheet.
  const completionBeingViewed =
    completedEntries.find(
      ({ completion }) => completion.id === viewing.viewingCompletionId,
    ) ?? null

  const handleChoreClick = (choreInQueue: ChoreInQueue) => {
    open(choreInQueue.choreId)
  }

  return (
    <Box>
      <Grid container spacing={2}>
        <Grid size={{ xs: 12, md: 6 }}>
          <ChoreQueueColumn
            title="To do"
            count={toDoChores.length}
            emptyMessage="Nothing to do right now"
            footer={
              <Button
                variant="text"
                color="inherit"
                startIcon={<TaskAltIcon />}
                onClick={() => picking.open()}
              >
                Mark something else done
              </Button>
            }
          >
            <List aria-label="Chores to do" disablePadding>
              {toDoChores.map((choreInQueue) => (
                <ChoreTodoListItem
                  key={choreInQueue.choreId}
                  choreInQueue={choreInQueue}
                  now={now}
                  onClick={() => handleChoreClick(choreInQueue)}
                />
              ))}
            </List>
          </ChoreQueueColumn>
        </Grid>

        <Grid size={{ xs: 12, md: 6 }}>
          <ChoreQueueColumn
            title="Done"
            count={completedEntries.length}
            emptyMessage="Nothing done yet today"
          >
            <List aria-label="Done chores" disablePadding>
              {completedEntries.map(({ completion, chore, names }) => (
                <CompletionListItem
                  key={completion.id}
                  choreName={chore.name}
                  emoji={chore.emoji}
                  completorNames={names}
                  onClick={() => viewing.open(completion.id)}
                />
              ))}
            </List>
          </ChoreQueueColumn>
        </Grid>
      </Grid>

      {completionBeingViewed && (
        <UndoCompletionFlow
          key={completionBeingViewed.completion.id}
          completionId={completionBeingViewed.completion.id}
          choreName={completionBeingViewed.chore.name}
          completorNames={completionBeingViewed.names}
          completedAt={completionBeingViewed.completion.createdAt}
          onClose={viewing.close}
        />
      )}

      {isFlowOpen && (
        <CompleteChoreFlow
          key={presetChoreId ?? 'picking'}
          presetChoreId={presetChoreId}
          onClose={closeFlow}
        />
      )}
    </Box>
  )
}

export default function ChoreQueuePage() {
  return (
    <>
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
