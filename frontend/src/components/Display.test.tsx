import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import Display from './Display'

describe('Display', () => {
  it('renders the current value when idle (RF-16)', () => {
    render(<Display display="42" status="idle" errorMessage={null} />)
    expect(screen.getByText('42')).toBeInTheDocument()
  })

  it('keeps showing the value while a request is in flight (RF-26 does not hide it)', () => {
    render(<Display display="42" status="loading" errorMessage={null} />)
    expect(screen.getByText('42')).toBeInTheDocument()
  })

  it('renders the error message instead of the value on a validation/domain error (RF-17)', () => {
    render(<Display display="0" status="error" errorMessage="cannot divide by zero" />)
    expect(screen.getByText('cannot divide by zero')).toBeInTheDocument()
    expect(screen.queryByText('0')).not.toBeInTheDocument()
  })

  it('renders a distinct connection-error message on a network failure (RF-25)', () => {
    render(
      <Display
        display="0"
        status="networkError"
        errorMessage="Unable to reach the server. Please check your connection and try again."
      />,
    )
    expect(
      screen.getByText(
        'Unable to reach the server. Please check your connection and try again.',
      ),
    ).toBeInTheDocument()
    expect(screen.queryByText('0')).not.toBeInTheDocument()
  })
})
