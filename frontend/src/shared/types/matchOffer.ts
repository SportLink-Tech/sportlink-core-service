export interface MatchOffer {
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
    latitude?: number
    longitude?: number
  }
  admitted_categories: {
    type: 'SPECIFIC' | 'GREATER_THAN' | 'LESS_THAN' | 'BETWEEN'
    categories?: number[]
    min_level?: number
    max_level?: number
  }
  status: 'PENDING' | 'CONFIRMED' | 'CANCELLED' | 'EXPIRED'
  created_at: string
  owner_account_id?: string
}

export interface CreateMatchOfferRequest {
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
    latitude?: number
    longitude?: number
  }
  admitted_categories: {
    type: 'SPECIFIC' | 'GREATER_THAN' | 'LESS_THAN' | 'BETWEEN'
    categories?: number[]
    min_level?: number
    max_level?: number
  }
}

export interface GeoFilter {
  latitude: number
  longitude: number
  radiusKm: number
}

export interface FindMatchOffersQuery {
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
  geoFilter?: GeoFilter
  limit?: number
  offset?: number
}

export interface PaginationInfo {
  number: number // Current page number (1-based)
  out_of: number // Total number of pages
  total: number // Total number of items matching the query
}

export interface PaginatedMatchOffersResponse {
  data: MatchOffer[]
  pagination: PaginationInfo
}

export interface ApiError {
  code: string
  message: string
}

