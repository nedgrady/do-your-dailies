import { createFileRoute } from '@tanstack/react-router'
import ChoreListPage from '../../features/chore/manage/ChoreListPage'

export interface ChoresManageSearch {
  // The chore id being edited, or the literal 'new' while creating one.
  // Undefined means the edit/create dialog is closed. Kept in the URL so
  // the phone's back button closes the dialog (browser history) instead of
  // leaving the page.
  //
  // A number, not a numeric string: the router's default search
  // serializer JSON-quotes any string that looks like a number (to
  // preserve its type on round-trip), which is what produced ugly
  // `?chore=%2284%22` URLs — a plain number serializes unquoted.
  chore?: number | 'new'
}

export const Route = createFileRoute('/chores/manage')({
  component: ChoreListPage,
  validateSearch: (search: Record<string, unknown>): ChoresManageSearch => {
    if (search.chore === 'new') return { chore: 'new' }
    if (typeof search.chore === 'number') return { chore: search.chore }
    return {}
  },
})
