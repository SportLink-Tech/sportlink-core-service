import { AuthRepository } from '../../domain/ports/AuthRepository'

const API_BASE_URL = '/api'

export class AuthApiAdapter implements AuthRepository {
  async googleLogin(idToken: string): Promise<{ accountId: string }> {
    const response = await fetch(`${API_BASE_URL}/auth/google`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include', // sends/receives cookies
      body: JSON.stringify({ id_token: idToken }),
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.message || 'Error al iniciar sesión con Google')
    }

    return { accountId: data.account_id }
  }
}
