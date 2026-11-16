import {
  Box,
  List,
  ListItem,
  ListItemText,
  Paper,
  Stack,
  Typography,
} from '@mui/material'
import { QueryErrorResetBoundary } from '@tanstack/react-query'
import { Suspense } from 'react'
import { ErrorBoundary } from 'react-error-boundary'
import { useChorePreviewQuery } from '../../../domain/chorePreview'

const PREVIEW_DAYS = 7

function ChorePreviewContent() {
  const { data } = useChorePreviewQuery(PREVIEW_DAYS)

  return (
    <Stack spacing={2}>
      {data.days.map((day) => (
        <Paper key={day.dayOffset} variant="outlined" sx={{ p: 2 }}>
          <Typography variant="h6" gutterBottom component="h2">
            Day {day.dayOffset}
          </Typography>
          <List aria-label={`Day ${day.dayOffset} chores`} disablePadding>
            {day.chores.map((chore, index) => (
              <ListItem key={index} disablePadding sx={{ mb: 1 }}>
                <ListItemText primary={chore.choreName} />
              </ListItem>
            ))}
          </List>
        </Paper>
      ))}
    </Stack>
  )
}

export default function ChorePreviewPage() {
  return (
    <Box>
      <Typography variant="h5" sx={{ mb: 1 }}>
        Preview upcoming days
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
        Assuming you complete every recommended chore each day, here's what the
        next {PREVIEW_DAYS} days would look like.
      </Typography>
      <QueryErrorResetBoundary>
        {({ reset }) => (
          <ErrorBoundary
            onReset={reset}
            fallbackRender={({ resetErrorBoundary }) => (
              <>
                <p>Could not load chore preview.</p>
                <button type="button" onClick={resetErrorBoundary}>
                  Retry
                </button>
              </>
            )}
          >
            <Suspense fallback={<p>Loading chore preview...</p>}>
              <ChorePreviewContent />
            </Suspense>
          </ErrorBoundary>
        )}
      </QueryErrorResetBoundary>
    </Box>
  )
}
