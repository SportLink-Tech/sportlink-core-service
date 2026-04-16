import { MatchOfferRepository } from '../ports/MatchOfferRepository'
import { FindMatchOffersQuery, PaginatedMatchOffersResponse } from '../../../../shared/types/matchOffer'
import { getErrorMessage } from '../../../../shared/utils/errorMessages'

export interface FindMatchOffersResult {
  announcements: PaginatedMatchOffersResponse['data']
  pagination: {
    number: number
    outOf: number
    total: number
  }
  success: boolean
  error?: string
}

export class FindMatchOffersUseCase {
  constructor(private repository: MatchOfferRepository) {}

  async execute(accountId: string, query: FindMatchOffersQuery): Promise<FindMatchOffersResult> {
    try {
      const queryWithDefaults = {
        ...query,
        fromDate: query.fromDate || new Date().toISOString().split('T')[0],
      }

      const response = await this.repository.search(accountId, queryWithDefaults)

      if (response.status === 200) {
        const paginatedResponse: PaginatedMatchOffersResponse = response.data
        return {
          announcements: paginatedResponse.data,
          pagination: {
            number: paginatedResponse.pagination.number,
            outOf: paginatedResponse.pagination.out_of,
            total: paginatedResponse.pagination.total,
          },
          success: true,
        }
      }

      if (response.status === 404) {
        return {
          announcements: [],
          pagination: { number: 1, outOf: 0, total: 0 },
          success: true,
        }
      }

      return {
        announcements: [],
        pagination: { number: 1, outOf: 0, total: 0 },
        success: false,
        error: 'Error al obtener las ofertas',
      }
    } catch (error) {
      return {
        announcements: [],
        pagination: { number: 1, outOf: 0, total: 0 },
        success: false,
        error: getErrorMessage(error),
      }
    }
  }
}
