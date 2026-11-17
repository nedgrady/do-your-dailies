// Idempotency-Key handling for POST requests. The generated client sends every
// request through customInstance (orval-axios-mutator.ts), which uses this module
// to pick the key, so components never handle keys themselves.

// Keys are GUIDs. randomUUID exists only in secure contexts; getRandomValues works
// everywhere, including plain http:// on a LAN address used for phone testing.
export function newIdempotencyKey(): string {
  if (typeof crypto.randomUUID === 'function') {
    return crypto.randomUUID()
  }
  const bytes = crypto.getRandomValues(new Uint8Array(16))
  bytes[6] = (bytes[6] & 0x0f) | 0x40
  bytes[8] = (bytes[8] & 0x3f) | 0x80
  const hex = Array.from(bytes, (byte) => byte.toString(16).padStart(2, '0')).join('')
  return [hex.slice(0, 8), hex.slice(8, 12), hex.slice(12, 16), hex.slice(16, 20), hex.slice(20)].join('-')
}

export interface RecentKeyStore {
  // The key to send for this request. The same key is returned for the same
  // fingerprint while it is in flight, and for windowMs after it succeeded.
  keyFor: (fingerprint: string) => string
  markSucceeded: (fingerprint: string) => void
  // A failed attempt forgets its key, so the next attempt gets a new one.
  markFailed: (fingerprint: string) => void
}

interface Entry {
  key: string
  // null while in flight; otherwise the time after which the entry is forgotten.
  expiresAt: number | null
}

// Pure apart from the clock and key generator, so it can be tested without HTTP.
export function createRecentKeyStore({
  now = () => Date.now(),
  windowMs = 5000,
  newKey = newIdempotencyKey,
}: {
  now?: () => number
  windowMs?: number
  newKey?: () => string
} = {}): RecentKeyStore {
  const entries = new Map<string, Entry>()

  const sweep = () => {
    const current = now()
    for (const [fingerprint, entry] of entries) {
      if (entry.expiresAt !== null && entry.expiresAt <= current) {
        entries.delete(fingerprint)
      }
    }
  }

  return {
    keyFor(fingerprint) {
      sweep()
      const existing = entries.get(fingerprint)
      if (existing) return existing.key
      const key = newKey()
      entries.set(fingerprint, { key, expiresAt: null })
      return key
    },
    markSucceeded(fingerprint) {
      const entry = entries.get(fingerprint)
      if (entry) {
        entries.set(fingerprint, { key: entry.key, expiresAt: now() + windowMs })
      }
    },
    markFailed(fingerprint) {
      entries.delete(fingerprint)
    },
  }
}
