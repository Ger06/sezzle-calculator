import type { CSSProperties, Dispatch } from 'react'
import type { Action } from '../types/calculator'
import './buttonBase.css'
import './Keypad.css'

interface KeypadProps {
  dispatch: Dispatch<Action>
}

const DIGIT_AREAS: Record<string, string> = {
  '7': 'seven',
  '8': 'eight',
  '9': 'nine',
  '4': 'four',
  '5': 'five',
  '6': 'six',
  '1': 'one',
  '2': 'two',
  '3': 'three',
}

const DIGIT_ROWS: readonly (readonly string[])[] = [
  ['7', '8', '9'],
  ['4', '5', '6'],
  ['1', '2', '3'],
]

function area(name: string): CSSProperties {
  return { gridArea: name }
}

function Keypad({ dispatch }: KeypadProps) {
  return (
    <>
      <button
        type="button"
        className="key key-clear"
        style={area('clear')}
        onClick={() => dispatch({ type: 'CLEAR' })}
      >
        AC
      </button>
      <button
        type="button"
        className="key key-backspace"
        style={area('back')}
        aria-label="Backspace"
        onClick={() => dispatch({ type: 'BACKSPACE' })}
      >
        ⌫
      </button>
      <button
        type="button"
        className="key key-sign"
        style={area('sign')}
        aria-label="Toggle sign"
        onClick={() => dispatch({ type: 'SIGN_TOGGLE' })}
      >
        ±
      </button>

      {DIGIT_ROWS.flatMap((row) =>
        row.map((digit) => (
          <button
            key={digit}
            type="button"
            className="key key-digit"
            style={area(DIGIT_AREAS[digit])}
            onClick={() => dispatch({ type: 'DIGIT', digit })}
          >
            {digit}
          </button>
        )),
      )}

      <button
        type="button"
        className="key key-digit"
        style={area('zero')}
        onClick={() => dispatch({ type: 'DIGIT', digit: '0' })}
      >
        0
      </button>
      <button
        type="button"
        className="key key-digit"
        style={area('dot')}
        onClick={() => dispatch({ type: 'DECIMAL' })}
      >
        .
      </button>
    </>
  )
}

export default Keypad
