// @vitest-environment jsdom

import { apiBaseUrl } from '#/lib/axios'
import { server } from '#/mocks/server'
import { render, screen } from '#/test-utils'
import { http, HttpResponse } from 'msw'
import { describe, expect, it } from 'vitest'
import ChoreQueuePage from './ChoreQueuePage'

describe('chore queue page', () => {
  it('renders chores from the backend current-day queue endpoint', async () => {
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
    )

    render(<ChoreQueuePage />)

    expect(await screen.findByText('Dishes')).toBeTruthy()
  })
})
