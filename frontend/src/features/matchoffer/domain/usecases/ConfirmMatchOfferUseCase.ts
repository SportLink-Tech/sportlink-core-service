import { MatchOfferRepository } from '../ports/MatchOfferRepository'

export class ConfirmMatchOfferUseCase {
  constructor(private readonly repository: MatchOfferRepository) {}

  async execute(accountId: string, offerId: string): Promise<{ success: boolean; error?: string }> {
    try {
      await this.repository.confirm(accountId, offerId)
      return { success: true }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Error al confirmar el partido',
      }
    }
  }
}
