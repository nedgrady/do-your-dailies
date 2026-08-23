import { createFileRoute } from '@tanstack/react-router'
import ChoreListPage from '../../features/chore/manage/ChoreListPage'

export interface ChoresManageSearch {
  // The chore id being edited, or the literal 'new' while creating one.
  // Undefined means the edit/create dialog is closed. Kept in the URL so
  // the phone's back button closes the dialog (browser history) instead of
  // leaving the page.
  chore?: string
}

export const Route = createFileRoute('/chores/manage')({
  component: ChoreListPage,
  validateSearch: (search: Record<string, unknown>): ChoresManageSearch => ({
    chore: typeof search.chore === 'string' ? search.chore : undefined,
  }),
})
