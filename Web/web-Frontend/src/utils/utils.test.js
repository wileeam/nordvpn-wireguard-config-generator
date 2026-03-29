import { describe, it, expect } from 'vitest'
import { formatName, sanitizeName, Validators } from './utils.js'

// ---------------------------------------------------------------------------
// formatName
// ---------------------------------------------------------------------------

describe('formatName', () => {
  it('returns empty string for falsy input', () => {
    expect(formatName('')).toBe('')
    expect(formatName(null)).toBe('')
    expect(formatName(undefined)).toBe('')
  })

  it('capitalizes a simple word', () => {
    expect(formatName('united')).toBe('United')
  })

  it('formats snake_case to Title Case', () => {
    expect(formatName('united_states')).toBe('United States')
  })

  it('handles multiple underscores', () => {
    expect(formatName('new_york_city')).toBe('New York City')
  })

  it('appends numeric suffix with space', () => {
    expect(formatName('us1')).toBe('Us 1')
  })

  it('handles server name with number', () => {
    expect(formatName('server42')).toBe('Server 42')
  })

  it('single word already capitalized', () => {
    expect(formatName('germany')).toBe('Germany')
  })

  it('handles already uppercase input', () => {
    // 'US' split as word 'US', no number
    expect(formatName('US')).toBe('US')
  })
})

// ---------------------------------------------------------------------------
// sanitizeName
// ---------------------------------------------------------------------------

describe('sanitizeName', () => {
  it('lowercases the input', () => {
    expect(sanitizeName('United States')).toBe('united states')
  })

  it('keeps spaces (only special chars replaced)', () => {
    expect(sanitizeName('New York')).toBe('new york')
  })

  it('replaces forbidden chars with underscores', () => {
    for (const ch of ['/', '\\', ':', '*', '?', '"', '<', '>', '|', '#']) {
      const result = sanitizeName(`foo${ch}bar`)
      expect(result).not.toContain(ch)
      expect(result).toContain('_')
    }
  })

  it('collapses multiple consecutive underscores', () => {
    const result = sanitizeName('foo//bar')
    expect(result).not.toMatch(/__+/)
  })

  it('trims leading and trailing underscores', () => {
    const result = sanitizeName('/test/')
    expect(result).not.toMatch(/^_|_$/)
  })

  it('already clean input unchanged', () => {
    expect(sanitizeName('already_clean')).toBe('already_clean')
  })

  it('returns empty string for empty input', () => {
    expect(sanitizeName('')).toBe('')
  })
})

// ---------------------------------------------------------------------------
// Validators.Key
// ---------------------------------------------------------------------------

describe('Validators.Key', () => {
  const validKey = 'QGWezOjmUUJkSGt4TqFWqHaGs5v0a2Atg050X+Ya6Vw='

  it('accepts a valid 44-char base64 key', () => {
    expect(Validators.Key.valid(validKey)).toBe(true)
  })

  it('accepts empty string (optional field)', () => {
    expect(Validators.Key.valid('')).toBe(true)
  })

  it('accepts falsy values as optional', () => {
    expect(Validators.Key.valid(null)).toBe(true)
    expect(Validators.Key.valid(undefined)).toBe(true)
  })

  it('rejects a key of wrong length', () => {
    expect(Validators.Key.valid('abc=')).toBe(false)
    expect(Validators.Key.valid(validKey + 'x')).toBe(false)
  })

  it('rejects key not ending with =', () => {
    expect(Validators.Key.valid(validKey.slice(0, 43) + 'x')).toBe(false)
  })

  it('has error message', () => {
    expect(typeof Validators.Key.err).toBe('string')
    expect(Validators.Key.err.length).toBeGreaterThan(0)
  })
})

// ---------------------------------------------------------------------------
// Validators.DNS
// ---------------------------------------------------------------------------

