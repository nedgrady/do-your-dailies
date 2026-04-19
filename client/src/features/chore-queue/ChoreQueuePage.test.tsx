// @vitest-environment jsdom

import { render, screen } from '#/test-utils'
import { describe, expect, it } from 'vitest'
import ChoreQueuePage from './ChoreQueuePage'

describe('chore queue page', () => {
  it('renders a title for required chores', () => {
    render(<ChoreQueuePage />)

    expect(screen.getByRole('heading', { level: 1 }).textContent).toBe(
      'Required Chores',
    )
  })
})
