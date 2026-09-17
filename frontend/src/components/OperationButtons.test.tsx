import { fireEvent, render, screen } from '@testing-library/react'
import type { ComponentProps } from 'react'
import { describe, expect, it, vi } from 'vitest'
import OperationButtons from './OperationButtons'

function renderButtons(
  overrides: Partial<ComponentProps<typeof OperationButtons>> = {},
) {
  const props = {
    operator: null,
    operandEntered: false,
    status: 'idle' as const,
    onOperator: vi.fn(),
    onEquals: vi.fn(),
    onSqrt: vi.fn(),
    onSquare: vi.fn(),
    ...overrides,
  }
  render(<OperationButtons {...props} />)
  return props
}

describe('OperationButtons — operator selection (RF-14)', () => {
  it.each([
    ['+', '+'],
    ['-', '-'],
    ['×', '×'],
    ['÷', '÷'],
    ['%', '%'],
    ['^', '^'],
  ])('pressing %s dispatches OPERATOR(%s)', (label, operator) => {
    const props = renderButtons()
    fireEvent.click(screen.getByRole('button', { name: label }))
    expect(props.onOperator).toHaveBeenCalledWith(operator)
  })
})

describe('OperationButtons — unary shortcuts (RF-15)', () => {
  it('pressing √x calls onSqrt', () => {
    const props = renderButtons()
    fireEvent.click(screen.getByRole('button', { name: '√x' }))
    expect(props.onSqrt).toHaveBeenCalled()
  })

  it('pressing x² calls onSquare', () => {
    const props = renderButtons()
    fireEvent.click(screen.getByRole('button', { name: 'x²' }))
    expect(props.onSquare).toHaveBeenCalled()
  })
})

describe('OperationButtons — equals disabled wiring (RF-18)', () => {
  it('is disabled when no operator is pending', () => {
    renderButtons({ operator: null, operandEntered: false })
    expect(screen.getByRole('button', { name: '=' })).toBeDisabled()
  })

  it('is disabled when an operator is pending but no second-operand digit was entered', () => {
    renderButtons({ operator: '+', operandEntered: false })
    expect(screen.getByRole('button', { name: '=' })).toBeDisabled()
  })

  it('is enabled once an operator is pending and a second-operand digit was entered', () => {
    renderButtons({ operator: '+', operandEntered: true })
    expect(screen.getByRole('button', { name: '=' })).toBeEnabled()
  })

  it('pressing an enabled equals calls onEquals', () => {
    const props = renderButtons({ operator: '+', operandEntered: true })
    fireEvent.click(screen.getByRole('button', { name: '=' }))
    expect(props.onEquals).toHaveBeenCalled()
  })
})

describe('OperationButtons — in-flight guard (RF-26)', () => {
  it('disables equals and both unary shortcuts while a request is in flight, even when a submission would otherwise be valid', () => {
    renderButtons({ operator: '+', operandEntered: true, status: 'loading' })
    expect(screen.getByRole('button', { name: '=' })).toBeDisabled()
    expect(screen.getByRole('button', { name: '√x' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'x²' })).toBeDisabled()
  })

  it('re-enables equals and unary shortcuts once the response has landed (idle)', () => {
    renderButtons({ operator: '+', operandEntered: true, status: 'idle' })
    expect(screen.getByRole('button', { name: '=' })).toBeEnabled()
    expect(screen.getByRole('button', { name: '√x' })).toBeEnabled()
    expect(screen.getByRole('button', { name: 'x²' })).toBeEnabled()
  })

  it('re-enables unary shortcuts after an error response too', () => {
    renderButtons({ status: 'error' })
    expect(screen.getByRole('button', { name: '√x' })).toBeEnabled()
    expect(screen.getByRole('button', { name: 'x²' })).toBeEnabled()
  })

  it('does not disable the binary operator buttons while loading (RF-26 only names equals and unary shortcuts)', () => {
    renderButtons({ status: 'loading' })
    expect(screen.getByRole('button', { name: '+' })).toBeEnabled()
  })
})
