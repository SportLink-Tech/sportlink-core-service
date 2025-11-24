import { Team, CreateTeamRequest, ApiError } from '../../../../shared/types/team'
import { TeamRepository } from '../../domain/ports/TeamRepository'

const API_BASE_URL = '/api'

/**
 * Adapter: Team API Adapter
 * Implements TeamRepository port
 * This is in the infrastructure layer, so it's tightly coupled to external APIs
 * Following Hexagonal Architecture - Adapter implements Port
 */
export class TeamApiAdapter implements TeamRepository {
  async createTeam(request: CreateTeamRequest): Promise<{ data: Team; status: number }> {
    const response = await fetch(`${API_BASE_URL}/team`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    })

    const data = await response.json()

    if (!response.ok) {
      const error: ApiError = data
      throw new Error(error.message || 'Error creating team')
    }

    return {
      data,
      status: response.status
    }
  }

  async findTeam(sport: string, teamName?: string, categories?: number[]): Promise<{ data: Team[]; status: number }> {
    const queryParams = new URLSearchParams()
    
    if (teamName) {
      queryParams.append('name', teamName)
    }
    
    if (categories && categories.length > 0) {
      queryParams.append('category', categories.join(','))
    }
    
    const queryString = queryParams.toString()
    const url = `${API_BASE_URL}/sport/${sport}/team${queryString ? `?${queryString}` : ''}`
    
    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    })

    const data = await response.json()

    if (!response.ok) {
      const error: ApiError = data
      throw new Error(error.message || 'Error finding teams')
    }

    return {
      data,
      status: response.status
    }
  }
}

