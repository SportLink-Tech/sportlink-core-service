import { Team } from '../../../../shared/types/team'
import { TeamRepository } from '../ports/TeamRepository'

/**
 * Use Case: Search Team
 * Contains domain logic and business rules
 * Depends on TeamRepository port (Dependency Inversion Principle)
 */
export class SearchTeamUseCase {
  constructor(private readonly teamRepository: TeamRepository) {}

  async execute(
    sport: string, 
    teamName?: string, 
    categories?: number[]
  ): Promise<{ teams: Team[]; success: boolean; error?: string }> {
    try {
      // Business rule: Validate inputs
      if (!sport || sport.trim().length === 0) {
        return {
          teams: [],
          success: false,
          error: 'Sport is required'
        }
      }

      // Sport is the only required parameter
      // Name and categories are optional - can search by sport alone
      const response = await this.teamRepository.findTeam(sport, teamName, categories)
      
      // Normalize Members array for each team (backend may return null)
      const teams = response.data.map(team => ({
        ...team,
        Members: team.Members || []
      }))

      return {
        teams,
        success: true
      }
    } catch (error) {
      return {
        teams: [],
        success: false,
        error: error instanceof Error ? error.message : 'Error searching teams'
      }
    }
  }
}

