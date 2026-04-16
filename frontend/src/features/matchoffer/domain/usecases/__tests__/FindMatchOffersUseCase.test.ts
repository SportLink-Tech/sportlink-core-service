import { describe, it, expect, vi, beforeEach } from 'vitest'
import { FindMatchOffersUseCase } from '../FindMatchOffersUseCase'
import { MatchOfferRepository } from '../../ports/MatchOfferRepository'
import { FindMatchOffersQuery, MatchOffer, PaginatedMatchOffersResponse } from '../../../../../shared/types/matchOffer'

const ACCOUNT_ID = 'account-123'

const mockRepository: MatchOfferRepository = {
  create: vi.fn(),
  search: vi.fn(),
  findByAccount: vi.fn(),
  retrieve: vi.fn(),
  delete: vi.fn(),
  confirm: vi.fn(),
}

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

const mockPaginatedResponse = (data: MatchOffer[]): PaginatedMatchOffersResponse => ({
  data,
  pagination: { number: 1, out_of: 1, total: data.length },
})

describe('FindMatchOffersUseCase', () => {
  let useCase: FindMatchOffersUseCase

  beforeEach(() => {
    useCase = new FindMatchOffersUseCase(mockRepository)
    vi.clearAllMocks()
  })

  describe('Success Cases', () => {
    it('should search announcements successfully', async () => {
      vi.mocked(mockRepository.search).mockResolvedValue({
        data: mockPaginatedResponse(mockAnnouncements),
        status: 200,
      })

      const query: FindMatchOffersQuery = { sports: ['Paddle'] }
      const result = await useCase.execute(ACCOUNT_ID, query)

      expect(result.success).toBe(true)
      expect(result.announcements).toEqual(mockAnnouncements)
      expect(result.error).toBeUndefined()
      expect(mockRepository.search).toHaveBeenCalledWith(
        ACCOUNT_ID,
        expect.objectContaining({ sports: ['Paddle'], fromDate: expect.any(String) })
      )
    })

    it('should add fromDate as today if not provided', async () => {
      vi.mocked(mockRepository.search).mockResolvedValue({
        data: mockPaginatedResponse(mockAnnouncements),
        status: 200,
      })

      const today = new Date().toISOString().split('T')[0]
      await useCase.execute(ACCOUNT_ID, {})

      expect(mockRepository.search).toHaveBeenCalledWith(
        ACCOUNT_ID,
        expect.objectContaining({ fromDate: today })
      )
    })

    it('should preserve provided fromDate', async () => {
      vi.mocked(mockRepository.search).mockResolvedValue({
        data: mockPaginatedResponse(mockAnnouncements),
        status: 200,
      })

      const customDate = '2025-12-10'
      await useCase.execute(ACCOUNT_ID, { fromDate: customDate })

      expect(mockRepository.search).toHaveBeenCalledWith(
        ACCOUNT_ID,
        expect.objectContaining({ fromDate: customDate })
      )
    })

    it('should return empty array for 404 status', async () => {
      vi.mocked(mockRepository.search).mockResolvedValue({
        data: mockPaginatedResponse([]),
        status: 404,
      })

      const result = await useCase.execute(ACCOUNT_ID, {})

      expect(result.success).toBe(true)
      expect(result.announcements).toEqual([])
      expect(result.error).toBeUndefined()
    })

    it('should return pagination metadata', async () => {
      vi.mocked(mockRepository.search).mockResolvedValue({
        data: { data: mockAnnouncements, pagination: { number: 2, out_of: 5, total: 45 } },
        status: 200,
      })

      const result = await useCase.execute(ACCOUNT_ID, {})

      expect(result.pagination).toEqual({ number: 2, outOf: 5, total: 45 })
    })
  })

  describe('Error Cases', () => {
    it('should handle repository error', async () => {
      vi.mocked(mockRepository.search).mockRejectedValue(new Error('Network error'))

      const result = await useCase.execute(ACCOUNT_ID, {})

      expect(result.success).toBe(false)
      expect(result.error).toBeDefined()
      expect(result.announcements).toEqual([])
    })

    it('should handle non-200/404 status codes', async () => {
      vi.mocked(mockRepository.search).mockResolvedValue({
        data: mockPaginatedResponse([]),
        status: 500,
      })

      const result = await useCase.execute(ACCOUNT_ID, {})

      expect(result.success).toBe(false)
      expect(result.error).toBe('Error al obtener las ofertas')
    })
  })
})
