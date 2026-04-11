import { Match } from '../domain/types'

const API_BASE_URL = '/api'

export async function fetchMatches(accountId: string, statuses?: string[]): Promise<Match[]> {
  const params = new URLSearchParams()
  if (statuses && statuses.length > 0) {
    params.set('statuses', statuses.join(','))
  }

  const url = `${API_BASE_URL}/account/${accountId}/match${params.toString() ? '?' + params.toString() : ''}`
  const response = await fetch(url)

  if (!response.ok) {
    const data = await response.json().catch(() => ({}))
    throw new Error(data.message || 'Error al cargar los partidos')
  }

  return response.json()
}
