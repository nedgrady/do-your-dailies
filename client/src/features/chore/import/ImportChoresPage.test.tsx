// @vitest-environment jsdom

import { render, screen, waitFor } from '#/test-utils'
import type * as ReactRouter from '@tanstack/react-router'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import ImportChoresPage from './ImportChoresPage'

const { navigateMock } = vi.hoisted(() => ({ navigateMock: vi.fn() }))

vi.mock('@tanstack/react-router', async (importOriginal) => {
  const actual = await importOriginal<typeof ReactRouter>()
  return { ...actual, useNavigate: () => navigateMock }
})

// jsdom doesn't implement execCommand at all, so it's stubbed per-test below
// rather than faithfully re-simulated — what actually matters here is that
// our own code calls it correctly and doesn't move focus, not whether jsdom
// can round-trip a real browser's insertText behavior.
afterEach(() => {
  Reflect.deleteProperty(document, 'execCommand')
  navigateMock.mockClear()
})

async function pasteAndParse(text: string) {
  const user = userEvent.setup()
  const pasteBox = screen.getByRole('textbox', { name: 'Chores to import' })
  await user.click(pasteBox)
  await user.paste(text)
  await user.click(screen.getByRole('button', { name: 'Parse' }))
}

describe('ImportChoresPage', () => {
  it('navigates back to the chore list when Back to chores is clicked', async () => {
    render(<ImportChoresPage />)

    const user = userEvent.setup()
    await user.click(screen.getByRole('button', { name: 'Back to chores' }))

    expect(navigateMock).toHaveBeenCalledWith({ to: '/chores/manage' })
  })

  it('lists why each invalid row was rejected, and blocks submit', async () => {
    render(<ImportChoresPage />)

    await pasteAndParse('Dishes\t3d\nWeird\t3 x y')

    expect(await screen.findByText('Fix 1 row before importing')).toBeTruthy()
    expect(
      screen.getByText('Weird: Couldn\'t understand cadence "3 x y"'),
    ).toBeTruthy()
    expect(
      screen.getByRole('button', { name: /Add \d+ chores?/ }),
    ).toHaveProperty('disabled', true)

    const goodRow = screen.getByText('Dishes').closest('[role="row"]')
    const badRow = screen.getByText('Weird').closest('[role="row"]')
    expect(goodRow?.className).not.toContain('invalid-row')
    expect(badRow?.className).toContain('invalid-row')
  })

  it('shows no errors and enables submit once every row is valid', async () => {
    render(<ImportChoresPage />)

    await pasteAndParse('Dishes\t3d\nVacuum\t1w')

    expect(await screen.findByRole('grid')).toBeTruthy()
    expect(screen.queryByText(/before importing/)).toBeNull()
    expect(
      screen.getByRole('button', { name: /Add \d+ chores?/ }),
    ).toHaveProperty('disabled', false)
  })

  it('demotes Parse to a secondary button once rows exist, so only one primary button shows at a time', async () => {
    render(<ImportChoresPage />)

    const parseButton = screen.getByRole('button', { name: 'Parse' })
    expect(parseButton.className).toContain('MuiButton-contained')

    await pasteAndParse('Dishes\t3d')

    expect(parseButton.className).not.toContain('MuiButton-contained')
    expect(parseButton.className).toContain('MuiButton-outlined')
    expect(
      screen.getByRole('button', { name: /Add \d+ chores?/ }).className,
    ).toContain('MuiButton-contained')
  })

  it('reclaims primary for Parse once the box is edited again after parsing', async () => {
    render(<ImportChoresPage />)

    const parseButton = screen.getByRole('button', { name: 'Parse' })
    await pasteAndParse('Dishes\t3d')
    expect(parseButton.className).toContain('MuiButton-outlined')

    const user = userEvent.setup()
    await user.click(screen.getByRole('textbox', { name: 'Chores to import' }))
    await user.type(
      screen.getByRole('textbox', { name: 'Chores to import' }),
      'x',
    )

    expect(parseButton.className).toContain('MuiButton-contained')
    expect(parseButton.className).not.toContain('MuiButton-outlined')

    const addButton = screen.getByRole('button', { name: /Add \d+ chores?/ })
    expect(addButton.className).toContain('MuiButton-outlined')
    expect(addButton.className).not.toContain('MuiButton-contained')
  })

  it('explains a blank chore name separately from a bad cadence', async () => {
    render(<ImportChoresPage />)

    await pasteAndParse('\t3d\nWeird\t3 x y')

    expect(await screen.findByText('Fix 2 rows before importing')).toBeTruthy()
    expect(screen.getByText('Row 1: Chore name is empty')).toBeTruthy()
  })

  it('marks a row invalid once an edit commits, not while still typing', async () => {
    render(<ImportChoresPage />)

    await pasteAndParse('Dishes\t3d')

    const cell = screen
      .getByText('Dishes')
      .closest('[role="row"]')
      ?.querySelector('[role="gridcell"][data-field="cadenceAmount"]')
    expect(cell).toBeTruthy()

    const user = userEvent.setup()
    await user.click(cell as HTMLElement)
    const input = await screen.findByRole('spinbutton')
    await user.clear(input)
    await user.type(input, '0')

    // Still mid-edit — nothing has committed yet, so no error is shown.
    expect(
      screen.getByText('Dishes').closest('[role="row"]')?.className,
    ).not.toContain('invalid-row')

    await user.keyboard('{Enter}')

    await waitFor(() => {
      expect(
        screen.getByText('Dishes').closest('[role="row"]')?.className,
      ).toContain('invalid-row')
    })
    expect(await screen.findByText('Dishes: Cadence is missing')).toBeTruthy()
  })

  it('keeps the delimiter picker hidden until custom mode is selected', async () => {
    render(<ImportChoresPage />)

    expect(screen.queryByRole('radio')).toBeNull()

    const user = userEvent.setup()
    await user.click(screen.getByRole('button', { name: 'Custom settings' }))

    expect(screen.getByRole('radio', { name: 'Auto-detect' })).toBeTruthy()
    expect(screen.getByRole('radio', { name: 'New line' })).toBeTruthy()
  })

  it('parses comma-separated input once Comma is selected as the column delimiter', async () => {
    render(<ImportChoresPage />)

    const user = userEvent.setup()
    await user.click(screen.getByRole('button', { name: 'Custom settings' }))
    await user.click(screen.getByRole('radio', { name: 'Comma' }))
    await user.click(screen.getByRole('textbox', { name: 'Chores to import' }))
    await user.paste('Dishes,3d\nVacuum,1w')
    await user.click(screen.getByRole('button', { name: 'Parse' }))

    expect(await screen.findByText('Dishes')).toBeTruthy()
    expect(screen.getByText('Vacuum')).toBeTruthy()
    expect(screen.queryByText(/before importing/)).toBeNull()
  })

  it('parses a custom row delimiter once selected', async () => {
    render(<ImportChoresPage />)

    const user = userEvent.setup()
    await user.click(screen.getByRole('button', { name: 'Custom settings' }))
    await user.click(screen.getByRole('radio', { name: 'Semicolon' }))
    await user.click(screen.getByRole('textbox', { name: 'Chores to import' }))
    await user.paste('Dishes\t3d;Vacuum\t1w')
    await user.click(screen.getByRole('button', { name: 'Parse' }))

    expect(await screen.findByText('Dishes')).toBeTruthy()
    expect(screen.getByText('Vacuum')).toBeTruthy()
    expect(screen.queryByText(/before importing/)).toBeNull()
  })

  it('resets custom choices when switching back to Auto-detect', async () => {
    render(<ImportChoresPage />)

    const user = userEvent.setup()
    await user.click(screen.getByRole('button', { name: 'Custom settings' }))
    await user.click(screen.getByRole('radio', { name: 'Comma' }))

    await user.click(screen.getByRole('button', { name: 'Auto-detect' }))

    // The panel is hidden again, not just visually — the underlying choice
    // is genuinely back to auto-detect, so tab-delimited input still works.
    expect(screen.queryByRole('radio')).toBeNull()
    await user.click(screen.getByRole('textbox', { name: 'Chores to import' }))
    await user.paste('Dishes\t3d\nVacuum\t1w')
    await user.click(screen.getByRole('button', { name: 'Parse' }))

    expect(await screen.findByText('Dishes')).toBeTruthy()
    expect(screen.getByText('Vacuum')).toBeTruthy()
    expect(screen.queryByText(/before importing/)).toBeNull()

    // Reopening the panel shows the reset defaults, not the stale "Comma".
    await user.click(screen.getByRole('button', { name: 'Custom settings' }))
    expect(screen.getByRole('radio', { name: 'Auto-detect' })).toHaveProperty(
      'checked',
      true,
    )
  })

  it('inserts a literal tab instead of moving focus, and Esc+Tab moves focus normally', async () => {
    const execCommand = vi.fn(() => true)
    document.execCommand = execCommand
    render(<ImportChoresPage />)

    const pasteBox = screen.getByRole('textbox', {
      name: 'Chores to import',
    })
    const user = userEvent.setup()
    await user.click(pasteBox)
    await user.keyboard('{Tab}')

    // We don't re-simulate jsdom's missing native insertText behavior —
    // what matters is that our code invokes it correctly and doesn't let
    // focus move away (the keyboard-trap risk this is all about).
    expect(execCommand).toHaveBeenCalledWith('insertText', false, '\t')
    expect(document.activeElement).toBe(pasteBox)

    await user.keyboard('{Escape}{Tab}')
    expect(document.activeElement).not.toBe(pasteBox)
  })

  it('falls back to normal Tab behavior when execCommand is unavailable', async () => {
    render(<ImportChoresPage />)

    // eslint-disable-next-line @typescript-eslint/no-unnecessary-type-assertion -- tsc disagrees: without this cast, `.value` below fails to compile.
    const pasteBox = screen.getByRole('textbox', {
      name: 'Chores to import',
    }) as HTMLTextAreaElement
    const user = userEvent.setup()
    await user.click(pasteBox)
    await user.keyboard('a{Tab}b')

    expect(pasteBox.value).toBe('a')
    expect(document.activeElement).not.toBe(pasteBox)
  })
})
