import { Team, CreateTeamRequest, UpdateTeamRequest } from '../../../../shared/types/team'

/**
 * Port (Interface) for Team Repository
 * This is the contract that adapters must implement
 * Following Hexagonal Architecture - Domain defines the interface
 */
export interface TeamRepository {
  createTeam(accountId: string, request: CreateTeamRequest): Promise<{ data: Team; status: number }>
  findTeam(sport: string, teamName?: string, categories?: number[]): Promise<{ data: Team[]; status: number }>
  listAccountTeams(accountId: string): Promise<{ data: Team[]; status: number }>
  updateTeam(sport: string, currentName: string, request: UpdateTeamRequest): Promise<{ data: Team; status: number }>
}

