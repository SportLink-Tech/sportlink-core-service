export interface AuthRepository {
  googleLogin(idToken: string): Promise<{ accountId: string }>
}
