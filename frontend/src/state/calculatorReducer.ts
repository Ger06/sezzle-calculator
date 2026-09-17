import type { Action, State } from '../types/calculator'

export const NETWORK_ERROR_MESSAGE =
  'Unable to reach the server. Please check your connection and try again.'

export const initialState: State = {
  display: '0',
  operandEntered: false,
  firstOperand: null,
  operator: null,
  status: 'idle',
  errorMessage: null,
}

function isErrorStatus(status: State['status']): boolean {
  return status === 'error' || status === 'networkError'
}

export function calculatorReducer(state: State, action: Action): State {
  switch (action.type) {
    case 'DIGIT': {
      if (isErrorStatus(state.status)) {
        return { ...initialState, display: action.digit, operandEntered: true }
      }
      const display = state.operandEntered
        ? state.display + action.digit
        : action.digit
      return { ...state, display, operandEntered: true }
    }

    case 'DECIMAL': {
      if (isErrorStatus(state.status)) return state
      if (state.display.includes('.')) return state
      const display = state.operandEntered ? state.display + '.' : '0.'
      return { ...state, display, operandEntered: true }
    }

    case 'SIGN_TOGGLE': {
      if (isErrorStatus(state.status)) return state
      return { ...state, display: String(Number(state.display) * -1) }
    }

    case 'BACKSPACE': {
      if (isErrorStatus(state.status)) return state
      if (!state.operandEntered) return state
      const next = state.display.slice(0, -1)
      if (next === '' || next === '-') {
        return { ...state, display: '0', operandEntered: false }
      }
      return { ...state, display: next }
    }

    case 'CLEAR':
      return { ...initialState }

    case 'OPERATOR': {
      if (isErrorStatus(state.status)) {
        return { ...initialState, firstOperand: 0, operator: action.operator }
      }

      if (state.operator !== null && !state.operandEntered) {
        return { ...state, operator: action.operator }
      }
      return {
        ...state,
        firstOperand: Number(state.display),
        operator: action.operator,
        operandEntered: false,
      }
    }

    case 'SUBMIT_START':
      return { ...state, status: 'loading', errorMessage: null }

    case 'SUBMIT_SUCCESS':
      return {
        ...state,
        status: 'idle',
        display: String(action.result),
        operator: null,
        firstOperand: null,
        operandEntered: false,
        errorMessage: null,
      }

    case 'SUBMIT_ERROR':
      return { ...state, status: 'error', errorMessage: action.message }

    case 'NETWORK_ERROR':
      return {
        ...state,
        status: 'networkError',
        errorMessage: NETWORK_ERROR_MESSAGE,
      }

    default:
      return state
  }
}
