import {
  Box,
  Grid,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  Paper,
  Typography,
} from '@mui/material'
import { QueryErrorResetBoundary } from '@tanstack/react-query'
import { Suspense } from 'react'
import { ErrorBoundary } from 'react-error-boundary'
import { nullChore, useAllChoresQuery, type Chore } from '../../../domain/chore'
import {
  useChoreCompletionsBetweenQuery,
  useCreateChoreCompletionMutation,
} from '../../../domain/choreCompletion'
import {
  useCurrentDayChoreQueueQuery,
  type ChoreInQueue,
} from '../../../domain/choreQueue'

function ChoreQueueContent() {
  const { data: allChoresLookup = [] } = useAllChoresQuery()
  const { data: choresInQueue = [] } = useCurrentDayChoreQueueQuery()
  const { data: completions = [] } = useChoreCompletionsBetweenQuery(
    new Date(new Date().setHours(0, 0, 0, 0)),
    new Date(new Date().setHours(24, 0, 0, 0)),
  )

  const { mutate: createChoreCompletion } = useCreateChoreCompletionMutation()

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

  // TODO: wire this up to your actual "complete chore" mutation
  const handleChoreClick = (choreInQueue: ChoreInQueue) => {
    createChoreCompletion(choreInQueue.choreId)
  }

  return (
    <Box sx={{ mt: 2 }}>
      <Grid container spacing={2}>
        <Grid size={{ xs: 12, md: 6 }}>
          <Paper variant="outlined" sx={{ p: 2, height: '100%' }}>
            <Typography variant="h6" gutterBottom component="h2">
              To do
            </Typography>
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
              Completed
            </Typography>
            <List aria-label="Completed chores" disablePadding>
              {completedChores.map((chore) => (
                <ListItem key={chore.id} disablePadding sx={{ mb: 1 }}>
                  <Paper variant="outlined" sx={{ p: 2, width: '100%' }}>
                    <ListItemText primary={chore.name} />
                  </Paper>
                </ListItem>
              ))}
            </List>
          </Paper>
        </Grid>
      </Grid>
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
