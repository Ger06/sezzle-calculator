import { useEffect, useReducer } from 'react'
import * as calculatorApi from '../api/calculatorApi'
import { calculatorReducer, initialState } from '../state/calculatorReducer'
import type { ApiResult, Operator } from '../types/calculator'
import './Calculator.css'
import Display from './Display'
import Keypad from './Keypad'
import OperationButtons from './OperationButtons'

function callOperation(operator: Operator, a: number, b: number): Promise<ApiResult> {
  switch (operator) {
    case '+':
      return calculatorApi.add(a, b)
    case '-':
      return calculatorApi.subtract(a, b)
    case '×':
      return calculatorApi.multiply(a, b)
    case '÷':
      return calculatorApi.divide(a, b)
    case '%':
      return calculatorApi.percentage(a, b)
    case '^':
      return calculatorApi.power(a, b)
  }
}

const KEY_TO_OPERATOR: Partial<Record<string, Operator>> = {
  '+': '+',
  '-': '-',
  '*': '×',
  '/': '÷',
  '^': '^',
  '%': '%',
}

function Calculator() {
  const [state, dispatch] = useReducer(calculatorReducer, initialState)

  function applyResult(result: ApiResult): boolean {
    if (result.ok) {
      dispatch({ type: 'SUBMIT_SUCCESS', result: result.result })
      return true
    }
    if (result.kind === 'domain') {
      dispatch({ type: 'SUBMIT_ERROR', message: result.message })
      return false
    }
    dispatch({ type: 'NETWORK_ERROR' })
    return false
  }

  async function submit(operator: Operator, a: number, b: number): Promise<boolean> {
    dispatch({ type: 'SUBMIT_START' })
    const result = await callOperation(operator, a, b)
    return applyResult(result)
  }

  async function handleEquals() {

    if (state.status === 'loading' || state.operator === null || !state.operandEntered) {
      return
    }
    await submit(state.operator, state.firstOperand ?? 0, Number(state.display))
  }

  async function handleOperator(operator: Operator) {

    if (state.status === 'loading') return

    if (state.operator !== null && state.operandEntered) {

      const succeeded = await submit(
        state.operator,
        state.firstOperand ?? 0,
        Number(state.display),
      )
      if (!succeeded) return
      dispatch({ type: 'OPERATOR', operator })
      return
    }

    dispatch({ type: 'OPERATOR', operator })
  }

  async function handleUnary(call: (value: number) => Promise<ApiResult>) {
    if (state.status === 'loading') return
    dispatch({ type: 'SUBMIT_START' })
    const result = await call(Number(state.display))
    applyResult(result)
  }

  const handleSqrt = () => {
    void handleUnary((value) => calculatorApi.sqrt(value))
  }
  const handleSquare = () => {
    void handleUnary((value) => calculatorApi.power(value, 2))
  }

  useEffect(() => {
    function handleKeyDown(event: KeyboardEvent) {
      if (event.target instanceof HTMLButtonElement) return

      if (event.key >= '0' && event.key <= '9') {
        dispatch({ type: 'DIGIT', digit: event.key })
        return
      }

      const operator = KEY_TO_OPERATOR[event.key]
      if (operator) {
        void handleOperator(operator)
        return
      }

      switch (event.key) {
        case 'Enter':
          event.preventDefault()
          void handleEquals()
          break
        case 'Escape':
          dispatch({ type: 'CLEAR' })
          break
        case 'Backspace':
          event.preventDefault()
          dispatch({ type: 'BACKSPACE' })
          break
        default:
          break
      }
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  })

  return (
    <div className="calculator">
      <Display
        display={state.display}
        status={state.status}
        errorMessage={state.errorMessage}
      />
      <div className="calculator-grid">
        <Keypad dispatch={dispatch} />
        <OperationButtons
          operator={state.operator}
          operandEntered={state.operandEntered}
          status={state.status}
          onOperator={(operator) => void handleOperator(operator)}
          onEquals={() => void handleEquals()}
          onSqrt={handleSqrt}
          onSquare={handleSquare}
        />
      </div>
    </div>
  )
}

export default Calculator
