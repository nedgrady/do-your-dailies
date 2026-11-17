import {
  useNavigate,
  useRouter,
  useRouterState,
  useSearch,
} from '@tanstack/react-router'

// Opens and closes the "Who did it?" sheet through the shared `completing`
// search param (declared on the root route). Opening pushes a history entry so
// the browser/phone Back button closes the sheet; closing goes back one entry
// when we pushed it, so no duplicate entry is left behind.
export function useCompletingChore() {
  const navigate = useNavigate()
  const router = useRouter()
  const pathname = useRouterState({ select: (state) => state.location.pathname })
  const completingChoreId = useSearch({
    strict: false,
    select: (search) => search.completing,
  })

  // The sheet can live on any page, so navigate to the current pathname. The
  // `'/'` cast only tells the typed router that `search` is the root's shared
  // search schema, which is the only one this hook touches.
  const navigateSearch = (completing: number | undefined, replace: boolean) =>
    navigate({
      to: pathname as '/',
      search: completing === undefined ? {} : { completing },
      replace,
    })

  const open = (choreId: number) => {
    void navigateSearch(choreId, false)
  }

  const close = () => {
    if (router.history.canGoBack()) {
      router.history.back()
      return
    }
    void navigateSearch(undefined, true)
  }

  return { completingChoreId, open, close }
}
