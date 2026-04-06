import { MatchRequestRepository } from '../ports/MatchRequestRepository'

export class FindSentMatchRequestsUseCase {
  constructor(private readonly repository: MatchRequestRepository) {}

  async execute(requesterAccountId: string, statuses?: string[]): Promise<Set<string>> {
    try {
      const requests = await this.repository.findSent(requesterAccountId, statuses)
      return new Set(requests.map((r) => r.match_offer_id).filter(Boolean))
    } catch {
      return new Set()
    }
  }
}
