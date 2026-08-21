// @vitest-environment jsdom

import { render, screen } from '#/test-utils'
import {
  createMemoryHistory,
  createRootRoute,
  createRouter,
  RouterProvider,
} from '@tanstack/react-router'
import { describe, expect, it } from 'vitest'
import Header from './Header'

function renderHeader() {
  const rootRoute = createRootRoute({ component: Header })
  const router = createRouter({
    routeTree: rootRoute,
    history: createMemoryHistory({ initialEntries: ['/'] }),
  })

  return render(<RouterProvider router={router} />)
}

describe('header', () => {
  it('renders a theme toggle control', async () => {
    renderHeader()

    expect(
      await screen.findByRole('button', { name: /toggle color mode/i }),
    ).toBeTruthy()
  })
})
