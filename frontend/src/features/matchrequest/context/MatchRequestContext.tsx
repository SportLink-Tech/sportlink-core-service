import { createContext, useContext, ReactNode } from 'react'
import { CreateMatchRequestUseCase } from '../domain/usecases/CreateMatchRequestUseCase'
import { FindSentMatchRequestsUseCase } from '../domain/usecases/FindSentMatchRequestsUseCase'
import { MatchRequestApiAdapter } from '../infrastructure/adapters/MatchRequestApiAdapter'

interface MatchRequestContextType {
  createMatchRequestUseCase: CreateMatchRequestUseCase
  findSentMatchRequestsUseCase: FindSentMatchRequestsUseCase
}

const MatchRequestContext = createContext<MatchRequestContextType | undefined>(undefined)

export function MatchRequestProvider({ children }: { children: ReactNode }) {
  const adapter = new MatchRequestApiAdapter()
  const createMatchRequestUseCase = new CreateMatchRequestUseCase(adapter)
  const findSentMatchRequestsUseCase = new FindSentMatchRequestsUseCase(adapter)

  return (
    <MatchRequestContext.Provider value={{ createMatchRequestUseCase, findSentMatchRequestsUseCase }}>
      {children}
    </MatchRequestContext.Provider>
  )
}

export function useMatchRequestContext(): MatchRequestContextType {
  const context = useContext(MatchRequestContext)
  if (!context) {
    throw new Error('useMatchRequestContext must be used within MatchRequestProvider')
  }
  return context
}
