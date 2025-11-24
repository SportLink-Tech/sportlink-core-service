import { Team, CreateTeamRequest } from '../../../../shared/types/team'
import { TeamRepository } from '../ports/TeamRepository'

/**
 * Use Case: Create Team
 * Contains domain logic and business rules
 * Depends on TeamRepository port (Dependency Inversion Principle)
 */
export class CreateTeamUseCase {
  constructor(private readonly teamRepository: TeamRepository) {}

  async execute(request: CreateTeamRequest): Promise<{ team: Team | null; success: boolean; error?: string }> {
    try {
      // Business rule: Validate team name
      if (!request.name || request.name.trim().length === 0) {
        return {
          team: null,
          success: false,
          error: 'Team name is required'
        }
      }

      // Business rule: Validate sport
      if (!request.sport) {
        return {
          team: null,
          success: false,
          error: 'Sport is required'
        }
      }

      const response = await this.teamRepository.createTeam(request)
      
      // Check HTTP status
      if (response.status === 201) {
        return {
          team: response.data,
          success: true
        }
      }

      return {
        team: null,
        success: false,
        error: 'Failed to create team'
      }
    } catch (error) {
      return {
        team: null,
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error occurred'
      }
    }
  }
}

