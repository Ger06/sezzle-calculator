import { fireEvent, render, screen, within } from '@testing-library/react'
import { useReducer } from 'react'
import { describe, expect, it, vi } from 'vitest'
import { calculatorReducer, initialState } from '../state/calculatorReducer'
import Display from './Display'
import Keypad from './Keypad'

describe('Keypad', () => {
  it('dispatches DIGIT when a digit button is pressed (RF-13)', () => {
    const dispatch = vi.fn()
    render(<Keypad dispatch={dispatch} />)
    fireEvent.click(screen.getByRole('button', { name: '7' }))
    expect(dispatch).toHaveBeenCalledWith({ type: 'DIGIT', digit: '7' })
  })

  it('dispatches DIGIT for the 0 button', () => {
    const dispatch = vi.fn()
    render(<Keypad dispatch={dispatch} />)
    fireEvent.click(screen.getByRole('button', { name: '0' }))
    expect(dispatch).toHaveBeenCalledWith({ type: 'DIGIT', digit: '0' })
  })

  it('dispatches DECIMAL when the decimal-point button is pressed (RF-22)', () => {
    const dispatch = vi.fn()
    render(<Keypad dispatch={dispatch} />)
    fireEvent.click(screen.getByRole('button', { name: '.' }))
    expect(dispatch).toHaveBeenCalledWith({ type: 'DECIMAL' })
  })

  it('dispatches SIGN_TOGGLE when the sign-toggle (±) button is pressed (RF-23)', () => {
    const dispatch = vi.fn()
    render(<Keypad dispatch={dispatch} />)
    fireEvent.click(screen.getByRole('button', { name: 'Toggle sign' }))
    expect(dispatch).toHaveBeenCalledWith({ type: 'SIGN_TOGGLE' })
  })

  it('dispatches BACKSPACE when the delete-last-digit button is pressed (RF-20)', () => {
    const dispatch = vi.fn()
    render(<Keypad dispatch={dispatch} />)
    fireEvent.click(screen.getByRole('button', { name: 'Backspace' }))
    expect(dispatch).toHaveBeenCalledWith({ type: 'BACKSPACE' })
  })

  it('dispatches CLEAR when the AC button is pressed (RF-19)', () => {
    const dispatch = vi.fn()
    render(<Keypad dispatch={dispatch} />)
    fireEvent.click(screen.getByRole('button', { name: 'AC' }))
    expect(dispatch).toHaveBeenCalledWith({ type: 'CLEAR' })
  })
})

function Harness() {
  const [state, dispatch] = useReducer(calculatorReducer, initialState)
  return (
    <>
      <Display
        display={state.display}
        status={state.status}
        errorMessage={state.errorMessage}
      />
      <Keypad dispatch={dispatch} />
    </>
  )
}

describe('Keypad + Display wired through the real reducer', () => {
  it('a digit press updates the rendered display', () => {
    render(<Harness />)
    fireEvent.click(screen.getByRole('button', { name: '7' }))
    expect(within(screen.getByRole('status')).getByText('7')).toBeInTheDocument()
  })
})
