import { describe, it, expect, vi, beforeEach } from 'vitest'
import { CreateMatchAnnouncementUseCase } from '../CreateMatchAnnouncementUseCase'
import { MatchAnnouncementRepository } from '../../ports/MatchAnnouncementRepository'
import { CreateMatchAnnouncementRequest, MatchAnnouncement } from '../../../../../shared/types/matchAnnouncement'

// Mock del repositorio
const mockRepository: MatchAnnouncementRepository = {
  create: vi.fn(),
  find: vi.fn(),
}

describe('CreateMatchAnnouncementUseCase', () => {
  let useCase: CreateMatchAnnouncementUseCase

  beforeEach(() => {
    useCase = new CreateMatchAnnouncementUseCase(mockRepository)
    vi.clearAllMocks()
  })

  const validRequest: CreateMatchAnnouncementRequest = {
    team_name: 'Boca Junior',
    sport: 'Paddle',
    day: '2025-12-05',
    time_slot: {
      start_time: '2025-12-05T13:00:00',
      end_time: '2025-12-05T14:00:00',
    },
    location: {
      country: 'Argentina',
      province: 'Buenos Aires',
      locality: 'CABA',
    },
    admitted_categories: {
      type: 'GREATER_THAN',
      min_level: 5,
    },
  }

  const mockAnnouncement: MatchAnnouncement = {
    id: '123',
    team_name: 'Boca Junior',
    sport: 'Paddle',
    day: '2025-12-05',
    time_slot: {
      start_time: '2025-12-05T13:00:00',
      end_time: '2025-12-05T14:00:00',
    },
    location: {
      country: 'Argentina',
      province: 'Buenos Aires',
      locality: 'CABA',
    },
    admitted_categories: {
      type: 'GREATER_THAN',
      min_level: 5,
    },
    status: 'PENDING',
    created_at: '2025-12-01T10:00:00',
  }

  describe('Validations', () => {
    it('should fail when team name is empty', async () => {
      const request = { ...validRequest, team_name: '' }
      const result = await useCase.execute(request)

      expect(result.success).toBe(false)
      expect(result.error).toBe('El nombre del equipo es obligatorio')
      expect(mockRepository.create).not.toHaveBeenCalled()
    })

    it('should fail when sport is empty', async () => {
      const request = { ...validRequest, sport: '' }
      const result = await useCase.execute(request)

      expect(result.success).toBe(false)
      expect(result.error).toBe('El deporte es obligatorio')
      expect(mockRepository.create).not.toHaveBeenCalled()
    })

    it('should fail when day is empty', async () => {
      const request = { ...validRequest, day: '' }
      const result = await useCase.execute(request)

      expect(result.success).toBe(false)
      expect(result.error).toBe('La fecha del partido es obligatoria')
      expect(mockRepository.create).not.toHaveBeenCalled()
    })

    it('should fail when location is incomplete', async () => {
      const request = { ...validRequest, location: { country: '', province: '', locality: '' } }
      const result = await useCase.execute(request)

      expect(result.success).toBe(false)
      expect(result.error).toBe('La ubicación completa es obligatoria')
      expect(mockRepository.create).not.toHaveBeenCalled()
    })
  })

  describe('Success Cases', () => {
    it('should create announcement successfully with status 201', async () => {
      vi.mocked(mockRepository.create).mockResolvedValue({
        data: mockAnnouncement,
        status: 201,
      })

      const result = await useCase.execute(validRequest)

      expect(result.success).toBe(true)
      expect(result.announcement).toEqual(mockAnnouncement)
      expect(result.error).toBeUndefined()
      expect(mockRepository.create).toHaveBeenCalledWith(validRequest)
      expect(mockRepository.create).toHaveBeenCalledTimes(1)
    })

    it('should create announcement successfully with status 200', async () => {
      vi.mocked(mockRepository.create).mockResolvedValue({
        data: mockAnnouncement,
        status: 200,
      })

      const result = await useCase.execute(validRequest)

      expect(result.success).toBe(true)
      expect(result.announcement).toEqual(mockAnnouncement)
      expect(result.error).toBeUndefined()
    })
  })

  describe('Error Cases', () => {
    it('should handle repository error with Error object', async () => {
      const errorMessage = 'Network error'
      vi.mocked(mockRepository.create).mockRejectedValue(new Error(errorMessage))

      const result = await useCase.execute(validRequest)

      expect(result.success).toBe(false)
      expect(result.error).toBeDefined()
      expect(result.announcement).toEqual({})
    })

    it('should handle repository error with API error', async () => {
      vi.mocked(mockRepository.create).mockRejectedValue({
        code: 'request_validation_failed',
        message: 'team does not exist',
      })

      const result = await useCase.execute(validRequest)

      expect(result.success).toBe(false)
      expect(result.error).toBeDefined()
      expect(result.announcement).toEqual({})
    })

    it('should handle non-201/200 status codes', async () => {
      vi.mocked(mockRepository.create).mockResolvedValue({
        data: mockAnnouncement,
        status: 400,
      })

      const result = await useCase.execute(validRequest)

      expect(result.success).toBe(false)
      expect(result.error).toBe('Error al crear el anuncio')
    })
  })
})

