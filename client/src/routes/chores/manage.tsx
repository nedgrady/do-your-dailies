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
  // The chore id whose "Delete?" confirm sheet is open. A number, for the same
  // serialization reason as `chore`.
  deleting?: number
}

export const Route = createFileRoute('/chores/manage')({
  component: ChoreListPage,
  // Each param is validated on its own so one never drops the other.
  validateSearch: (search: Record<string, unknown>): ChoresManageSearch => {
    const result: ChoresManageSearch = {}
    if (search.chore === 'new') result.chore = 'new'
    else if (typeof search.chore === 'number') result.chore = search.chore
    if (typeof search.deleting === 'number') result.deleting = search.deleting
    return result
  },
})
