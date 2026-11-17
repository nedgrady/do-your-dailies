import type { AxiosRequestConfig } from 'axios'
import axios from './axios'
import { createRecentKeyStore } from './idempotencyKey'

const keys = createRecentKeyStore()

// Identical POSTs that are already in flight share one promise, so a double
// click or a repeated commit never reaches the server twice.
const inFlight = new Map<string, Promise<unknown>>()

function postFingerprint(config: AxiosRequestConfig): string {
  const body = typeof config.data === 'string' ? config.data : JSON.stringify(config.data ?? null)
  return `POST ${config.url ?? ''} ${body}`
}

export const customInstance = async <T>(config: AxiosRequestConfig): Promise<T> => {
  if (config.method?.toUpperCase() !== 'POST') {
    const response = await axios(config)
    return response.data
  }

  const fingerprint = postFingerprint(config)
  const pending = inFlight.get(fingerprint)
  if (pending) return pending as Promise<T>

  const request = axios({
    ...config,
    headers: { ...config.headers, 'Idempotency-Key': keys.keyFor(fingerprint) },
  })
    .then((response) => {
      keys.markSucceeded(fingerprint)
      return response.data as T
    })
    .catch((error: unknown) => {
      keys.markFailed(fingerprint)
      throw error
    })
    .finally(() => {
      inFlight.delete(fingerprint)
    })

  inFlight.set(fingerprint, request)
  return request
}
