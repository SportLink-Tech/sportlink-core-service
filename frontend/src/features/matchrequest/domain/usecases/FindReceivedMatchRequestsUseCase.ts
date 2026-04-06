import { MatchRequest, MatchRequestRepository } from '../ports/MatchRequestRepository'

export class FindReceivedMatchRequestsUseCase {
  constructor(private readonly repository: MatchRequestRepository) {}

  async execute(ownerAccountId: string): Promise<{ success: boolean; requests: MatchRequest[]; error?: string }> {
    try {
      const requests = await this.repository.findReceived(ownerAccountId)
      return { success: true, requests }
    } catch (error) {
      return {
        success: false,
        requests: [],
        error: error instanceof Error ? error.message : 'Error al obtener las solicitudes',
      }
    }
  }
}
