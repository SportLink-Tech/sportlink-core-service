import { describe, it, expect } from 'vitest'
import { translateError, getErrorMessage } from '../errorMessages'

describe('Error Messages', () => {
  describe('translateError', () => {
    it('should translate known error codes', () => {
      expect(translateError('request_validation_failed')).toBe('Error de validación en la solicitud')
      expect(translateError('not_found')).toBe('No se encontró el recurso solicitado')
      expect(translateError('invalid_request_format')).toBe('Formato de solicitud inválido')
    })

    it('should handle pattern matching for team not found', () => {
      const message = "team 'Boca Junior' for sport 'Paddle' does not exist"
      const result = translateError(undefined, message)
      expect(result).toBe("El equipo \"Boca Junior\" para el deporte \"Paddle\" no existe")
    })

    it('should handle pattern matching for cannot be empty', () => {
      const message = 'team name cannot be empty'
      const result = translateError(undefined, message)
      expect(result).toBe('team name no puede estar vacío')
    })

    it('should handle pattern matching for is required', () => {
      const message = 'Sport is required'
      const result = translateError(undefined, message)
      expect(result).toBe('Sport es obligatorio')
    })

    it('should handle generic validation failed message', () => {
      const message = 'validation failed: invalid data'
      const result = translateError(undefined, message)
      expect(result).toBe('Error de validación. Por favor verifica los datos ingresados.')
    })

    it('should handle generic invalid format message', () => {
      const message = 'invalid format detected'
      const result = translateError(undefined, message)
      expect(result).toBe('Formato inválido. Por favor verifica los datos ingresados.')
    })

    it('should return default message for unknown errors', () => {
      const result = translateError(undefined, 'some random error')
      expect(result).toBe('Ocurrió un error. Por favor intenta nuevamente.')
    })
  })

  describe('getErrorMessage', () => {
    it('should handle Error objects', () => {
      const error = new Error('Network connection failed')
      const result = getErrorMessage(error)
      expect(result).toBeDefined()
      expect(typeof result).toBe('string')
    })

    it('should handle API error objects with code', () => {
      const error = {
        code: 'request_validation_failed',
        message: 'Invalid data',
      }
      const result = getErrorMessage(error)
      expect(result).toBe('Error de validación en la solicitud')
    })

    it('should handle API error objects with errorCode', () => {
      const error = {
        errorCode: 'not_found',
        errorMessage: 'Resource not found',
      }
      const result = getErrorMessage(error)
      expect(result).toBe('No se encontró el recurso solicitado')
    })

    it('should handle string errors', () => {
      const error = 'validation failed: invalid input'
      const result = getErrorMessage(error)
      expect(result).toBe('Error de validación. Por favor verifica los datos ingresados.')
    })

    it('should handle unknown error types', () => {
      const error = 123
      const result = getErrorMessage(error)
      expect(result).toBe('Ocurrió un error inesperado.')
    })

    it('should handle null/undefined', () => {
      expect(getErrorMessage(null)).toBe('Ocurrió un error inesperado.')
      expect(getErrorMessage(undefined)).toBe('Ocurrió un error inesperado.')
    })
  })
})

