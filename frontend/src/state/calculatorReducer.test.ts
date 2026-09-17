import { describe, expect, it } from 'vitest'
import type { State } from '../types/calculator'
import {
  calculatorReducer as reducer,
  initialState,
  NETWORK_ERROR_MESSAGE,
} from './calculatorReducer'

function run(state: State, ...actions: Parameters<typeof reducer>[1][]): State {
  return actions.reduce(reducer, state)
}


describe('DIGIT (RF-13)', () => {
  it('replaces the default 0 with the first digit typed', () => {
    const state = run(initialState, { type: 'DIGIT', digit: '5' })
    expect(state.display).toBe('5')
    expect(state.operandEntered).toBe(true)
  })

  it('appends subsequent digits to the operand being entered', () => {
    const state = run(
      initialState,
      { type: 'DIGIT', digit: '5' },
      { type: 'DIGIT', digit: '3' },
    )
    expect(state.display).toBe('53')
  })

  it('replaces a freshly displayed result rather than appending to it', () => {
    const afterResult = run(
      initialState,
      { type: 'DIGIT', digit: '9' },
      { type: 'OPERATOR', operator: '×' },
      { type: 'DIGIT', digit: '8' },
      { type: 'SUBMIT_START' },
      { type: 'SUBMIT_SUCCESS', result: 72 },
    )
    expect(afterResult.display).toBe('72')
    expect(afterResult.operandEntered).toBe(false)

    const state = reducer(afterResult, { type: 'DIGIT', digit: '5' })
    expect(state.display).toBe('5')
    expect(state.operandEntered).toBe(true)
  })
})

describe('DECIMAL (RF-22)', () => {
  it('appends a decimal point to the operand being entered', () => {
    const state = run(
      initialState,
      { type: 'DIGIT', digit: '3' },
      { type: 'DECIMAL' },
      { type: 'DIGIT', digit: '5' },
    )
    expect(state.display).toBe('3.5')
  })

  it('is a no-op if the operand already contains a decimal point', () => {
    const state = run(
      initialState,
      { type: 'DIGIT', digit: '3' },
      { type: 'DECIMAL' },
      { type: 'DECIMAL' },
    )
    expect(state.display).toBe('3.')
  })
})

describe('SIGN_TOGGLE (RF-23)', () => {
  it('multiplies the currently displayed value by -1', () => {
    const state = run(
      initialState,
      { type: 'DIGIT', digit: '5' },
      { type: 'SIGN_TOGGLE' },
    )
    expect(state.display).toBe('-5')
  })

  it('toggles back to positive on a second press', () => {
    const state = run(
      initialState,
      { type: 'DIGIT', digit: '5' },
      { type: 'SIGN_TOGGLE' },
      { type: 'SIGN_TOGGLE' },
    )
    expect(state.display).toBe('5')
  })
})

describe('BACKSPACE (RF-20)', () => {
  it('removes the last digit of the operand being entered', () => {
    const state = run(
      initialState,
      { type: 'DIGIT', digit: '5' },
      { type: 'DIGIT', digit: '3' },
      { type: 'BACKSPACE' },
    )
    expect(state.display).toBe('5')
  })

  it('resets to the empty/zero state when deleting the only remaining digit', () => {
    const state = run(
      initialState,
      { type: 'DIGIT', digit: '5' },
      { type: 'BACKSPACE' },
    )
    expect(state.display).toBe('0')
    expect(state.operandEntered).toBe(false)
  })

  it('is a no-op immediately after a freshly displayed result (frozen, matching DIGIT — bug fix)', () => {
    const afterResult = run(
      initialState,
      { type: 'DIGIT', digit: '9' },
      { type: 'OPERATOR', operator: '×' },
      { type: 'DIGIT', digit: '8' },
      { type: 'SUBMIT_START' },
      { type: 'SUBMIT_SUCCESS', result: 72 },
    )
    expect(afterResult.display).toBe('72')

    const state = reducer(afterResult, { type: 'BACKSPACE' })
    expect(state).toEqual(afterResult)
  })

  it('is also a no-op right after selecting an operator, before any second-operand digit is entered', () => {
    const afterOperator = run(
      initialState,
      { type: 'DIGIT', digit: '5' },
      { type: 'OPERATOR', operator: '+' },
    )
    expect(afterOperator.display).toBe('5')
    expect(afterOperator.operandEntered).toBe(false)

    const state = reducer(afterOperator, { type: 'BACKSPACE' })
    expect(state).toEqual(afterOperator)
  })
})

describe('CLEAR (RF-19)', () => {
  it('resets the calculator to its initial state from any prior state', () => {
    const state = run(
      initialState,
      { type: 'DIGIT', digit: '5' },
      { type: 'OPERATOR', operator: '+' },
      { type: 'DIGIT', digit: '3' },
      { type: 'CLEAR' },
    )
    expect(state).toEqual(initialState)
  })
})

