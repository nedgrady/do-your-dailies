import TaskAltIcon from '@mui/icons-material/TaskAlt'
import {
  Box,
  Button,
  Grid,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  Paper,
  Stack,
  Typography,
} from '@mui/material'
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
import CompleteChoreFlow from './CompleteChoreFlow'
import UndoCompletionFlow from './UndoCompletionFlow'

function ChoreQueueContent() {
  const { data: allChoresLookup = [] } = useAllChoresQuery()
  const { data: choresInQueue = [] } = useCurrentDayChoreQueueQuery()
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
          <Paper variant="outlined" sx={{ p: 2, height: '100%' }}>
            <Stack
              direction="row"
              sx={{
                alignItems: 'center',
                justifyContent: 'space-between',
                mb: 1,
              }}
            >
              <Typography variant="h6" component="h2">
                To do
              </Typography>
              <Button
                startIcon={<TaskAltIcon />}
                onClick={() => picking.open()}
              >
                Mark something else done
              </Button>
            </Stack>
            <List aria-label="Chores to do" disablePadding>
              {toDoChores.map((choreInQueue) => (
                <ListItem
                  key={choreInQueue.choreId}
                  disablePadding
                  sx={{ mb: 1 }}
                >
                  <ListItemButton
                    onClick={() => handleChoreClick(choreInQueue)}
                    component={Paper}
                    variant="outlined"
                    sx={{ p: 2 }}
                  >
                    <ListItemText primary={choreInQueue.choreName} />
                  </ListItemButton>
                </ListItem>
              ))}
            </List>
          </Paper>
        </Grid>

        <Grid size={{ xs: 12, md: 6 }}>
          <Paper variant="outlined" sx={{ p: 2, height: '100%' }}>
            <Typography variant="h6" gutterBottom component="h2">
              Done
            </Typography>
            <List aria-label="Done chores" disablePadding>
              {completedEntries.map(({ completion, chore, names }) => (
                <ListItem key={completion.id} disablePadding sx={{ mb: 1 }}>
                  <ListItemButton
                    onClick={() => viewing.open(completion.id)}
                    component={Paper}
                    variant="outlined"
                    sx={{ p: 2 }}
                  >
                    <ListItemText
                      primary={chore.name}
                      secondary={names || undefined}
                    />
                  </ListItemButton>
                </ListItem>
              ))}
            </List>
          </Paper>
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
