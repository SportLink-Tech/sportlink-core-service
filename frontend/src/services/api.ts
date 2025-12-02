import { CreateTeamRequest, Team } from '../types/team'

interface ApiError {
  code: string
  message: string
}

const API_BASE_URL = '/api'

interface ApiResponse<T> {
  data: T
  status: number
}

class ApiService {
  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<ApiResponse<T>> {
    const url = `${API_BASE_URL}${endpoint}`
    
    try {
      const response = await fetch(url, {
        ...options,
        headers: {
          'Content-Type': 'application/json',
          ...options?.headers,
        },
      })

      const data = await response.json()

      if (!response.ok) {
        const error = data as ApiError
        throw new Error(error.message || 'Error en la petición')
      }

      return {
        data,
        status: response.status
      }
    } catch (error) {
      if (error instanceof Error) {
        throw error
      }
      throw new Error('Error desconocido')
    }
  }

  async createTeam(data: CreateTeamRequest): Promise<ApiResponse<Team>> {
    return this.request<Team>('/team', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async getTeam(sport: string, teamName: string): Promise<ApiResponse<Team>> {
    return this.request<Team>(`/sport/${sport}/team/${teamName}`, {
      method: 'GET',
    })
  }

  async findTeams(sport: string, teamName?: string, categories?: number[]): Promise<ApiResponse<Team[]>> {
    const queryParams = new URLSearchParams()
    
    if (teamName) {
      queryParams.append('name', teamName)
    }
    
    if (categories && categories.length > 0) {
      queryParams.append('category', categories.join(','))
    }
    
    const queryString = queryParams.toString()
    const url = `/sport/${sport}/team${queryString ? `?${queryString}` : ''}`
    
    return this.request<Team[]>(url, {
      method: 'GET',
    })
  }
}

export type { ApiResponse }

export const apiService = new ApiService()

