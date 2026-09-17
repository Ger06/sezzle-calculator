import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import {
  add,
  divide,
  multiply,
  percentage,
  power,
  sqrt,
  subtract,
} from './calculatorApi'

function mockResponse(ok: boolean, body: unknown): Response {
  return { ok, json: () => Promise.resolve(body) } as Response
}

let fetchMock: ReturnType<typeof vi.fn>

beforeEach(() => {
  fetchMock = vi.fn()
  vi.stubGlobal('fetch', fetchMock)
})

afterEach(() => {
  vi.unstubAllGlobals()
})

function lastRequest() {
  const [url, init] = fetchMock.mock.calls[0] as [string, RequestInit]
  return { url, body: JSON.parse(init.body as string) as unknown }
}

describe('calculatorApi request field names (Data contract)', () => {
  it('add sends {a, b} to /add', async () => {
    fetchMock.mockResolvedValueOnce(mockResponse(true, { result: 5 }))
    await add(2, 3)
    const { url, body } = lastRequest()
    expect(url).toMatch(/\/add$/)
    expect(body).toEqual({ a: 2, b: 3 })
  })

  it('subtract sends {a, b} to /subtract', async () => {
    fetchMock.mockResolvedValueOnce(mockResponse(true, { result: 2 }))
    await subtract(5, 3)
    const { url, body } = lastRequest()
    expect(url).toMatch(/\/subtract$/)
    expect(body).toEqual({ a: 5, b: 3 })
  })

  it('multiply sends {a, b} to /multiply', async () => {
    fetchMock.mockResolvedValueOnce(mockResponse(true, { result: 10 }))
    await multiply(4, 2.5)
    const { url, body } = lastRequest()
    expect(url).toMatch(/\/multiply$/)
    expect(body).toEqual({ a: 4, b: 2.5 })
  })

  it('divide sends {a, b} to /divide', async () => {
    fetchMock.mockResolvedValueOnce(mockResponse(true, { result: 2.5 }))
    await divide(10, 4)
    const { url, body } = lastRequest()
    expect(url).toMatch(/\/divide$/)
    expect(body).toEqual({ a: 10, b: 4 })
  })

  it('power sends {base, exponent} to /power', async () => {
    fetchMock.mockResolvedValueOnce(mockResponse(true, { result: 1024 }))
    await power(2, 10)
    const { url, body } = lastRequest()
    expect(url).toMatch(/\/power$/)
    expect(body).toEqual({ base: 2, exponent: 10 })
  })

  it('sqrt sends {value} to /sqrt', async () => {
    fetchMock.mockResolvedValueOnce(mockResponse(true, { result: 4 }))
    await sqrt(16)
    const { url, body } = lastRequest()
    expect(url).toMatch(/\/sqrt$/)
    expect(body).toEqual({ value: 16 })
  })

  it('percentage sends {value, percentage} to /percentage', async () => {
    fetchMock.mockResolvedValueOnce(mockResponse(true, { result: 30 }))
    await percentage(200, 15)
    const { url, body } = lastRequest()
    expect(url).toMatch(/\/percentage$/)
    expect(body).toEqual({ value: 200, percentage: 15 })
  })
})

describe('calculatorApi outcome discrimination', () => {
  it('discriminates a successful response', async () => {
    fetchMock.mockResolvedValueOnce(mockResponse(true, { result: 5 }))
    const outcome = await add(2, 3)
    expect(outcome).toEqual({ ok: true, result: 5 })
  })

  it('discriminates a domain/validation error response despite fetch resolving "successfully" (the response.ok gotcha, docs/decisions.md #10)', async () => {
    fetchMock.mockResolvedValueOnce(
      mockResponse(false, {
        error: { code: 'DIVISION_BY_ZERO', message: 'cannot divide by zero' },
      }),
    )
    const outcome = await divide(1, 0)
    expect(outcome).toEqual({
      ok: false,
      kind: 'domain',
      code: 'DIVISION_BY_ZERO',
      message: 'cannot divide by zero',
    })
  })

  it('discriminates an INVALID_INPUT validation error the same way as a domain error', async () => {
    fetchMock.mockResolvedValueOnce(
      mockResponse(false, {
        error: { code: 'INVALID_INPUT', message: 'b is required' },
      }),
    )
    const outcome = await add(2, Number.NaN)
    expect(outcome).toEqual({
      ok: false,
      kind: 'domain',
      code: 'INVALID_INPUT',
      message: 'b is required',
    })
  })

  it('discriminates a network-level failure (fetch rejects — no HTTP response at all)', async () => {
    fetchMock.mockRejectedValueOnce(new TypeError('Failed to fetch'))
    const outcome = await add(2, 3)
    expect(outcome).toEqual({ ok: false, kind: 'network' })
  })
})
