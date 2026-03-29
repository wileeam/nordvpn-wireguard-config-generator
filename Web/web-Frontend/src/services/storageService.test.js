import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { storage } from './storageService.js'

describe('storage', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  afterEach(() => {
    localStorage.clear()
  })

  // -------------------------------------------------------------------------
  // get
  // -------------------------------------------------------------------------

  describe('get', () => {
    it('returns null for missing key', () => {
      expect(storage.get('nonexistent')).toBeNull()
    })

    it('returns parsed value for existing key', () => {
      localStorage.setItem('mykey', JSON.stringify({ a: 1 }))
      expect(storage.get('mykey')).toEqual({ a: 1 })
    })

    it('returns primitive values correctly', () => {
      localStorage.setItem('num', JSON.stringify(42))
      expect(storage.get('num')).toBe(42)

      localStorage.setItem('str', JSON.stringify('hello'))
      expect(storage.get('str')).toBe('hello')

      localStorage.setItem('bool', JSON.stringify(true))
      expect(storage.get('bool')).toBe(true)
    })

    it('returns null on invalid JSON', () => {
      localStorage.setItem('bad', 'not-valid-json{{{')
      expect(storage.get('bad')).toBeNull()
    })

    it('returns null for null stored value', () => {
      localStorage.setItem('nullkey', 'null')
      // JSON.parse('null') === null, so this is null (but item exists)
      expect(storage.get('nullkey')).toBeNull()
    })
  })

  // -------------------------------------------------------------------------
  // set
  // -------------------------------------------------------------------------

  describe('set', () => {
    it('stores an object', () => {
      storage.set('obj', { x: 10, y: 20 })
      const raw = localStorage.getItem('obj')
      expect(JSON.parse(raw)).toEqual({ x: 10, y: 20 })
    })

    it('stores a string', () => {
      storage.set('str', 'hello world')
      expect(JSON.parse(localStorage.getItem('str'))).toBe('hello world')
    })

    it('stores a number', () => {
      storage.set('num', 99)
      expect(JSON.parse(localStorage.getItem('num'))).toBe(99)
    })

    it('stores a boolean', () => {
      storage.set('flag', false)
      expect(JSON.parse(localStorage.getItem('flag'))).toBe(false)
    })

    it('overwrites existing value', () => {
      storage.set('key', 'first')
      storage.set('key', 'second')
      expect(storage.get('key')).toBe('second')
    })
  })

  // -------------------------------------------------------------------------
  // clean
  // -------------------------------------------------------------------------

  describe('clean', () => {
    it('removes keys not in the allowed list', () => {
      localStorage.setItem('random_key', 'value')
      localStorage.setItem('another_key', 'value2')
      storage.clean()
      expect(localStorage.getItem('random_key')).toBeNull()
      expect(localStorage.getItem('another_key')).toBeNull()
    })

    it('preserves allowed keys', () => {
      storage.set('wg_gen_settings', { dns: '8.8.8.8' })
      storage.set('showIp', true)
      storage.clean()
      expect(storage.get('wg_gen_settings')).toEqual({ dns: '8.8.8.8' })
      expect(storage.get('showIp')).toBe(true)
    })

    it('leaves storage empty when no allowed keys exist', () => {
      localStorage.setItem('key1', 'val1')
      localStorage.setItem('key2', 'val2')
      storage.clean()
      expect(localStorage.length).toBe(0)
    })

    it('handles already-empty storage gracefully', () => {
      expect(() => storage.clean()).not.toThrow()
    })
  })

  // -------------------------------------------------------------------------
  // round-trip
  // -------------------------------------------------------------------------

  it('round-trips a complex settings object', () => {
    const settings = { dns: '8.8.8.8', endpoint: 'hostname', keepalive: 60 }
    storage.set('wg_gen_settings', settings)
    expect(storage.get('wg_gen_settings')).toEqual(settings)
  })
})
