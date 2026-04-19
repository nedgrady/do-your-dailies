import { createFileRoute } from '@tanstack/react-router'
import ChoreQueuePage from '../features/chore-queue/ChoreQueuePage'

export const Route = createFileRoute('/chore-queue')({
  component: ChoreQueuePage,
})
