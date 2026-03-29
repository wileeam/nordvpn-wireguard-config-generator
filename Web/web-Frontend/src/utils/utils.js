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

// Comprehensive ISO 3166-1 alpha-2 country code mapping for ALL countries worldwide
// This ensures any country that NordVPN adds in the future will automatically display flags
const countryCodeMap = {
  // A
  'Afghanistan': 'AF', 'Albania': 'AL', 'Algeria': 'DZ', 'Andorra': 'AD', 'Angola': 'AO',
  'Antigua And Barbuda': 'AG', 'Argentina': 'AR', 'Armenia': 'AM', 'Australia': 'AU', 'Austria': 'AT', 'Azerbaijan': 'AZ',
  // B
  'Bahamas': 'BS', 'Bahrain': 'BH', 'Bangladesh': 'BD', 'Barbados': 'BB', 'Belarus': 'BY',
  'Belgium': 'BE', 'Belize': 'BZ', 'Benin': 'BJ', 'Bhutan': 'BT', 'Bolivia': 'BO',
  'Bosnia And Herzegovina': 'BA', 'Botswana': 'BW', 'Brazil': 'BR', 'Brunei': 'BN', 'Bulgaria': 'BG',
  'Burkina Faso': 'BF', 'Burundi': 'BI',
  // C
  'Cambodia': 'KH', 'Cameroon': 'CM', 'Canada': 'CA', 'Cape Verde': 'CV',
  'Central African Republic': 'CF', 'Chad': 'TD', 'Chile': 'CL', 'China': 'CN', 'Colombia': 'CO',
  'Comoros': 'KM', 'Congo': 'CG', 'Costa Rica': 'CR', 'Croatia': 'HR', 'Cuba': 'CU',
  'Cyprus': 'CY', 'Czech Republic': 'CZ', 'Czechia': 'CZ',
  // D
  'Democratic Republic Of The Congo': 'CD', 'Denmark': 'DK', 'Djibouti': 'DJ', 'Dominica': 'DM', 'Dominican Republic': 'DO',
  // E
  'East Timor': 'TL', 'Ecuador': 'EC', 'Egypt': 'EG', 'El Salvador': 'SV',
  'Equatorial Guinea': 'GQ', 'Eritrea': 'ER', 'Estonia': 'EE', 'Ethiopia': 'ET', 'Eswatini': 'SZ',
  // F
  'Fiji': 'FJ', 'Finland': 'FI', 'France': 'FR',
  // G
  'Gabon': 'GA', 'Gambia': 'GM', 'Georgia': 'GE', 'Germany': 'DE', 'Ghana': 'GH',
  'Greece': 'GR', 'Grenada': 'GD', 'Guatemala': 'GT', 'Guinea': 'GN', 'Guinea-Bissau': 'GW', 'Guyana': 'GY',
  // H
  'Haiti': 'HT', 'Honduras': 'HN', 'Hong Kong': 'HK', 'Hungary': 'HU',
  // I
  'Iceland': 'IS', 'India': 'IN', 'Indonesia': 'ID', 'Iran': 'IR', 'Iraq': 'IQ',
  'Ireland': 'IE', 'Israel': 'IL', 'Italy': 'IT', 'Ivory Coast': 'CI',
  // J
  'Jamaica': 'JM', 'Japan': 'JP', 'Jordan': 'JO',
  // K
  'Kazakhstan': 'KZ', 'Kenya': 'KE', 'Kiribati': 'KI', 'Kosovo': 'XK', 'Kuwait': 'KW', 'Kyrgyzstan': 'KG',
  // L
  'Laos': 'LA', 'Latvia': 'LV', 'Lebanon': 'LB', 'Lesotho': 'LS', 'Liberia': 'LR',
  'Libya': 'LY', 'Liechtenstein': 'LI', 'Lithuania': 'LT', 'Luxembourg': 'LU',
  // M
  'Macao': 'MO', 'Madagascar': 'MG', 'Malawi': 'MW', 'Malaysia': 'MY', 'Maldives': 'MV',
  'Mali': 'ML', 'Malta': 'MT', 'Marshall Islands': 'MH', 'Mauritania': 'MR', 'Mauritius': 'MU',
  'Mexico': 'MX', 'Micronesia': 'FM', 'Moldova': 'MD', 'Monaco': 'MC', 'Mongolia': 'MN',
  'Montenegro': 'ME', 'Morocco': 'MA', 'Mozambique': 'MZ', 'Myanmar': 'MM',
  // N
  'Namibia': 'NA', 'Nauru': 'NR', 'Nepal': 'NP', 'Netherlands': 'NL', 'New Zealand': 'NZ',
  'Nicaragua': 'NI', 'Niger': 'NE', 'Nigeria': 'NG', 'North Korea': 'KP', 'North Macedonia': 'MK', 'Norway': 'NO',
  // O
  'Oman': 'OM',
  // P
  'Pakistan': 'PK', 'Palau': 'PW', 'Palestine': 'PS', 'Panama': 'PA', 'Papua New Guinea': 'PG',
  'Paraguay': 'PY', 'Peru': 'PE', 'Philippines': 'PH', 'Poland': 'PL', 'Portugal': 'PT',
  // Q
  'Qatar': 'QA',
  // R
  'Romania': 'RO', 'Russia': 'RU', 'Rwanda': 'RW',
  // S
  'Saint Kitts And Nevis': 'KN', 'Saint Lucia': 'LC', 'Saint Vincent And The Grenadines': 'VC',
  'Samoa': 'WS', 'San Marino': 'SM', 'Sao Tome And Principe': 'ST', 'Saudi Arabia': 'SA',
  'Senegal': 'SN', 'Serbia': 'RS', 'Seychelles': 'SC', 'Sierra Leone': 'SL', 'Singapore': 'SG',
  'Slovakia': 'SK', 'Slovenia': 'SI', 'Solomon Islands': 'SB', 'Somalia': 'SO', 'South Africa': 'ZA',
  'South Korea': 'KR', 'South Sudan': 'SS', 'Spain': 'ES', 'Sri Lanka': 'LK', 'Sudan': 'SD',
  'Suriname': 'SR', 'Sweden': 'SE', 'Switzerland': 'CH', 'Syria': 'SY',
  // T
  'Taiwan': 'TW', 'Tajikistan': 'TJ', 'Tanzania': 'TZ', 'Thailand': 'TH', 'Timor-Leste': 'TL',
  'Togo': 'TG', 'Tonga': 'TO', 'Trinidad And Tobago': 'TT', 'Tunisia': 'TN', 'Turkey': 'TR',
  'Turkmenistan': 'TM', 'Tuvalu': 'TV',
  // U
  'Uganda': 'UG', 'Ukraine': 'UA', 'United Arab Emirates': 'AE', 'United Kingdom': 'GB', 'United States': 'US',
  'Uruguay': 'UY', 'Uzbekistan': 'UZ',
  // V
  'Vanuatu': 'VU', 'Vatican City': 'VA', 'Venezuela': 'VE', 'Vietnam': 'VN',
  // Y
  'Yemen': 'YE',
  // Z
  'Zambia': 'ZM', 'Zimbabwe': 'ZW'
}

