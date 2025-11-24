import { Team, CreateTeamRequest } from '../../../../shared/types/team'

/**
 * Port (Interface) for Team Repository
 * This is the contract that adapters must implement
 * Following Hexagonal Architecture - Domain defines the interface
 */
export interface TeamRepository {
  createTeam(request: CreateTeamRequest): Promise<{ data: Team; status: number }>
  findTeam(sport: string, teamName?: string, categories?: number[]): Promise<{ data: Team[]; status: number }>
}

