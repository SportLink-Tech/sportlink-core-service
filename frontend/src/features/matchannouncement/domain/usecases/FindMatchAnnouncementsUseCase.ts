import { MatchAnnouncementRepository } from '../ports/MatchAnnouncementRepository'
import { MatchAnnouncement, FindMatchAnnouncementsQuery } from '../../../../shared/types/matchAnnouncement'
import { getErrorMessage } from '../../../../shared/utils/errorMessages'

export class FindMatchAnnouncementsUseCase {
  constructor(private repository: MatchAnnouncementRepository) {}

  async execute(query: FindMatchAnnouncementsQuery): Promise<{ announcements: MatchAnnouncement[]; success: boolean; error?: string }> {
    try {
      // Always add fromDate as today if not provided
      const queryWithDefaults = {
        ...query,
        fromDate: query.fromDate || new Date().toISOString().split('T')[0],
      }

      const response = await this.repository.find(queryWithDefaults)

      if (response.status === 200) {
        return { announcements: response.data, success: true }
      }

      if (response.status === 404) {
        return { announcements: [], success: true }
      }

      return { announcements: [], success: false, error: 'Error al obtener los anuncios' }
    } catch (error) {
      return {
        announcements: [],
        success: false,
        error: getErrorMessage(error),
      }
    }
  }
}

