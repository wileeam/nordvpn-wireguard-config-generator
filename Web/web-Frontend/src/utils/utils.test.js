import { describe, it, expect } from 'vitest'
import { formatName, sanitizeName, Validators } from '../utils'

describe('formatName', () => {
  it('should format snake_case to Title Case', () => {
    expect(formatName('united_states')).toBe('United States')
  })

  it('should handle snake_case with numbers', () => {
    expect(formatName('server_1234')).toBe('Server 1234')
  })

  it('should handle multiple words with numbers', () => {
    expect(formatName('new_york_99')).toBe('New York 99')
  })

  it('should handle empty string', () => {
    expect(formatName('')).toBe('')
  })

  it('should handle null/undefined', () => {
    expect(formatName(null)).toBe('')
    expect(formatName(undefined)).toBe('')
  })

  it('should handle single word', () => {
    expect(formatName('tokyo')).toBe('Tokyo')
  })

  it('should handle uppercase input', () => {
    expect(formatName('UNITED_STATES')).toBe('UNITED STATES')
  })

  it('should handle mixed case', () => {
    expect(formatName('UnItEd_StAtEs')).toBe('UnItEd StAtEs')
  })
})

describe('sanitizeName', () => {
  it('should convert to lowercase', () => {
    expect(sanitizeName('HELLO')).toBe('hello')
  })

  it('should replace invalid characters with underscore', () => {
    expect(sanitizeName('hello/world')).toBe('hello_world')
    expect(sanitizeName('hello\\world')).toBe('hello_world')
    expect(sanitizeName('hello:world')).toBe('hello_world')
    expect(sanitizeName('hello*world')).toBe('hello_world')
    expect(sanitizeName('hello?world')).toBe('hello_world')
    expect(sanitizeName('hello"world')).toBe('hello_world')
    expect(sanitizeName('hello<world')).toBe('hello_world')
    expect(sanitizeName('hello>world')).toBe('hello_world')
    expect(sanitizeName('hello|world')).toBe('hello_world')
    expect(sanitizeName('hello#world')).toBe('hello_world')
  })

  it('should replace multiple underscores with single underscore', () => {
    expect(sanitizeName('hello___world')).toBe('hello_world')
  })

  it('should trim leading and trailing underscores', () => {
    expect(sanitizeName('_hello_')).toBe('hello')
    expect(sanitizeName('___hello___')).toBe('hello')
  })

  it('should handle complex sanitization', () => {
    expect(sanitizeName('Test/File*Name<2024>.txt')).toBe('test_file_name_2024_txt')
  })

  it('should handle empty string', () => {
    expect(sanitizeName('')).toBe('')
  })

  it('should handle valid names unchanged (except lowercase)', () => {
    expect(sanitizeName('hello_world_123')).toBe('hello_world_123')
  })
})

describe('Validators.Key', () => {
  it('should accept valid WireGuard keys', () => {
    expect(Validators.Key.valid('ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno=')).toBe(true)
    expect(Validators.Key.valid('ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+/abcde=')).toBe(true)
    expect(Validators.Key.valid('1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefg=')).toBe(true)
  })

  it('should accept empty string', () => {
    expect(Validators.Key.valid('')).toBe(true)
  })

  it('should reject invalid keys', () => {
    expect(Validators.Key.valid('too_short=')).toBe(false)
    expect(Validators.Key.valid('ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnop')).toBe(false) // no =
    expect(Validators.Key.valid('ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno==')).toBe(false) // double =
    expect(Validators.Key.valid('ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmn@=')).toBe(false) // invalid char
  })

  it('should have error message', () => {
    expect(Validators.Key.err).toBe('Invalid private key format')
  })
})

