import { MatchAnnouncementRepository } from '../ports/MatchAnnouncementRepository'

export class FindAccountMatchAnnouncementsUseCase {
  constructor(private readonly repository: MatchAnnouncementRepository) {}

  async execute(accountId: string, statuses?: string[]): Promise<{ success: boolean; announcementIds: Set<string>; error?: string }> {
    try {
      const result = await this.repository.findByAccount(accountId, statuses)
      const ids = new Set(result.data.map((a) => a.id).filter((id): id is string => Boolean(id)))
      return { success: true, announcementIds: ids }
    } catch (error) {
      return {
        success: false,
        announcementIds: new Set(),
        error: error instanceof Error ? error.message : 'Error al obtener los anuncios de la cuenta',
      }
    }
  }
}
