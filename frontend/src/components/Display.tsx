import type { Status } from '../types/calculator'
import './Display.css'

interface DisplayProps {
  display: string
  status: Status
  errorMessage: string | null
}

function Display({ display, status, errorMessage }: DisplayProps) {
  const showError = status === 'error' || status === 'networkError'

  return (
    <div className="display" role="status" aria-live="polite">
      {showError ? (
        <span className="display-error">{errorMessage}</span>
      ) : (
        <span className="display-value">{display}</span>
      )}
    </div>
  )
}

export default Display
