import { MatchRequestRepository } from '../ports/MatchRequestRepository'

export class CancelMatchRequestUseCase {
  constructor(private readonly repository: MatchRequestRepository) {}

  async execute(requesterAccountId: string, requestId: string): Promise<{ success: boolean; error?: string }> {
    try {
      await this.repository.cancel(requesterAccountId, requestId)
      return { success: true }
    } catch (e) {
      return { success: false, error: e instanceof Error ? e.message : 'Error al cancelar la solicitud' }
    }
  }
}
