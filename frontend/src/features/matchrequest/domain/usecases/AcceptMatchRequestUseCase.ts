import { MatchRequestRepository } from '../ports/MatchRequestRepository'

export class AcceptMatchRequestUseCase {
  constructor(private readonly repository: MatchRequestRepository) {}

  async execute(ownerAccountId: string, requestId: string): Promise<{ success: boolean; error?: string }> {
    try {
      await this.repository.accept(ownerAccountId, requestId)
      return { success: true }
    } catch (e) {
      return { success: false, error: e instanceof Error ? e.message : 'Error al aceptar la solicitud' }
    }
  }
}
