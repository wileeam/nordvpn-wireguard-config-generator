const RX = {
  WORD: /^([a-z]+)(\d+)?$/i,
  NAME: /[\/\\:*?"<>|#]/g,
  MULTI: /_+/g,
  TRIM: /^_|_$/g,
  IPV4: /^(\d{1,3}\.){3}\d{1,3}$/,
  KEY: /^[A-Za-z0-9+/]{43}=$/,
  TOKEN: /^[a-f0-9]{64}$/i,
  HEX: /[^a-f0-9]/g
}

export const formatName = s => {
  if (!s) return ''
  return s.split('_').map(p => {
    const [, w, n] = p.match(RX.WORD) || [null, p]
    return w.charAt(0).toUpperCase() + w.slice(1) + (n ? ` ${n}` : '')
  }).join(' ')
}

export const sanitizeName = s => s.toLowerCase()
  .replace(RX.NAME, '_')
  .replace(RX.MULTI, '_')
  .replace(RX.TRIM, '')

// Map of country names (formatted) to ISO 3166-1 alpha-2 codes for flag emojis
const countryCodeMap = {
  'Albania': 'AL', 'Argentina': 'AR', 'Australia': 'AU', 'Austria': 'AT',
  'Belgium': 'BE', 'Bosnia And Herzegovina': 'BA', 'Brazil': 'BR', 'Bulgaria': 'BG',
  'Canada': 'CA', 'Chile': 'CL', 'Costa Rica': 'CR', 'Croatia': 'HR', 'Cyprus': 'CY', 'Czech Republic': 'CZ',
  'Denmark': 'DK',
  'Estonia': 'EE',
  'Finland': 'FI', 'France': 'FR',
  'Georgia': 'GE', 'Germany': 'DE', 'Greece': 'GR',
  'Hong Kong': 'HK', 'Hungary': 'HU',
  'Iceland': 'IS', 'India': 'IN', 'Indonesia': 'ID', 'Ireland': 'IE', 'Israel': 'IL', 'Italy': 'IT',
  'Japan': 'JP',
  'Latvia': 'LV', 'Lithuania': 'LT', 'Luxembourg': 'LU',
  'Malaysia': 'MY', 'Mexico': 'MX', 'Moldova': 'MD',
  'Netherlands': 'NL', 'New Zealand': 'NZ', 'North Macedonia': 'MK', 'Norway': 'NO',
  'Poland': 'PL', 'Portugal': 'PT',
  'Romania': 'RO',
  'Serbia': 'RS', 'Singapore': 'SG', 'Slovakia': 'SK', 'Slovenia': 'SI', 'South Africa': 'ZA', 'South Korea': 'KR', 'Spain': 'ES', 'Sweden': 'SE', 'Switzerland': 'CH',
  'Taiwan': 'TW', 'Thailand': 'TH', 'Turkey': 'TR',
  'Ukraine': 'UA', 'United Arab Emirates': 'AE', 'United Kingdom': 'GB', 'United States': 'US',
  'Vietnam': 'VN'
}

// Convert ISO 3166-1 alpha-2 country code to flag emoji
export const getCountryFlag = countryName => {
  const code = countryCodeMap[countryName]
  if (!code) return ''
  // Convert country code to flag emoji using regional indicator symbols
  return String.fromCodePoint(
    ...[...code].map(c => c.charCodeAt(0) + 127397)
  )
}

export const Validators = {
  Key: {
    valid: k => !k || RX.KEY.test(k),
    err: 'Invalid private key format'
  },
  DNS: {
    valid: d => !d || d.split(',').every(ip => {
      const t = ip.trim()
      return RX.IPV4.test(t) && t.split('.').every(n => {
        const i = parseInt(n, 10)
        return i >= 0 && i <= 255
      })
    }),
    err: 'Invalid IPv4 address'
  },
  Keepalive: {
    valid: v => !v || (!isNaN(v) && v >= 15 && v <= 120),
    min: 15,
    max: 120,
    err: 'Must be between 15 and 120'
  },
  Token: {
    valid: t => !t || RX.TOKEN.test(t),
    clean: t => t ? t.toLowerCase().replace(RX.HEX, '').slice(0, 64) : '',
    err: 'Token must be 64 hex characters'
  }
}