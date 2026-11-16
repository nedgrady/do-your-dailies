// @vitest-environment jsdom

import { apiBaseUrl } from '#/lib/axios'
import { server } from '#/mocks/server'
import { render, screen, within } from '#/test-utils'
import { http, HttpResponse } from 'msw'
import { describe, expect, it } from 'vitest'
import ChorePreviewPage from './ChorePreviewPage'

describe('ChorePreviewPage', () => {
  it("renders each day heading with that day's chores, including a daily chore repeated across every day", async () => {
    server.use(
      http.get(`${apiBaseUrl}/api/chore-preview`, () => {
        return HttpResponse.json({
          days: [
            {
              dayOffset: 1,
              chores: [
                { choreId: 1, choreName: 'Sweep floor' },
                { choreId: 2, choreName: 'Water plants' },
              ],
            },
            {
              dayOffset: 2,
              chores: [{ choreId: 1, choreName: 'Sweep floor' }],
            },
          ],
        })
      }),
    )

    render(<ChorePreviewPage />)

    const day1 = await screen.findByRole('heading', { name: 'Day 1' })
    const day2 = await screen.findByRole('heading', { name: 'Day 2' })

    expect(
      within(day1.closest('.MuiPaper-root') as HTMLElement).getByText(
        'Sweep floor',
      ),
    ).toBeTruthy()
    expect(
      within(day1.closest('.MuiPaper-root') as HTMLElement).getByText(
        'Water plants',
      ),
    ).toBeTruthy()
    expect(
      within(day2.closest('.MuiPaper-root') as HTMLElement).getByText(
        'Sweep floor',
      ),
    ).toBeTruthy()
  })

  it('shows an error state when the preview fails to load', async () => {
    server.use(
      http.get(`${apiBaseUrl}/api/chore-preview`, () => {
        return HttpResponse.json({ error: 'boom' }, { status: 500 })
      }),
    )

    render(<ChorePreviewPage />)

    expect(
      await screen.findByText('Could not load chore preview.'),
    ).toBeTruthy()
  })
})