describe('Validators.DNS', () => {
  it('accepts a valid IPv4', () => {
    expect(Validators.DNS.valid('8.8.8.8')).toBe(true)
  })

  it('accepts comma-separated valid IPs', () => {
    expect(Validators.DNS.valid('8.8.8.8,1.1.1.1')).toBe(true)
  })

  it('accepts empty string (optional)', () => {
    expect(Validators.DNS.valid('')).toBe(true)
  })

  it('accepts falsy (optional)', () => {
    expect(Validators.DNS.valid(null)).toBe(true)
    expect(Validators.DNS.valid(undefined)).toBe(true)
  })

  it('rejects non-IP string', () => {
    expect(Validators.DNS.valid('notanip')).toBe(false)
  })

  it('rejects out-of-range octets', () => {
    expect(Validators.DNS.valid('256.0.0.1')).toBe(false)
  })

  it('rejects partial IP', () => {
    expect(Validators.DNS.valid('192.168.1')).toBe(false)
  })

  it('rejects invalid comma-separated list', () => {
    expect(Validators.DNS.valid('8.8.8.8,notanip')).toBe(false)
  })

  it('accepts IPs with spaces around commas', () => {
    expect(Validators.DNS.valid('8.8.8.8, 1.1.1.1')).toBe(true)
  })

  it('has error message', () => {
    expect(typeof Validators.DNS.err).toBe('string')
  })
})

// ---------------------------------------------------------------------------
// Validators.Keepalive
// ---------------------------------------------------------------------------

describe('Validators.Keepalive', () => {
  it('accepts value in valid range', () => {
    expect(Validators.Keepalive.valid(25)).toBe(true)
    expect(Validators.Keepalive.valid(60)).toBe(true)
  })

  it('accepts boundary values', () => {
    expect(Validators.Keepalive.valid(15)).toBe(true)
    expect(Validators.Keepalive.valid(120)).toBe(true)
  })

  it('accepts falsy (optional)', () => {
    expect(Validators.Keepalive.valid(null)).toBe(true)
    expect(Validators.Keepalive.valid(undefined)).toBe(true)
    expect(Validators.Keepalive.valid('')).toBe(true)
  })

  it('rejects below minimum', () => {
    expect(Validators.Keepalive.valid(14)).toBe(false)
  })

  it('rejects above maximum', () => {
    expect(Validators.Keepalive.valid(121)).toBe(false)
  })

  it('rejects non-numeric string', () => {
    expect(Validators.Keepalive.valid('abc')).toBe(false)
  })

  it('has min/max and error message', () => {
    expect(Validators.Keepalive.min).toBe(15)
    expect(Validators.Keepalive.max).toBe(120)
    expect(typeof Validators.Keepalive.err).toBe('string')
  })
})

// ---------------------------------------------------------------------------
// Validators.Token
// ---------------------------------------------------------------------------

describe('Validators.Token', () => {
  const validToken = 'a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2'

  it('accepts a valid 64 hex char token', () => {
    expect(Validators.Token.valid(validToken)).toBe(true)
  })

  it('accepts uppercase hex', () => {
    expect(Validators.Token.valid(validToken.toUpperCase())).toBe(true)
  })

  it('accepts falsy (optional)', () => {
    expect(Validators.Token.valid(null)).toBe(true)
    expect(Validators.Token.valid(undefined)).toBe(true)
    expect(Validators.Token.valid('')).toBe(true)
  })

  it('rejects token with wrong length', () => {
    expect(Validators.Token.valid(validToken.slice(0, 63))).toBe(false)
    expect(Validators.Token.valid(validToken + 'a')).toBe(false)
  })

  it('rejects non-hex characters', () => {
    expect(Validators.Token.valid('g'.repeat(64))).toBe(false)
  })

  it('clean() strips non-hex chars and limits to 64', () => {
    const dirty = '!@#$' + validToken + 'EXTRA'
    const cleaned = Validators.Token.clean(dirty)
    expect(cleaned.length).toBeLessThanOrEqual(64)
    expect(/^[a-f0-9]*$/.test(cleaned)).toBe(true)
  })

  it('clean() lowercases the token', () => {
    const cleaned = Validators.Token.clean(validToken.toUpperCase())
    expect(cleaned).toBe(validToken.toLowerCase())
  })

  it('clean() returns empty string for falsy input', () => {
    expect(Validators.Token.clean('')).toBe('')
    expect(Validators.Token.clean(null)).toBe('')
  })

  it('has error message', () => {
    expect(typeof Validators.Token.err).toBe('string')
  })
})
