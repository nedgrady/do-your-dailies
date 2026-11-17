import {
  useNavigate,
  useRouter,
  useRouterState,
  useSearch,
} from '@tanstack/react-router'
import type { RootSearch } from '../routes/__root'

type SheetParam = 'completing' | 'completion'

// Opens and closes a sheet through one of the shared search params declared on
// the root route. Opening pushes a history entry so the browser/phone Back
// button closes the sheet; closing goes back one entry when we pushed it, so no
// duplicate entry is left behind. Only one sheet is open at a time, so setting
// a param replaces the whole search object.
function useSheetSearchParam(param: SheetParam) {
  const navigate = useNavigate()
  const router = useRouter()
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  })
  const value = useSearch({
    strict: false,
    select: (search) => search[param],
  })

  // The sheet can live on any page, so navigate to the current pathname. The
  // cast only tells the typed router that `search` is the root's shared search
  // schema, which is the only one these hooks touch.
  const navigateSearch = (id: number | undefined, replace: boolean) =>
    navigate({
      to: pathname as '/',
      search: (id === undefined ? {} : { [param]: id }) as RootSearch,
      replace,
    })

  const open = (id: number) => {
    void navigateSearch(id, false)
  }

  const close = () => {
    if (router.history.canGoBack()) {
      router.history.back()
      return
    }
    void navigateSearch(undefined, true)
  }

  return { value, open, close }
}

// The "Who did it?" sheet, keyed by the chore being completed.
export function useCompletingChore() {
  const { value, open, close } = useSheetSearchParam('completing')
  return { completingChoreId: value, open, close }
}

// The "Mark not done?" sheet, keyed by the completion being viewed.
export function useViewingCompletion() {
  const { value, open, close } = useSheetSearchParam('completion')
  return { viewingCompletionId: value, open, close }
}
