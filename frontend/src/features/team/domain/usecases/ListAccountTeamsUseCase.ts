import { Team } from '../../../../shared/types/team'
import { TeamRepository } from '../ports/TeamRepository'

export class ListAccountTeamsUseCase {
  constructor(private readonly teamRepository: TeamRepository) {}

  async execute(accountId: string): Promise<{ teams: Team[]; success: boolean; error?: string }> {
    try {
      const response = await this.teamRepository.listAccountTeams(accountId)
      return {
        teams: response.data || [],
        success: true,
      }
    } catch (error) {
      return {
        teams: [],
        success: false,
        error: error instanceof Error ? error.message : 'Error listing account teams',
      }
    }
  }
}
