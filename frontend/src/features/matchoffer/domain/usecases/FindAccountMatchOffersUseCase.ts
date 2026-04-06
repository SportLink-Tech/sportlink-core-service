import { MatchOffer } from '../../../../shared/types/matchOffer'
import { MatchOfferRepository } from '../ports/MatchOfferRepository'

export class FindAccountMatchOffersUseCase {
  constructor(private readonly repository: MatchOfferRepository) {}

  async execute(accountId: string, statuses?: string[]): Promise<{ success: boolean; offerIds: Set<string>; offers: MatchOffer[]; error?: string }> {
    try {
      const result = await this.repository.findByAccount(accountId, statuses)
      const ids = new Set(result.data.map((a) => a.id).filter((id): id is string => Boolean(id)))
      return { success: true, offerIds: ids, offers: result.data }
    } catch (error) {
      return {
        success: false,
        offerIds: new Set(),
        offers: [],
        error: error instanceof Error ? error.message : 'Error al obtener las ofertas de la cuenta',
      }
    }
  }
}
