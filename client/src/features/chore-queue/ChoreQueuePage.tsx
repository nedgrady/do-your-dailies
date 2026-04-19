import { useEffect, useState } from 'react'

type Chore = {
  id: number
  name: string
  cadenceInDays: number
  createdAt: string
  updatedAt: string
}

export default function ChoreQueuePage() {
  const [chores, setChores] = useState<Chore[]>([])
  const [hasError, setHasError] = useState(false)

  useEffect(() => {
    let isMounted = true

    async function loadQueue() {
      try {
        const response = await fetch('/api/chore-queue', {
          headers: {
            Accept: 'application/json',
          },
        })

        if (!response.ok) {
          throw new Error('queue request failed')
        }

        const queue = (await response.json()) as Chore[]
        if (isMounted) {
          setChores(queue)
        }
      } catch {
        if (isMounted) {
          setHasError(true)
        }
      }
    }

    void loadQueue()

    return () => {
      isMounted = false
    }
  }, [])

  return (
    <main>
      <h1>Required Chores</h1>
      {hasError ? <p>Unable to load chore queue.</p> : null}
      <ul aria-label="Current day chore queue">
        {chores.map((chore) => (
          <li key={chore.id}>{chore.name}</li>
        ))}
      </ul>
    </main>
  )
}
