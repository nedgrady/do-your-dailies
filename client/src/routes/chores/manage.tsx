import { createFileRoute } from '@tanstack/react-router'
import ChoreListPage from '../../features/chore/manage/ChoreListPage'

export const Route = createFileRoute('/chores/manage')({
  component: ChoreListPage,
})
