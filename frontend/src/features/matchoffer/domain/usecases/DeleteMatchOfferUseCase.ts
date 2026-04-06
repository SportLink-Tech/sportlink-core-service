import { MatchOfferRepository } from '../ports/MatchOfferRepository'

export class DeleteMatchOfferUseCase {
  constructor(private readonly repository: MatchOfferRepository) {}

  async execute(accountId: string, offerId: string): Promise<{ success: boolean; error?: string }> {
    try {
      await this.repository.delete(accountId, offerId)
      return { success: true }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Error al eliminar la oferta',
      }
    }
  }
}
