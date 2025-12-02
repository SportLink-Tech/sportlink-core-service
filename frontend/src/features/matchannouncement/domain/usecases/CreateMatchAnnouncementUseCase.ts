import { MatchAnnouncementRepository } from '../ports/MatchAnnouncementRepository'
import { MatchAnnouncement, CreateMatchAnnouncementRequest } from '../../../../shared/types/matchAnnouncement'
import { getErrorMessage } from '../../../../shared/utils/errorMessages'

export class CreateMatchAnnouncementUseCase {
  constructor(private repository: MatchAnnouncementRepository) {}

  async execute(request: CreateMatchAnnouncementRequest): Promise<{ announcement: MatchAnnouncement; success: boolean; error?: string }> {
    try {
      // Validations
      if (!request.team_name || request.team_name.trim().length === 0) {
        return { announcement: {} as MatchAnnouncement, success: false, error: 'El nombre del equipo es obligatorio' }
      }

      if (!request.sport || request.sport.trim().length === 0) {
        return { announcement: {} as MatchAnnouncement, success: false, error: 'El deporte es obligatorio' }
      }

      if (!request.day) {
        return { announcement: {} as MatchAnnouncement, success: false, error: 'La fecha del partido es obligatoria' }
      }

      if (!request.location.country || !request.location.province || !request.location.locality) {
        return { announcement: {} as MatchAnnouncement, success: false, error: 'La ubicación completa es obligatoria' }
      }

      const response = await this.repository.create(request)

      if (response.status === 201 || response.status === 200) {
        return { announcement: response.data, success: true }
      }

      return { announcement: {} as MatchAnnouncement, success: false, error: 'Error al crear el anuncio' }
    } catch (error) {
      return {
        announcement: {} as MatchAnnouncement,
        success: false,
        error: getErrorMessage(error),
      }
    }
  }
}

