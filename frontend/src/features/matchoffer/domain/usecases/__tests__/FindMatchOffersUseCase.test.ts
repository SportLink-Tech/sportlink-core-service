import { describe, it, expect, vi, beforeEach } from 'vitest'
import { FindMatchOffersUseCase } from '../FindMatchOffersUseCase'
import { MatchOfferRepository } from '../../ports/MatchOfferRepository'
import { FindMatchOffersQuery, MatchOffer } from '../../../../../shared/types/matchOffer'

// Mock del repositorio
const mockRepository: MatchOfferRepository = {
  create: vi.fn(),
  find: vi.fn(),
}

describe('FindMatchOffersUseCase', () => {
  let useCase: FindMatchOffersUseCase

  beforeEach(() => {
    useCase = new FindMatchOffersUseCase(mockRepository)
    vi.clearAllMocks()
  })

  const mockAnnouncements: MatchOffer[] = [
    {
      id: '1',
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
    },
  ]

  describe('Success Cases', () => {
    it('should find announcements successfully', async () => {
      vi.mocked(mockRepository.find).mockResolvedValue({
        data: mockAnnouncements,
        status: 200,
      })

      const query: FindMatchOffersQuery = {
        sports: ['Paddle'],
      }

      const result = await useCase.execute(query)

      expect(result.success).toBe(true)
      expect(result.announcements).toEqual(mockAnnouncements)
      expect(result.error).toBeUndefined()
      expect(mockRepository.find).toHaveBeenCalledWith(
        expect.objectContaining({
          sports: ['Paddle'],
          fromDate: expect.any(String),
        })
      )
    })

    it('should add fromDate as today if not provided', async () => {
      vi.mocked(mockRepository.find).mockResolvedValue({
        data: mockAnnouncements,
        status: 200,
      })

      const today = new Date().toISOString().split('T')[0]
      const query: FindMatchOffersQuery = {}

      await useCase.execute(query)

      expect(mockRepository.find).toHaveBeenCalledWith(
        expect.objectContaining({
          fromDate: today,
        })
      )
    })

    it('should preserve provided fromDate', async () => {
      vi.mocked(mockRepository.find).mockResolvedValue({
        data: mockAnnouncements,
        status: 200,
      })

      const customDate = '2025-12-10'
      const query: FindMatchOffersQuery = {
        fromDate: customDate,
      }

      await useCase.execute(query)

      expect(mockRepository.find).toHaveBeenCalledWith(
        expect.objectContaining({
          fromDate: customDate,
        })
      )
    })

    it('should return empty array for 404 status', async () => {
      vi.mocked(mockRepository.find).mockResolvedValue({
        data: [],
        status: 404,
      })

      const result = await useCase.execute({})

      expect(result.success).toBe(true)
      expect(result.announcements).toEqual([])
      expect(result.error).toBeUndefined()
    })
  })

  describe('Error Cases', () => {
    it('should handle repository error', async () => {
      const errorMessage = 'Network error'
      vi.mocked(mockRepository.find).mockRejectedValue(new Error(errorMessage))

      const result = await useCase.execute({})

      expect(result.success).toBe(false)
      expect(result.error).toBeDefined()
      expect(result.announcements).toEqual([])
    })

    it('should handle API error with code', async () => {
      vi.mocked(mockRepository.find).mockRejectedValue({
        code: 'use_case_execution_failed',
        message: 'Database connection failed',
      })

      const result = await useCase.execute({})

      expect(result.success).toBe(false)
      expect(result.error).toBeDefined()
      expect(result.announcements).toEqual([])
    })

    it('should handle non-200/404 status codes', async () => {
      vi.mocked(mockRepository.find).mockResolvedValue({
        data: [],
        status: 500,
      })

      const result = await useCase.execute({})

      expect(result.success).toBe(false)
      expect(result.error).toBe('Error al obtener las ofertas')
    })
  })
})

