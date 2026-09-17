import type { CSSProperties } from 'react'
import type { Operator, Status } from '../types/calculator'
import './buttonBase.css'
import './OperationButtons.css'

interface OperationButtonsProps {
  operator: Operator | null
  operandEntered: boolean
  status: Status
  onOperator: (operator: Operator) => void
  onEquals: () => void
  onSqrt: () => void
  onSquare: () => void
}

const OPERATOR_AREAS: Record<Operator, string> = {
  '÷': 'divide',
  '×': 'multiply',
  '-': 'subtract',
  '+': 'add',
  '%': 'percent',
  '^': 'power',
}

const BINARY_OPERATORS: readonly Operator[] = ['÷', '×', '-', '+', '%', '^']

function area(name: string): CSSProperties {
  return { gridArea: name }
}

function OperationButtons({
  operator,
  operandEntered,
  status,
  onOperator,
  onEquals,
  onSqrt,
  onSquare,
}: OperationButtonsProps) {
  const isLoading = status === 'loading'
  const canSubmit = operator !== null && operandEntered
  const equalsDisabled = isLoading || !canSubmit
  const unaryDisabled = isLoading

  return (
    <>
      <button
        type="button"
        className="key key-unary"
        style={area('sqrt')}
        onClick={onSqrt}
        disabled={unaryDisabled}
      >
        √x
      </button>
      <button
        type="button"
        className="key key-unary"
        style={area('square')}
        onClick={onSquare}
        disabled={unaryDisabled}
      >
        x²
      </button>
      {BINARY_OPERATORS.map((op) => (
        <button
          key={op}
          type="button"
          className="key key-operator"
          style={area(OPERATOR_AREAS[op])}
          onClick={() => onOperator(op)}
        >
          {op}
        </button>
      ))}
      <button
        type="button"
        className="key key-equals"
        style={area('equals')}
        onClick={onEquals}
        disabled={equalsDisabled}
      >
        =
      </button>
    </>
  )
}

export default OperationButtons
