import { setupWorker } from 'msw/browser'
import { handlers } from './handlers'

export const worker = setupWorker(...handlers)

let workerStartPromise: Promise<void> | undefined

export function enableMocking() {
  if (!import.meta.env.DEV || import.meta.env.VITE_ENABLE_MSW !== 'true') {
    return Promise.resolve()
  }

  if (!workerStartPromise) {
    workerStartPromise = worker
      .start({ onUnhandledRequest: 'bypass' })
      .then(() => undefined)
  }

  return workerStartPromise
}
