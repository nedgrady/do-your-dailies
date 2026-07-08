// @vitest-environment jsdom

import { apiBaseUrl } from '#/lib/axios'
import { server } from '#/mocks/server'
import { render, screen } from '#/test-utils'
import userEvent from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { describe, expect, it } from 'vitest'
import ChoreQueuePage from './ChoreQueuePage'

describe('chore queue page', () => {
  it('shows chores completed today even when not in the current-day queue', async () => {
    server.use(
      http.get(`${apiBaseUrl}/api/chore-queue`, () => {
        return HttpResponse.json([])
      }),
      http.get(`${apiBaseUrl}/api/chore-completions`, () => {
        return HttpResponse.json([
          {
            id: 100,
            choreId: 1,
            createdAt: '2026-04-21T00:00:00Z',
            updatedAt: '2026-04-21T00:00:00Z',
          },
        ])
      }),
      http.get(`${apiBaseUrl}/api/chores`, () => {
        return HttpResponse.json([
          {
            id: 1,
            name: 'Dishes',
            cadenceInDays: 1,
            createdAt: '2026-04-01T00:00:00Z',
            updatedAt: '2026-04-01T00:00:00Z',
          },
        ])
      }),
    )

    render(<ChoreQueuePage />)

    expect(
      await screen.findByRole('button', { name: 'Dishes completed' }),
    ).toBeTruthy()
  })

  it('shows day-off success state when no chores are due today', async () => {
    server.use(
      http.get(`${apiBaseUrl}/api/chore-queue`, () => {
        return HttpResponse.json([])
      }),
      http.get(`${apiBaseUrl}/api/chore-completions`, () => {
        return HttpResponse.json([])
      }),
      http.get(`${apiBaseUrl}/api/chores`, () => {
        return HttpResponse.json([])
      }),
    )

    render(<ChoreQueuePage />)

    expect(
      await screen.findByRole('status', { name: 'Day off today' }),
    ).toBeTruthy()
  })

  it('shows all-done success state when all queued chores are completed today', async () => {
    server.use(
      http.get(`${apiBaseUrl}/api/chore-queue`, () => {
        return HttpResponse.json([
          {
            id: 1,
            name: 'Dishes',
            cadenceInDays: 1,
            createdAt: '2026-04-01T00:00:00Z',
            updatedAt: '2026-04-01T00:00:00Z',
          },
        ])
      }),
      http.get(`${apiBaseUrl}/api/chore-completions`, () => {
        return HttpResponse.json([
          {
            id: 100,
            choreId: 1,
            createdAt: '2026-04-21T00:00:00Z',
            updatedAt: '2026-04-21T00:00:00Z',
          },
        ])
      }),
      http.get(`${apiBaseUrl}/api/chores`, () => {
        return HttpResponse.json([
          {
            id: 1,
            name: 'Dishes',
            cadenceInDays: 1,
            createdAt: '2026-04-01T00:00:00Z',
            updatedAt: '2026-04-01T00:00:00Z',
          },
        ])
      }),
    )

    render(<ChoreQueuePage />)

    expect(
      await screen.findByRole('status', { name: 'All chores completed today' }),
    ).toBeTruthy()
  })

  it('posts completion and keeps chore visible as completed when tapped', async () => {
    let postedChoreID = 0
    let completionChoreIDs: number[] = []
    server.use(
      http.get(`${apiBaseUrl}/api/chore-queue`, () => {
        return HttpResponse.json([
          {
            id: 1,
            name: 'Dishes',
            cadenceInDays: 1,
            createdAt: '2026-04-01T00:00:00Z',
            updatedAt: '2026-04-01T00:00:00Z',
          },
        ])
      }),
      http.get(`${apiBaseUrl}/api/chore-completions`, ({ request }) => {
        const date = new URL(request.url).searchParams.get('date')
        if (!date || !/^\d{4}-\d{2}-\d{2}$/.test(date)) {
          return HttpResponse.json(
            { error: 'invalid date query' },
            { status: 400 },
          )
        }

        return HttpResponse.json(
          completionChoreIDs.map((choreId, index) => ({
            id: index + 1,
            choreId,
            createdAt: '2026-04-21T00:00:00Z',
            updatedAt: '2026-04-21T00:00:00Z',
          })),
        )
      }),
      http.get(`${apiBaseUrl}/api/chores`, () => {
        return HttpResponse.json([
          {
            id: 1,
            name: 'Dishes',
            cadenceInDays: 1,
            createdAt: '2026-04-01T00:00:00Z',
            updatedAt: '2026-04-01T00:00:00Z',
          },
        ])
      }),
      http.post(`${apiBaseUrl}/api/chore-completions`, async ({ request }) => {
        const body = (await request.json()) as { choreId: number }
        postedChoreID = body.choreId
        completionChoreIDs = [body.choreId]

        return HttpResponse.json(
          {
            id: 999,
            choreId: body.choreId,
            createdAt: '2026-04-21T00:00:00Z',
            updatedAt: '2026-04-21T00:00:00Z',
          },
          { status: 201 },
        )
      }),
    )

    render(<ChoreQueuePage />)
    const user = userEvent.setup()

    await user.click(
      await screen.findByRole('button', {
        name: 'Mark Dishes as done',
      }),
    )

    expect(
      await screen.findByRole('button', { name: 'Dishes completed' }),
    ).toBeTruthy()
    expect(postedChoreID).toEqual(1)
  })

  it('shows inline item error when marking a chore as done fails', async () => {
    server.use(
      http.get(`${apiBaseUrl}/api/chore-queue`, () => {
        return HttpResponse.json([
          {
            id: 1,
            name: 'Dishes',
            cadenceInDays: 1,
            createdAt: '2026-04-01T00:00:00Z',
            updatedAt: '2026-04-01T00:00:00Z',
          },
        ])
      }),
      http.get(`${apiBaseUrl}/api/chore-completions`, () => {
        return HttpResponse.json([])
      }),
      http.get(`${apiBaseUrl}/api/chores`, () => {
        return HttpResponse.json([
          {
            id: 1,
            name: 'Dishes',
            cadenceInDays: 1,
            createdAt: '2026-04-01T00:00:00Z',
            updatedAt: '2026-04-01T00:00:00Z',
          },
        ])
      }),
      http.post(`${apiBaseUrl}/api/chore-completions`, () => {
        return HttpResponse.error()
      }),
    )

    render(<ChoreQueuePage />)
    const user = userEvent.setup()

    await user.click(
      await screen.findByRole('button', { name: 'Mark Dishes as done' }),
    )

    expect((await screen.findByRole('alert')).textContent).toEqual(
      'Could not mark this chore as done.',
    )
  })
})
