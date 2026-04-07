import { MatchOffer } from '../../../../shared/types/matchOffer'
import { MatchOfferRepository } from '../ports/MatchOfferRepository'

export class RetrieveMatchOfferUseCase {
  constructor(private readonly repository: MatchOfferRepository) {}

  async execute(offerId: string): Promise<MatchOffer | null> {
    try {
      return await this.repository.retrieve(offerId)
    } catch {
      return null
    }
  }
}
