export type Sport = 'Football' | 'Paddle' | 'Tennis'

export interface Player {
  ID: string
  Category: number
  Sport: Sport
}

export interface TeamStats {
  Wins: number
  Losses: number
  Draws: number
}

export interface Team {
  Name: string
  Category: number
  Stats: TeamStats
  Sport: Sport
  Members: Player[]
}

export interface CreateTeamRequest {
  sport: Sport
  name: string
  category?: number
  players?: string[]
}

export interface ApiError {
  error: {
    code: string
    message: string
  }
}

