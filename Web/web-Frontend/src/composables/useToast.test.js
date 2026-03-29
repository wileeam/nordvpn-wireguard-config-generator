import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useToast } from './useToast.js'

describe('useToast', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('initializes with null toast', () => {
    const { toast } = useToast()
    expect(toast.value).toBeNull()
  })

  it('show() sets toast with success type by default', () => {
    const { toast, show } = useToast()
    show('Hello World')
    expect(toast.value).not.toBeNull()
    expect(toast.value.message).toBe('Hello World')
    expect(toast.value.type).toBe('success')
  })

  it('show() sets toast with error type', () => {
    const { toast, show } = useToast()
    show('Something failed', 'error')
    expect(toast.value.type).toBe('error')
  })

  it('show() defaults unknown type to success', () => {
    const { toast, show } = useToast()
    show('msg', 'unknown_type')
    expect(toast.value.type).toBe('success')
  })

  it('show() does nothing for empty message', () => {
    const { toast, show } = useToast()
    show('')
    expect(toast.value).toBeNull()
  })

  it('show() does nothing for null message', () => {
    const { toast, show } = useToast()
    show(null)
    expect(toast.value).toBeNull()
  })

  it('show() converts Error objects to message', () => {
    const { toast, show } = useToast()
    show(new Error('Something went wrong'))
    expect(toast.value.message).toBe('Something went wrong')
  })

  it('toast auto-dismisses after 2 seconds', () => {
    const { toast, show } = useToast()
    show('Auto dismiss me')
    expect(toast.value).not.toBeNull()
    vi.advanceTimersByTime(2000)
    expect(toast.value).toBeNull()
  })

  it('toast is still visible before 2 seconds', () => {
    const { toast, show } = useToast()
    show('Still visible')
    vi.advanceTimersByTime(1999)
    expect(toast.value).not.toBeNull()
  })

  it('calling show() again resets the timer', () => {
    const { toast, show } = useToast()
    show('First message')
    vi.advanceTimersByTime(1500)
    show('Second message')
    vi.advanceTimersByTime(1500)
    // Still visible because timer was reset
    expect(toast.value).not.toBeNull()
    expect(toast.value.message).toBe('Second message')
  })

  it('show() truncates message to 100 characters', () => {
    const { toast, show } = useToast()
    const long = 'x'.repeat(200)
    show(long)
    expect(toast.value.message.length).toBeLessThanOrEqual(100)
  })

  it('show() uses only the first line of multi-line message', () => {
    const { toast, show } = useToast()
    show('Line 1\nLine 2\nLine 3')
    expect(toast.value.message).toBe('Line 1')
  })
})
