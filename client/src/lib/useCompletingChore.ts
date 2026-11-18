import {
  useNavigate,
  useRouter,
  useRouterState,
  useSearch,
} from '@tanstack/react-router'
import type { RootSearch } from '../routes/__root'

type SheetParam = keyof RootSearch

// Opens and closes a sheet through one of the shared search params declared on
// the root route. Opening pushes a history entry so the browser/phone Back
// button closes the sheet; closing goes back one entry when we pushed it, so no
// duplicate entry is left behind. Only one sheet is open at a time, so setting
// a param replaces the whole search object.
function useSheetSearchParam<TParam extends SheetParam>(param: TParam) {
  type Value = NonNullable<RootSearch[TParam]>

  const navigate = useNavigate()
  const router = useRouter()
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  })
  const value = useSearch({
    strict: false,
    select: (search) => search[param] as Value | undefined,
  })

  // The sheet can live on any page, so navigate to the current pathname. The
  // cast only tells the typed router that `search` is the root's shared search
  // schema, which is the only one these hooks touch.
  const navigateSearch = (next: Value | undefined, replace: boolean) =>
    navigate({
      to: pathname as '/',
      search: (next === undefined ? {} : { [param]: next }) as RootSearch,
      replace,
    })

  const open = (next: Value) => {
    void navigateSearch(next, false)
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

// The "Add completed chore" flow. Presence-only: the param is `true` when open.
export function usePickingChore() {
  const { value, open, close } = useSheetSearchParam('picking')
  return { isPicking: value === true, open: () => open(true), close }
}
