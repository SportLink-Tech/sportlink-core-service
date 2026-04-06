import { Team, UpdateTeamRequest } from '../../../../shared/types/team'
import { TeamRepository } from '../ports/TeamRepository'

export class UpdateTeamUseCase {
  constructor(private readonly teamRepository: TeamRepository) {}

  async execute(
    sport: string,
    currentName: string,
    request: UpdateTeamRequest,
  ): Promise<{ team: Team | null; success: boolean; error?: string }> {
    try {
      const response = await this.teamRepository.updateTeam(sport, currentName, request)
      return { team: response.data, success: true }
    } catch (error) {
      return {
        team: null,
        success: false,
        error: error instanceof Error ? error.message : 'Error updating team',
      }
    }
  }
}