// Normalize country name for lookups (handles case, spacing, diacritics, and special characters)
const normalizeCountryName = name => {
  if (!name) return ''
  // Normalize Unicode, strip diacritics, and remove most punctuation
  const normalized = name
    .normalize('NFD') // separate base chars and diacritics
    .replace(/[\u0300-\u036f]/g, '') // strip combining diacritical marks
    .replace(/[^\p{L}\p{N}\s-]+/gu, ' ') // keep letters, numbers, spaces, and hyphens
    .trim()
    .replace(/\s+/g, ' ')

  // Convert to title case, including hyphenated words (e.g., "bosnia-herzegovina")
  const titleCased = normalized
    .split(' ')
    .map(word =>
      word
        .split('-')
        .map(part =>
          part
            ? part.charAt(0).toUpperCase() + part.slice(1).toLowerCase()
            : part
        )
        .join('-')
    )
    .join(' ')

  return titleCased
    .replace(/\bAnd\b/g, 'And')
    .replace(/\bOf\b/g, 'Of')
    .replace(/\bThe\b/g, 'The')
}

// Convert ISO 3166-1 alpha-2 country code to flag emoji
export const getCountryFlag = countryName => {
  const normalized = normalizeCountryName(countryName)
  const code = countryCodeMap[normalized]
  if (!code) {
    // If not found, log for debugging in development only (avoid noisy logs in production)
    if (typeof import.meta !== 'undefined' && import.meta.env && import.meta.env.DEV) {
      console.warn(`Country code not found for: "${countryName}" (normalized: "${normalized}")`)
    }
    return ''
  }
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