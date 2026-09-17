import { fireEvent, render, screen, within } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import * as calculatorApi from '../api/calculatorApi'
import Calculator from './Calculator'

vi.mock('../api/calculatorApi')

beforeEach(() => {
  vi.resetAllMocks()
})

function display() {
  return screen.getByRole('status')
}

describe('Calculator — full binary-operation flow', () => {
  it('computes base^exponent end-to-end and displays the result (^ reaching general power, not just x²)', async () => {
    vi.mocked(calculatorApi.power).mockResolvedValueOnce({ ok: true, result: 1024 })

    render(<Calculator />)

    fireEvent.click(screen.getByRole('button', { name: '2' }))
    fireEvent.click(screen.getByRole('button', { name: '^' }))
    fireEvent.click(screen.getByRole('button', { name: '1' }))
    fireEvent.click(screen.getByRole('button', { name: '0' }))
    fireEvent.click(screen.getByRole('button', { name: '=' }))

    expect(await within(display()).findByText('1024')).toBeInTheDocument()
    expect(calculatorApi.power).toHaveBeenCalledWith(2, 10)
  })
})

describe('Calculator — unary shortcut', () => {
  it('√x submits the currently displayed value and shows the result', async () => {
    vi.mocked(calculatorApi.sqrt).mockResolvedValueOnce({ ok: true, result: 4 })

    render(<Calculator />)

    fireEvent.click(screen.getByRole('button', { name: '1' }))
    fireEvent.click(screen.getByRole('button', { name: '6' }))
    fireEvent.click(screen.getByRole('button', { name: '√x' }))

    expect(await within(display()).findByText('4')).toBeInTheDocument()
    expect(calculatorApi.sqrt).toHaveBeenCalledWith(16)
  })
})

describe('Calculator — error path', () => {
  it('displays the domain error message instead of a result', async () => {
    vi.mocked(calculatorApi.divide).mockResolvedValueOnce({
      ok: false,
      kind: 'domain',
      code: 'DIVISION_BY_ZERO',
      message: 'cannot divide by zero',
    })

    render(<Calculator />)

    fireEvent.click(screen.getByRole('button', { name: '1' }))
    fireEvent.click(screen.getByRole('button', { name: '÷' }))
    fireEvent.click(screen.getByRole('button', { name: '0' }))
    fireEvent.click(screen.getByRole('button', { name: '=' }))

    expect(
      await within(display()).findByText('cannot divide by zero'),
    ).toBeInTheDocument()
  })
})

describe('Calculator — keyboard equivalence (RF-21)', () => {
  it('typing on the keyboard produces the same result as clicking the equivalent buttons', async () => {
    vi.mocked(calculatorApi.add).mockResolvedValueOnce({ ok: true, result: 8 })

    render(<Calculator />)

    fireEvent.keyDown(window, { key: '5' })
    fireEvent.keyDown(window, { key: '+' })
    fireEvent.keyDown(window, { key: '3' })
    fireEvent.keyDown(window, { key: 'Enter' })

    expect(await within(display()).findByText('8')).toBeInTheDocument()
    expect(calculatorApi.add).toHaveBeenCalledWith(5, 3)
  })
})

describe('Calculator — keyboard focus vs. global shortcuts (accessibility bug fix)', () => {
  it('does not fire the global equals shortcut when Enter is pressed while focus is on a button', () => {
    render(<Calculator />)

    fireEvent.click(screen.getByRole('button', { name: '5' }))
    fireEvent.click(screen.getByRole('button', { name: '+' }))
    fireEvent.click(screen.getByRole('button', { name: '3' }))

    const sevenButton = screen.getByRole('button', { name: '7' })
    sevenButton.focus()
    fireEvent.keyDown(sevenButton, { key: 'Enter' })

    expect(calculatorApi.add).not.toHaveBeenCalled()

  })

  it('digit keys still work via the global shortcut when the keydown does not originate from a button', () => {
    render(<Calculator />)

    fireEvent.keyDown(window, { key: '5' })

    expect(within(display()).getByText('5')).toBeInTheDocument()
  })
})
