export type Operator = '+' | '-' | '×' | '÷' | '%' | '^'

export type Status = 'idle' | 'loading' | 'error' | 'networkError'

export interface State {
  display: string
  operandEntered: boolean
  firstOperand: number | null
  operator: Operator | null
  status: Status
  errorMessage: string | null
}

export type Action =
  | { type: 'DIGIT'; digit: string }
  | { type: 'DECIMAL' }
  | { type: 'SIGN_TOGGLE' }
  | { type: 'BACKSPACE' }
  | { type: 'CLEAR' }
  | { type: 'OPERATOR'; operator: Operator }
  | { type: 'SUBMIT_START' }
  | { type: 'SUBMIT_SUCCESS'; result: number }
  | { type: 'SUBMIT_ERROR'; message: string }
  | { type: 'NETWORK_ERROR' }

export type ErrorCode =
  | 'INVALID_INPUT'
  | 'DIVISION_BY_ZERO'
  | 'NEGATIVE_SQRT'
  | 'NON_FINITE_RESULT'

export type ApiResult =
  | { ok: true; result: number }
  | { ok: false; kind: 'domain'; code: ErrorCode; message: string }
  | { ok: false; kind: 'network' }
