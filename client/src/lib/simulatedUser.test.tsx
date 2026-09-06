// @vitest-environment jsdom

import { render, screen } from '#/test-utils'
import { renderHook } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { afterEach, describe, expect, it } from 'vitest'
import { server } from '#/mocks/server'
import axiosInstance, { apiBaseUrl } from './axios'
import { SimulatedUserProvider, useSimulatedUser } from './simulatedUser'

function Consumer() {
  const { userId, setUserId } = useSimulatedUser()
  return (
    <div>
      <span>current: {userId ?? 'none'}</span>
      <button onClick={() => setUserId(42)}>set to 42</button>
      <button onClick={() => setUserId(null)}>clear</button>
    </div>
  )
}

function mockProbe(): { capturedHeader: () => string | null } {
  let captured: string | null = null
  server.use(
    http.get(`${apiBaseUrl}/api/probe`, ({ request }) => {
      captured = request.headers.get('X-Simulated-User-Id')
      return HttpResponse.json({})
    }),
  )
  return { capturedHeader: () => captured }
}

describe('SimulatedUserProvider / useSimulatedUser', () => {
  afterEach(() => {
    delete axiosInstance.defaults.headers.common['X-Simulated-User-Id']
  })

  it('does not send the header by default', async () => {
    const probe = mockProbe()
    render(
      <SimulatedUserProvider>
        <Consumer />
      </SimulatedUserProvider>,
    )

    await axiosInstance.get('/api/probe')
    expect(probe.capturedHeader()).toBeNull()
  })

  it('sends the header once setUserId is called, and updates context state', async () => {
    const probe = mockProbe()
    render(
      <SimulatedUserProvider>
        <Consumer />
      </SimulatedUserProvider>,
    )

    const user = userEvent.setup()
    await user.click(screen.getByRole('button', { name: 'set to 42' }))

    expect(screen.getByText('current: 42')).toBeTruthy()
    await axiosInstance.get('/api/probe')
    expect(probe.capturedHeader()).toBe('42')
  })

  it('stops sending the header once cleared - the "must never leak" case', async () => {
    const probe = mockProbe()
    render(
      <SimulatedUserProvider>
        <Consumer />
      </SimulatedUserProvider>,
    )

    const user = userEvent.setup()
    await user.click(screen.getByRole('button', { name: 'set to 42' }))
    await user.click(screen.getByRole('button', { name: 'clear' }))

    expect(screen.getByText('current: none')).toBeTruthy()
    await axiosInstance.get('/api/probe')
    expect(probe.capturedHeader()).toBeNull()
  })

  it('clears the header when the provider itself unmounts', async () => {
    const probe = mockProbe()
    const { unmount } = render(
      <SimulatedUserProvider>
        <Consumer />
      </SimulatedUserProvider>,
    )

    const user = userEvent.setup()
    await user.click(screen.getByRole('button', { name: 'set to 42' }))

    unmount()

    await axiosInstance.get('/api/probe')
    expect(probe.capturedHeader()).toBeNull()
  })

  it('resolves to userId: null when rendered outside any provider - real pages act as self', () => {
    render(<Consumer />)
    expect(screen.getByText('current: none')).toBeTruthy()
  })

  it('throws if setUserId is ever called outside a provider', () => {
    const { result } = renderHook(() => useSimulatedUser())
    expect(() => result.current.setUserId(42)).toThrow(/SimulatedUserProvider/)
  })
})
