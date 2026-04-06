import { AuthRepository } from '../ports/AuthRepository'

export class GoogleAuthUseCase {
  constructor(private readonly authRepository: AuthRepository) {}

  async execute(idToken: string): Promise<{ accountId: string; success: boolean; error?: string }> {
    try {
      const result = await this.authRepository.googleLogin(idToken)
      return { accountId: result.accountId, success: true }
    } catch (error) {
      return {
        accountId: '',
        success: false,
        error: error instanceof Error ? error.message : 'Error al iniciar sesión',
      }
    }
  }
}
