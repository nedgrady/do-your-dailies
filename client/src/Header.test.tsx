// @vitest-environment jsdom

import { render, screen } from '#/test-utils'
import { describe, expect, it } from 'vitest'
import Header from './Header'

describe('header', () => {
  it('renders a theme toggle control', () => {
    render(<Header />)

    expect(
      screen.getByRole('button', { name: /toggle color mode/i }),
    ).toBeTruthy()
  })
})
