import {
  Box,
  Button,
  Divider,
  Stack,
  TextField,
  Typography,
} from '@mui/material'
import { useState } from 'react'
import ImportChoresPage from '../../chore/import/ImportChoresPage'
import { ChoreInsights } from '../../chore/manage/ChoreInsights'
import ChoreQueuePage from '../../chore/queue/ChoreQueuePage'
import {
  SimulatedUserProvider,
  useSimulatedUser,
} from '../../../lib/simulatedUser'

const RECENT_IDS_STORAGE_KEY = 'simulator.recentUserIds'

function loadRecentIds(): number[] {
  try {
    const raw = localStorage.getItem(RECENT_IDS_STORAGE_KEY)
    const parsed: unknown = raw ? JSON.parse(raw) : []
    return Array.isArray(parsed)
      ? parsed.filter((n): n is number => typeof n === 'number')
      : []
  } catch {
    return []
  }
}

function rememberRecentId(id: number): number[] {
  const updated = [id, ...loadRecentIds().filter((n) => n !== id)].slice(0, 8)
  try {
    localStorage.setItem(RECENT_IDS_STORAGE_KEY, JSON.stringify(updated))
  } catch {
    // Convenience only - a full/blocked localStorage just means no history.
  }
  return updated
}

// Reuses the app's own QueryClient and the real bulk-import/queue/insights
// pages unmodified - isolation from real data (and between different
// synthetic ids) comes entirely from every query/mutation key being scoped
// by the acting user id (see domain/queryKeys.ts + lib/simulatedUser), not
// from a separate cache instance. That also means switching back to a
// previously-viewed synthetic id is an instant cache hit, not a refetch.
function SimulatorPageContent() {
  const { userId: activeUserId, setUserId } = useSimulatedUser()
  const [recentIds, setRecentIds] = useState<number[]>(() => loadRecentIds())
  const [userIdInput, setUserIdInput] = useState('')

  const activate = (id: number) => {
    setRecentIds(rememberRecentId(id))
    setUserId(id)
  }

  return (
    <Box>
      <Typography variant="h5" sx={{ mb: 1 }}>
        Cold-start simulator
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
        Acts as an arbitrary synthetic user id so you can bulk-import backdated
        chores and see the resulting queue and insights, without touching real
        data. Only does anything if the server also has SIMULATION_MODE_ENABLED
        set.
      </Typography>
      <Stack
        direction="row"
        spacing={1}
        sx={{ mb: 3, alignItems: 'center', flexWrap: 'wrap' }}
      >
        <TextField
          label="Synthetic user id"
          size="small"
          value={userIdInput}
          onChange={(e) => setUserIdInput(e.target.value)}
        />
        <Button
          variant="contained"
          disabled={!/^\d+$/.test(userIdInput)}
          onClick={() => activate(Number(userIdInput))}
        >
          Act as this user
        </Button>
        {recentIds.map((id) => (
          <Button
            key={id}
            size="small"
            variant={id === activeUserId ? 'contained' : 'outlined'}
            onClick={() => activate(id)}
          >
            {id}
          </Button>
        ))}
      </Stack>

      {activeUserId !== null && (
        <Stack spacing={3}>
          <ImportChoresPage />
          <Divider />
          <Typography variant="h6">Resulting queue</Typography>
          <ChoreQueuePage />
          <Divider />
          <Typography variant="h6">Resulting insights</Typography>
          <ChoreInsights />
        </Stack>
      )}
    </Box>
  )
}

export default function SimulatorPage() {
  const enabled = import.meta.env.VITE_ENABLE_SIMULATOR === 'true'

  if (!enabled) {
    return (
      <Typography color="text.secondary">
        The cold-start simulator isn't enabled in this build.
      </Typography>
    )
  }

  return (
    <SimulatedUserProvider>
      <SimulatorPageContent />
    </SimulatedUserProvider>
  )
}
