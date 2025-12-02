/**
 * Error Messages Mapper
 * Maps backend error codes and messages to Spanish user-friendly messages
 */

export interface ErrorMapping {
  code: string
  message: string
}

// Error code to Spanish message mapping
export const ERROR_MESSAGES: Record<string, string> = {
  // Request validation errors
  'request_validation_failed': 'Error de validación en la solicitud',
  'invalid_request_format': 'Formato de solicitud inválido',
  
  // Use case execution errors
  'use_case_execution_failed': 'Error al procesar la solicitud',
  
  // Not found errors
  'not_found': 'No se encontró el recurso solicitado',
  
  // Match announcement specific errors
  'team_not_found': 'El equipo especificado no existe',
  'invalid_start_time_format': 'Formato de hora de inicio inválido',
  'invalid_end_time_format': 'Formato de hora de fin inválido',
  'invalid_day_format': 'Formato de fecha inválido',
  'invalid_category_format': 'Formato de categoría inválido',
  'invalid_category_range_type': 'Tipo de rango de categoría inválido',
  'day_cannot_be_in_the_past': 'La fecha del partido no puede ser en el pasado',
  'end_time_cannot_be_before_start_time': 'La hora de fin no puede ser anterior a la hora de inicio',
  'location_must_have_country_province_and_locality': 'La ubicación debe tener país, provincia y localidad',
  'team_name_cannot_be_empty': 'El nombre del equipo no puede estar vacío',
  'sport_cannot_be_empty': 'El deporte no puede estar vacío',
  'time_slot_cannot_be_empty': 'El horario no puede estar vacío',
  'invalid_status': 'Estado inválido',
  'created_at_cannot_be_empty': 'La fecha de creación no puede estar vacía',
}

// Common error message patterns to translate
export const ERROR_PATTERNS: Array<{ pattern: RegExp; message: string }> = [
  {
    pattern: /team '(.+)' for sport '(.+)' does not exist/i,
    message: 'El equipo "$1" para el deporte "$2" no existe',
  },
  {
    pattern: /invalid (.+) format/i,
    message: 'Formato de $1 inválido',
  },
  {
    pattern: /(.+) cannot be empty/i,
    message: '$1 no puede estar vacío',
  },
  {
    pattern: /(.+) is required/i,
    message: '$1 es obligatorio',
  },
  {
    pattern: /unable to parse datetime: (.+)/i,
    message: 'No se pudo parsear la fecha/hora: $1',
  },
]

/**
 * Translates backend error to Spanish user-friendly message
 * @param errorCode - Error code from backend
 * @param errorMessage - Original error message from backend
 * @returns Spanish translated message
 */
export function translateError(errorCode?: string, errorMessage?: string): string {
  // If we have a code, try to find exact match
  if (errorCode && ERROR_MESSAGES[errorCode.toLowerCase()]) {
    return ERROR_MESSAGES[errorCode.toLowerCase()]
  }

  // If we have a message, try pattern matching
  if (errorMessage) {
    for (const { pattern, message } of ERROR_PATTERNS) {
      const match = errorMessage.match(pattern)
      if (match) {
        // Replace $1, $2, etc with captured groups
        let translatedMessage = message
        for (let i = 1; i < match.length; i++) {
          translatedMessage = translatedMessage.replace(`$${i}`, match[i])
        }
        return translatedMessage
      }
    }
    
    // If message contains known keywords, try to extract and translate
    const lowerMessage = errorMessage.toLowerCase()
    
    if (lowerMessage.includes('team') && lowerMessage.includes('does not exist')) {
      return 'El equipo especificado no existe. Verifica el nombre del equipo.'
    }
    
    if (lowerMessage.includes('validation failed')) {
      return 'Error de validación. Por favor verifica los datos ingresados.'
    }
    
    if (lowerMessage.includes('invalid') && lowerMessage.includes('format')) {
      return 'Formato inválido. Por favor verifica los datos ingresados.'
    }
    
    if (lowerMessage.includes('cannot be empty')) {
      return 'Hay campos obligatorios sin completar.'
    }
  }

  // Default generic error
  return 'Ocurrió un error. Por favor intenta nuevamente.'
}

/**
 * Extracts error details from API error response
 * @param error - Error object from API
 * @returns Translated error message
 */
export function getErrorMessage(error: any): string {
  // If it's an Error object with message
  if (error instanceof Error) {
    return translateError(undefined, error.message)
  }

  // If it's an API error with code and message
  if (error && typeof error === 'object') {
    const code = error.code || error.errorCode
    const message = error.message || error.errorMessage || error.error
    return translateError(code, message)
  }

  // If it's just a string
  if (typeof error === 'string') {
    return translateError(undefined, error)
  }

  // Default
  return 'Ocurrió un error inesperado.'
}

