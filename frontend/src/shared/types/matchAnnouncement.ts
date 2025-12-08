export interface MatchAnnouncement {
  id?: string
  team_name: string
  sport: string
  day: string // ISO date string
  time_slot: {
    start_time: string // ISO datetime string
    end_time: string // ISO datetime string
  }
  location: {
    country: string
    province: string
    locality: string
  }
  admitted_categories: {
    type: 'SPECIFIC' | 'GREATER_THAN' | 'LESS_THAN' | 'BETWEEN'
    categories?: number[]
    min_level?: number
    max_level?: number
  }
  status: 'PENDING' | 'CONFIRMED' | 'CANCELLED' | 'EXPIRED'
  created_at: string
}

export interface CreateMatchAnnouncementRequest {
  team_name: string
  sport: string
  day: string
  time_slot: {
    start_time: string
    end_time: string
  }
  location: {
    country: string
    province: string
    locality: string
  }
  admitted_categories: {
    type: 'SPECIFIC' | 'GREATER_THAN' | 'LESS_THAN' | 'BETWEEN'
    categories?: number[]
    min_level?: number
    max_level?: number
  }
}

export interface FindMatchAnnouncementsQuery {
  sports?: string[]
  categories?: number[]
  statuses?: string[]
  fromDate?: string
  toDate?: string
  location?: {
    country?: string
    province?: string
    locality?: string
  }
  limit?: number
  offset?: number
}

export interface PaginationInfo {
  number: number // Current page number (1-based)
  out_of: number // Total number of pages
  total: number // Total number of items matching the query
}

export interface PaginatedMatchAnnouncementsResponse {
  data: MatchAnnouncement[]
  pagination: PaginationInfo
}

export interface ApiError {
  code: string
  message: string
}

