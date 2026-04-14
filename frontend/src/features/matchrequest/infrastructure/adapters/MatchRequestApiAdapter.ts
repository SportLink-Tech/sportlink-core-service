import { MatchRequest, MatchRequestRepository } from '../../domain/ports/MatchRequestRepository'

const API_BASE_URL = '/api'

export class MatchRequestApiAdapter implements MatchRequestRepository {
  async create(requesterAccountId: string, matchOfferId: string): Promise<void> {
    const response = await fetch(
      `${API_BASE_URL}/account/${requesterAccountId}/match-offer/${matchOfferId}/match-request`,
      { method: 'POST' },
    )

    if (!response.ok) {
      const data = await response.json().catch(() => ({}))
      throw new Error(data.message || 'Error al crear la solicitud de partido')
    }
  }

  async findReceived(ownerAccountId: string, statuses?: string[]): Promise<MatchRequest[]> {
    const queryParams = new URLSearchParams()
    if (statuses && statuses.length > 0) {
      queryParams.append('statuses', statuses.join(','))
    }
    const queryString = queryParams.toString()
    const url = `${API_BASE_URL}/account/${ownerAccountId}/match-request${queryString ? `?${queryString}` : ''}`
    const response = await fetch(url)
    if (!response.ok) {
      const data = await response.json().catch(() => ({}))
      throw new Error(data.message || 'Error al obtener solicitudes recibidas')
    }
    return response.json()
  }

  async accept(ownerAccountId: string, requestId: string): Promise<void> {
    const response = await fetch(
      `${API_BASE_URL}/account/${ownerAccountId}/match-request/${encodeURIComponent(requestId)}/accept`,
      { method: 'POST' },
    )
    if (!response.ok) {
      const data = await response.json().catch(() => ({}))
      throw new Error(data.message || 'Error al aceptar la solicitud')
    }
  }

  async cancel(requesterAccountId: string, requestId: string): Promise<void> {
    const response = await fetch(
      `${API_BASE_URL}/account/${requesterAccountId}/match-request/${encodeURIComponent(requestId)}/cancel`,
      { method: 'POST' },
    )
    if (!response.ok) {
      const data = await response.json().catch(() => ({}))
      throw new Error(data.message || 'Error al cancelar la solicitud')
    }
  }

  async findSent(requesterAccountId: string, statuses?: string[]): Promise<MatchRequest[]> {
    const queryParams = new URLSearchParams({ role: 'requester' })
    if (statuses && statuses.length > 0) {
      queryParams.append('statuses', statuses.join(','))
    }
    const url = `${API_BASE_URL}/account/${requesterAccountId}/match-request?${queryParams.toString()}`

    const response = await fetch(url)
    if (!response.ok) {
      const data = await response.json().catch(() => ({}))
      throw new Error(data.message || 'Error al obtener solicitudes enviadas')
    }
    return response.json()
  }
}