describe('Validators.DNS', () => {
  it('should accept valid single IPv4', () => {
    expect(Validators.DNS.valid('1.1.1.1')).toBe(true)
    expect(Validators.DNS.valid('8.8.8.8')).toBe(true)
    expect(Validators.DNS.valid('192.168.1.1')).toBe(true)
    expect(Validators.DNS.valid('255.255.255.255')).toBe(true)
    expect(Validators.DNS.valid('0.0.0.0')).toBe(true)
  })

  it('should accept multiple comma-separated IPs', () => {
    expect(Validators.DNS.valid('1.1.1.1,8.8.8.8')).toBe(true)
    expect(Validators.DNS.valid('1.1.1.1, 8.8.8.8, 9.9.9.9')).toBe(true)
  })

  it('should accept empty string', () => {
    expect(Validators.DNS.valid('')).toBe(true)
  })

  it('should reject invalid IPv4', () => {
    expect(Validators.DNS.valid('256.1.1.1')).toBe(false)
    expect(Validators.DNS.valid('1.1.1.256')).toBe(false)
    expect(Validators.DNS.valid('1.1.1')).toBe(false)
    expect(Validators.DNS.valid('1.1.1.1.1')).toBe(false)
    expect(Validators.DNS.valid('abc.def.ghi.jkl')).toBe(false)
    expect(Validators.DNS.valid('192.168.-1.1')).toBe(false)
  })

  it('should reject if any IP in list is invalid', () => {
    expect(Validators.DNS.valid('1.1.1.1,256.1.1.1')).toBe(false)
    expect(Validators.DNS.valid('8.8.8.8,invalid')).toBe(false)
  })

  it('should handle whitespace around IPs', () => {
    expect(Validators.DNS.valid(' 1.1.1.1 , 8.8.8.8 ')).toBe(true)
  })

  it('should have error message', () => {
    expect(Validators.DNS.err).toBe('Invalid IPv4 address')
  })
})

describe('Validators.Keepalive', () => {
  it('should accept valid keepalive values', () => {
    expect(Validators.Keepalive.valid(15)).toBe(true)
    expect(Validators.Keepalive.valid(25)).toBe(true)
    expect(Validators.Keepalive.valid(60)).toBe(true)
    expect(Validators.Keepalive.valid(120)).toBe(true)
  })

  it('should accept empty value', () => {
    expect(Validators.Keepalive.valid('')).toBe(true)
    expect(Validators.Keepalive.valid(null)).toBe(true)
    expect(Validators.Keepalive.valid(undefined)).toBe(true)
  })

  it('should reject values below 15', () => {
    expect(Validators.Keepalive.valid(14)).toBe(false)
    expect(Validators.Keepalive.valid(0)).toBe(false)
    expect(Validators.Keepalive.valid(-1)).toBe(false)
  })

  it('should reject values above 120', () => {
    expect(Validators.Keepalive.valid(121)).toBe(false)
    expect(Validators.Keepalive.valid(200)).toBe(false)
  })

  it('should reject non-numeric values', () => {
    expect(Validators.Keepalive.valid('abc')).toBe(false)
    expect(Validators.Keepalive.valid('25a')).toBe(false)
  })

  it('should have min and max values', () => {
    expect(Validators.Keepalive.min).toBe(15)
    expect(Validators.Keepalive.max).toBe(120)
  })

  it('should have error message', () => {
    expect(Validators.Keepalive.err).toBe('Must be between 15 and 120')
  })
})

describe('Validators.Token', () => {
  it('should accept valid 64-char hex tokens', () => {
    expect(Validators.Token.valid('abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789')).toBe(true)
    expect(Validators.Token.valid('ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789')).toBe(true)
    expect(Validators.Token.valid('0123456789abcdefABCDEF0123456789abcdefABCDEF0123456789abcdefABCD')).toBe(true)
  })

  it('should accept empty string', () => {
    expect(Validators.Token.valid('')).toBe(true)
  })

  it('should reject tokens with wrong length', () => {
    expect(Validators.Token.valid('abc123')).toBe(false)
    expect(Validators.Token.valid('abcdef0123456789abcdef0123456789abcdef0123456789abcdef012345678')).toBe(false) // 63 chars
    expect(Validators.Token.valid('abcdef0123456789abcdef0123456789abcdef0123456789abcdef01234567890')).toBe(false) // 65 chars
  })

  it('should reject non-hex characters', () => {
    expect(Validators.Token.valid('gbcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789')).toBe(false)
    expect(Validators.Token.valid('abcdef0123456789abcdef0123456789abcdef0123456789abcdef012345678z')).toBe(false)
    expect(Validators.Token.valid('abcdef0123456789-bcdef0123456789abcdef0123456789abcdef0123456789')).toBe(false)
  })

  it('should clean tokens correctly', () => {
    expect(Validators.Token.clean('ABCDEF0123456789')).toBe('abcdef0123456789')
    expect(Validators.Token.clean('abc-def-012')).toBe('abcdef012')
    expect(Validators.Token.clean('GHIJKLMNOP')).toBe('') // non-hex chars removed
    expect(Validators.Token.clean('abcdef0123456789abcdef0123456789abcdef0123456789abcdef01234567890000')).toBe('abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789') // truncated to 64
    expect(Validators.Token.clean('')).toBe('')
  })

  it('should have error message', () => {
    expect(Validators.Token.err).toBe('Token must be 64 hex characters')
  })
})
