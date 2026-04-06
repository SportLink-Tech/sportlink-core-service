import { MatchOfferRepository } from '../ports/MatchOfferRepository'
import { MatchOffer, CreateMatchOfferRequest } from '../../../../shared/types/matchOffer'
import { getErrorMessage } from '../../../../shared/utils/errorMessages'

export class CreateMatchOfferUseCase {
  constructor(private repository: MatchOfferRepository) {}

  async execute(accountId: string, request: CreateMatchOfferRequest): Promise<{ announcement: MatchOffer; success: boolean; error?: string }> {
    try {
      // Validations
      if (!request.team_name || request.team_name.trim().length === 0) {
        return { announcement: {} as MatchOffer, success: false, error: 'El nombre del equipo es obligatorio' }
      }

      if (!request.sport || request.sport.trim().length === 0) {
        return { announcement: {} as MatchOffer, success: false, error: 'El deporte es obligatorio' }
      }

      if (!request.day) {
        return { announcement: {} as MatchOffer, success: false, error: 'La fecha del partido es obligatoria' }
      }

      if (!request.location.country || !request.location.province || !request.location.locality) {
        return { announcement: {} as MatchOffer, success: false, error: 'La ubicación completa es obligatoria' }
      }

      const response = await this.repository.create(accountId, request)

      if (response.status === 201 || response.status === 200) {
        return { announcement: response.data, success: true }
      }

      return { announcement: {} as MatchOffer, success: false, error: 'Error al crear la oferta' }
    } catch (error) {
      return {
        announcement: {} as MatchOffer,
        success: false,
        error: getErrorMessage(error),
      }
    }
  }
}

