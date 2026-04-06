import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { GoogleAuthUseCase } from '../domain/usecases/GoogleAuthUseCase'
import { AuthApiAdapter } from '../infrastructure/adapters/AuthApiAdapter'
import { fetchAccount, Account } from '../infrastructure/adapters/AccountApiAdapter'

interface AuthContextType {
  accountId: string | null
  account: Account | null
  googleLogin: (idToken: string) => Promise<{ success: boolean; error?: string }>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

const ACCOUNT_ID_KEY = 'account_id'

export function AuthProvider({ children }: { children: ReactNode }) {
  const [accountId, setAccountId] = useState<string | null>(
    () => localStorage.getItem(ACCOUNT_ID_KEY)
  )
  const [account, setAccount] = useState<Account | null>(null)

  const googleAuthUseCase = new GoogleAuthUseCase(new AuthApiAdapter())

  useEffect(() => {
    if (accountId && !account) {
      fetchAccount(accountId)
        .then(setAccount)
        .catch(() => {})
    }
  }, [accountId])

  const googleLogin = async (idToken: string) => {
    const result = await googleAuthUseCase.execute(idToken)
    if (result.success) {
      setAccountId(result.accountId)
      localStorage.setItem(ACCOUNT_ID_KEY, result.accountId)
      try {
        const fetchedAccount = await fetchAccount(result.accountId)
        setAccount(fetchedAccount)
      } catch {
        // account fetch failure should not block login
      }
    }
    return { success: result.success, error: result.error }
  }

  const logout = () => {
    setAccountId(null)
    setAccount(null)
    localStorage.removeItem(ACCOUNT_ID_KEY)
  }

  return (
    <AuthContext.Provider value={{ accountId, account, googleLogin, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthContextType {
  const context = useContext(AuthContext)
  if (!context) throw new Error('useAuth must be used within AuthProvider')
  return context
}
