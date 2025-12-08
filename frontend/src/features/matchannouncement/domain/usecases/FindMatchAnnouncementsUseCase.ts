import { MatchAnnouncementRepository } from '../ports/MatchAnnouncementRepository'
import { MatchAnnouncement, FindMatchAnnouncementsQuery, PaginatedMatchAnnouncementsResponse } from '../../../../shared/types/matchAnnouncement'
import { getErrorMessage } from '../../../../shared/utils/errorMessages'

export interface FindMatchAnnouncementsResult {
  announcements: MatchAnnouncement[]
  pagination: {
    number: number // Current page number (1-based)
    outOf: number // Total number of pages
    total: number // Total number of items
  }
  success: boolean
  error?: string
}

export class FindMatchAnnouncementsUseCase {
  constructor(private repository: MatchAnnouncementRepository) {}

  async execute(query: FindMatchAnnouncementsQuery): Promise<FindMatchAnnouncementsResult> {
    try {
      // Always add fromDate as today if not provided
      const queryWithDefaults = {
        ...query,
        fromDate: query.fromDate || new Date().toISOString().split('T')[0],
      }

      const response = await this.repository.find(queryWithDefaults)

      if (response.status === 200) {
        const paginatedResponse: PaginatedMatchAnnouncementsResponse = response.data
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
          pagination: {
            number: 1,
            outOf: 0,
            total: 0,
          },
          success: true,
        }
      }

      return {
        announcements: [],
        pagination: {
          number: 1,
          outOf: 0,
          total: 0,
        },
        success: false,
        error: 'Error al obtener los anuncios',
      }
    } catch (error) {
      return {
        announcements: [],
        pagination: {
          number: 1,
          outOf: 0,
          total: 0,
        },
        success: false,
        error: getErrorMessage(error),
      }
    }
  }
}

