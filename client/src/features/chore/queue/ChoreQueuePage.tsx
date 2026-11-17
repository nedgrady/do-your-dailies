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
import { nullChore, useAllChoresQuery } from '../../../domain/chore'
import {
  useChoreCompletionsBetweenQuery,
  useCreateChoreCompletionMutation,
} from '../../../domain/choreCompletion'
import {
  useChoreCompletorsQuery,
  useCreateChoreCompletorMutation,
} from '../../../domain/choreCompletor'
import { useCurrentDayChoreQueueQuery } from '../../../domain/choreQueue'
import type { ChoreInQueue } from '../../../domain/choreQueue'
import { useCompletingChore } from '../../../lib/useCompletingChore'
import WhoDidItDialog from './WhoDidItDialog'

function ChoreQueueContent() {
  const { data: allChoresLookup = [] } = useAllChoresQuery()
  const { data: choresInQueue = [] } = useCurrentDayChoreQueueQuery()
  const { data: completions = [] } = useChoreCompletionsBetweenQuery(
    new Date(new Date().setHours(0, 0, 0, 0)),
    new Date(new Date().setHours(24, 0, 0, 0)),
  )

  const { data: completors = [] } = useChoreCompletorsQuery()
  const { mutate: createChoreCompletion } = useCreateChoreCompletionMutation()
  const { mutateAsync: createChoreCompletor } =
    useCreateChoreCompletorMutation()
  const { completingChoreId, open, close } = useCompletingChore()

  const toDoChores = choresInQueue.filter(
    (choreInQueue) =>
      !completions.some(
        (completion) => completion.choreId === choreInQueue.choreId,
      ),
  )

  // A stale, invalid or already-completed id in the URL just shows no sheet.
  const choreBeingCompleted =
    toDoChores.find((chore) => chore.choreId === completingChoreId) ?? null

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

  const handleChoreClick = (choreInQueue: ChoreInQueue) => {
    open(choreInQueue.choreId)
  }

  const handleConfirmCompletion = (choreCompletorIds: number[]) => {
    if (!choreBeingCompleted) return
    createChoreCompletion({
      data: {
        choreId: choreBeingCompleted.choreId,
        choreCompletorIds,
      },
    })
    close()
  }

  const handleAddPerson = async (name: string) => {
    const created = await createChoreCompletor({ data: { name } })
    return created.id
  }

  return (
    <Box>
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
              {completedEntries.map(({ completion, chore, names }) => (
                <ListItem key={completion.id} disablePadding sx={{ mb: 1 }}>
                  <Paper variant="outlined" sx={{ p: 2, width: '100%' }}>
                    <ListItemText
                      primary={chore.name}
                      secondary={names || undefined}
                    />
                  </Paper>
                </ListItem>
              ))}
            </List>
          </Paper>
        </Grid>
      </Grid>

      {choreBeingCompleted && (
        <WhoDidItDialog
          key={choreBeingCompleted.choreId}
          choreName={choreBeingCompleted.choreName}
          completors={completors}
          onClose={close}
          onConfirm={handleConfirmCompletion}
          onAddPerson={handleAddPerson}
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
