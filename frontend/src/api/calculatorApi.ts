import type { ApiResult, ErrorCode } from '../types/calculator'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'

interface ErrorBody {
  error: { code: ErrorCode; message: string }
}

interface ResultBody {
  result: number
}

async function postJSON(
  path: string,
  body: Record<string, number>,
): Promise<ApiResult> {
  let response: Response
  try {
    response = await fetch(`${API_BASE_URL}${path}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    })
  } catch {
    return { ok: false, kind: 'network' }
  }

  if (!response.ok) {
    const { error } = (await response.json()) as ErrorBody
    return { ok: false, kind: 'domain', code: error.code, message: error.message }
  }

  const { result } = (await response.json()) as ResultBody
  return { ok: true, result }
}

export function add(a: number, b: number): Promise<ApiResult> {
  return postJSON('/add', { a, b })
}

export function subtract(a: number, b: number): Promise<ApiResult> {
  return postJSON('/subtract', { a, b })
}

export function multiply(a: number, b: number): Promise<ApiResult> {
  return postJSON('/multiply', { a, b })
}

export function divide(a: number, b: number): Promise<ApiResult> {
  return postJSON('/divide', { a, b })
}

export function power(base: number, exponent: number): Promise<ApiResult> {
  return postJSON('/power', { base, exponent })
}

export function sqrt(value: number): Promise<ApiResult> {
  return postJSON('/sqrt', { value })
}

export function percentage(value: number, percentage: number): Promise<ApiResult> {
  return postJSON('/percentage', { value, percentage })
}
