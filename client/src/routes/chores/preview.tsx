import { createFileRoute } from '@tanstack/react-router'
import ChorePreviewPage from '../../features/chore/preview/ChorePreviewPage'

export const Route = createFileRoute('/chores/preview')({
  component: ChorePreviewPage,
})
