import {
  MatchOffer,
  CreateMatchOfferRequest,
  FindMatchOffersQuery,
  ApiError,
  PaginatedMatchOffersResponse,
} from '../../../../shared/types/matchOffer'
import { MatchOfferRepository } from '../../domain/ports/MatchOfferRepository'

const API_BASE_URL = '/api'

/**
 * Adapter: MatchOffer API Adapter
 * Implements MatchOfferRepository port
 * This is in the infrastructure layer, so it's tightly coupled to external APIs
 * Following Hexagonal Architecture - Adapter implements Port
 */
export class MatchOfferApiAdapter implements MatchOfferRepository {
  async create(accountId: string, request: CreateMatchOfferRequest): Promise<{ data: MatchOffer; status: number }> {
    const response = await fetch(`${API_BASE_URL}/account/${encodeURIComponent(accountId)}/match-offer`, {
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

  async search(accountId: string, query: FindMatchOffersQuery): Promise<{ data: PaginatedMatchOffersResponse; status: number }> {
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
      queryParams.append('from_date', query.fromDate)
    }

    if (query.toDate) {
      queryParams.append('to_date', query.toDate)
    }

    if (query.geoFilter) {
      queryParams.append('lat', query.geoFilter.latitude.toString())
      queryParams.append('lng', query.geoFilter.longitude.toString())
      queryParams.append('radius_km', query.geoFilter.radiusKm.toString())
    } else if (query.location) {
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
    const url = `${API_BASE_URL}/account/${encodeURIComponent(accountId)}/match-offer/search${queryString ? `?${queryString}` : ''}`

    const response = await fetch(url, { method: 'GET', headers: { 'Content-Type': 'application/json' } })
    const data = await response.json()

    if (!response.ok) {
      const error: ApiError = data
      throw new Error(error.message || 'Error finding match announcements')
    }

    return { data, status: response.status }
  }

  async findByAccount(accountId: string, statuses?: string[]): Promise<{ data: MatchOffer[]; status: number }> {
    const queryParams = new URLSearchParams()
    if (statuses && statuses.length > 0) {
      queryParams.append('statuses', statuses.join(','))
    }
    const queryString = queryParams.toString()
    const url = `${API_BASE_URL}/account/${encodeURIComponent(accountId)}/match-offer${queryString ? `?${queryString}` : ''}`

    const response = await fetch(url, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
    })

    const data = await response.json()

    if (!response.ok) {
      const error: ApiError = data
      throw new Error(error.message || 'Error finding account match announcements')
    }

    return { data, status: response.status }
  }

  async retrieve(offerId: string): Promise<MatchOffer> {
    const response = await fetch(`${API_BASE_URL}/match-offer/${offerId}`)
    const data = await response.json()
    if (!response.ok) {
      throw new Error(data.message || 'Error al obtener la publicación')
    }
    return data
  }

  async delete(accountId: string, offerId: string): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/account/${encodeURIComponent(accountId)}/match-offer/${offerId}`, {
      method: 'DELETE',
    })

    if (!response.ok) {
      const data = await response.json().catch(() => ({}))
      const error: ApiError = data
      throw new Error(error.message || 'Error al eliminar la oferta')
    }
  }

  async confirm(accountId: string, offerId: string): Promise<void> {
    const response = await fetch(
      `${API_BASE_URL}/account/${encodeURIComponent(accountId)}/match-offer/${offerId}/confirm`,
      { method: 'POST' }
    )

    if (!response.ok) {
      const data = await response.json().catch(() => ({}))
      const error: ApiError = data
      throw new Error(error.message || 'Error al confirmar el partido')
    }
  }
}

