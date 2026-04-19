import { QueryClient } from '@tanstack/react-query'
import { isAxiosError } from 'axios'

function shouldRetryQuery(failureCount: number, error: unknown) {
  if (failureCount >= 3) {
    return false
  }

  if (error instanceof Error && error.name === 'ZodError') {
    return false
  }

  if (isAxiosError(error)) {
    const status = error.response?.status

    if (typeof status === 'number' && status >= 400 && status < 500) {
      return status === 408 || status === 429
    }
  }

  return true
}

export function getContext() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: shouldRetryQuery,
      },
    },
  })

  return {
    queryClient,
  }
}
export default function TanstackQueryProvider() {}
