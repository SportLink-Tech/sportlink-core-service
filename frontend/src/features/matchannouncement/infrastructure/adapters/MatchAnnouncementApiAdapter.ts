import {
  MatchAnnouncement,
  CreateMatchAnnouncementRequest,
  FindMatchAnnouncementsQuery,
  ApiError,
  PaginatedMatchAnnouncementsResponse,
} from '../../../../shared/types/matchAnnouncement'
import { MatchAnnouncementRepository } from '../../domain/ports/MatchAnnouncementRepository'

const API_BASE_URL = '/api'

/**
 * Adapter: MatchAnnouncement API Adapter
 * Implements MatchAnnouncementRepository port
 * This is in the infrastructure layer, so it's tightly coupled to external APIs
 * Following Hexagonal Architecture - Adapter implements Port
 */
export class MatchAnnouncementApiAdapter implements MatchAnnouncementRepository {
  async create(request: CreateMatchAnnouncementRequest): Promise<{ data: MatchAnnouncement; status: number }> {
    const response = await fetch(`${API_BASE_URL}/match-announcement`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    })

    const data = await response.json()

    if (!response.ok) {
      const error: ApiError = data
      throw new Error(error.message || 'Error creating match announcement')
    }

    return {
      data,
      status: response.status,
    }
  }

  async find(query: FindMatchAnnouncementsQuery): Promise<{ data: PaginatedMatchAnnouncementsResponse; status: number }> {
    const queryParams = new URLSearchParams()

    if (query.sports && query.sports.length > 0) {
      queryParams.append('sports', query.sports.join(','))
    }

    if (query.categories && query.categories.length > 0) {
      queryParams.append('categories', query.categories.join(','))
    }

    if (query.statuses && query.statuses.length > 0) {
      queryParams.append('statuses', query.statuses.join(','))
    }

    if (query.fromDate) {
      queryParams.append('fromDate', query.fromDate)
    }

    if (query.toDate) {
      queryParams.append('toDate', query.toDate)
    }

    if (query.location) {
      if (query.location.country) {
        queryParams.append('country', query.location.country)
      }
      if (query.location.province) {
        queryParams.append('province', query.location.province)
      }
      if (query.location.locality) {
        queryParams.append('locality', query.location.locality)
      }
    }

    if (query.limit !== undefined && query.limit > 0) {
      queryParams.append('limit', query.limit.toString())
    }

    if (query.offset !== undefined && query.offset > 0) {
      queryParams.append('offset', query.offset.toString())
    }

    const queryString = queryParams.toString()
    const url = `${API_BASE_URL}/match-announcement${queryString ? `?${queryString}` : ''}`

    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    })

    const data = await response.json()

    if (!response.ok) {
      const error: ApiError = data
      throw new Error(error.message || 'Error finding match announcements')
    }

    return {
      data,
      status: response.status,
    }
  }
}

