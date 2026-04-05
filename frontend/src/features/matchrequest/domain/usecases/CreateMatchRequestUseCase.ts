import { MatchRequestRepository } from '../ports/MatchRequestRepository'

export class CreateMatchRequestUseCase {
  constructor(private readonly repository: MatchRequestRepository) {}

  async execute(
    requesterAccountId: string,
    matchAnnouncementId: string,
  ): Promise<{ success: boolean; error?: string }> {
    try {
      await this.repository.create(requesterAccountId, matchAnnouncementId)
      return { success: true }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Error al unirse al partido',
      }
    }
  }
}