describe('OPERATOR (RF-14)', () => {
  it('defaults to 0 as the first operand if none was entered', () => {
    const state = run(initialState, { type: 'OPERATOR', operator: '+' })
    expect(state.firstOperand).toBe(0)
    expect(state.operator).toBe('+')
  })

  it('captures the displayed value as the first operand', () => {
    const state = run(
      initialState,
      { type: 'DIGIT', digit: '5' },
      { type: 'OPERATOR', operator: '+' },
    )
    expect(state.firstOperand).toBe(5)
    expect(state.operator).toBe('+')
    expect(state.operandEntered).toBe(false)
  })

  it('replace branch: swaps the pending operator when no second-operand digit was entered', () => {
    const state = run(
      initialState,
      { type: 'DIGIT', digit: '5' },
      { type: 'OPERATOR', operator: '+' },
      { type: 'OPERATOR', operator: '×' },
    )
    expect(state.operator).toBe('×')
    expect(state.firstOperand).toBe(5)
    expect(state.display).toBe('5')
  })

  it('chain branch: after the caller resolves the pending computation (SUBMIT_SUCCESS), the new operator captures the result as the first operand', () => {
    const midChain = run(
      initialState,
      { type: 'DIGIT', digit: '5' },
      { type: 'OPERATOR', operator: '+' },
      { type: 'DIGIT', digit: '3' },
    )
    expect(midChain.operator).toBe('+')
    expect(midChain.operandEntered).toBe(true)

    const state = run(
      midChain,
      { type: 'SUBMIT_START' },
      { type: 'SUBMIT_SUCCESS', result: 8 },
      { type: 'OPERATOR', operator: '×' },
    )
    expect(state.firstOperand).toBe(8)
    expect(state.operator).toBe('×')
  })

  it('RF-27: selecting an operator after a displayed result uses that result as the first operand', () => {
    const afterResult = run(
      initialState,
      { type: 'DIGIT', digit: '2' },
      { type: 'OPERATOR', operator: '+' },
      { type: 'DIGIT', digit: '2' },
      { type: 'SUBMIT_START' },
      { type: 'SUBMIT_SUCCESS', result: 4 },
    )
    expect(afterResult.display).toBe('4')
    expect(afterResult.operator).toBeNull()

    const state = run(afterResult, { type: 'OPERATOR', operator: '+' })
    expect(state.firstOperand).toBe(4)
    expect(state.operator).toBe('+')
  })
})

describe('submission lifecycle (RF-16, RF-17, RF-25, RF-26)', () => {
  it('SUBMIT_START marks the calculator as loading and clears any prior error', () => {
    const state = run(
      { ...initialState, status: 'error', errorMessage: 'stale' },
      { type: 'SUBMIT_START' },
    )
    expect(state.status).toBe('loading')
    expect(state.errorMessage).toBeNull()
  })

  it('start → success lands in the correct terminal state', () => {
    const state = run(
      initialState,
      { type: 'DIGIT', digit: '2' },
      { type: 'OPERATOR', operator: '+' },
      { type: 'DIGIT', digit: '2' },
      { type: 'SUBMIT_START' },
      { type: 'SUBMIT_SUCCESS', result: 4 },
    )
    expect(state).toEqual({
      display: '4',
      operandEntered: false,
      firstOperand: null,
      operator: null,
      status: 'idle',
      errorMessage: null,
    })
  })

  it('start → error (domain/validation) lands in the correct terminal state', () => {
    const state = run(
      initialState,
      { type: 'DIGIT', digit: '1' },
      { type: 'OPERATOR', operator: '÷' },
      { type: 'DIGIT', digit: '0' },
      { type: 'SUBMIT_START' },
      { type: 'SUBMIT_ERROR', message: 'cannot divide by zero' },
    )
    expect(state.status).toBe('error')
    expect(state.errorMessage).toBe('cannot divide by zero')
  })

  it('start → network failure lands in a distinct terminal state with a connection-error message', () => {
    const state = run(
      initialState,
      { type: 'DIGIT', digit: '2' },
      { type: 'OPERATOR', operator: '+' },
      { type: 'DIGIT', digit: '2' },
      { type: 'SUBMIT_START' },
      { type: 'NETWORK_ERROR' },
    )
    expect(state.status).toBe('networkError')
    expect(state.errorMessage).toBe(NETWORK_ERROR_MESSAGE)
  })
})

describe('error recovery on new input (RF-24)', () => {
  it('a digit press clears the error and starts a fresh entry with that digit', () => {
    const errored = run(
      initialState,
      { type: 'SUBMIT_START' },
      { type: 'SUBMIT_ERROR', message: 'cannot divide by zero' },
    )
    const state = reducer(errored, { type: 'DIGIT', digit: '7' })
    expect(state.status).toBe('idle')
    expect(state.errorMessage).toBeNull()
    expect(state.display).toBe('7')
    expect(state.operandEntered).toBe(true)
    expect(state.operator).toBeNull()
  })

  it('an operator press clears the error and starts a fresh entry against a default first operand of 0', () => {
    const errored = run(
      initialState,
      { type: 'SUBMIT_START' },
      { type: 'NETWORK_ERROR' },
    )
    const state = reducer(errored, { type: 'OPERATOR', operator: '+' })
    expect(state.status).toBe('idle')
    expect(state.errorMessage).toBeNull()
    expect(state.firstOperand).toBe(0)
    expect(state.operator).toBe('+')
  })
})
