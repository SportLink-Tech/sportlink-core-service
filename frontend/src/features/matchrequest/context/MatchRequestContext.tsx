import { createContext, useContext, ReactNode } from 'react'
import { CreateMatchRequestUseCase } from '../domain/usecases/CreateMatchRequestUseCase'
import { FindSentMatchRequestsUseCase } from '../domain/usecases/FindSentMatchRequestsUseCase'
import { FindSentMatchRequestsListUseCase } from '../domain/usecases/FindSentMatchRequestsListUseCase'
import { FindReceivedMatchRequestsUseCase } from '../domain/usecases/FindReceivedMatchRequestsUseCase'
import { AcceptMatchRequestUseCase } from '../domain/usecases/AcceptMatchRequestUseCase'
import { MatchRequestApiAdapter } from '../infrastructure/adapters/MatchRequestApiAdapter'

interface MatchRequestContextType {
  createMatchRequestUseCase: CreateMatchRequestUseCase
  findSentMatchRequestsUseCase: FindSentMatchRequestsUseCase
  findSentMatchRequestsListUseCase: FindSentMatchRequestsListUseCase
  findReceivedMatchRequestsUseCase: FindReceivedMatchRequestsUseCase
  acceptMatchRequestUseCase: AcceptMatchRequestUseCase
}

const MatchRequestContext = createContext<MatchRequestContextType | undefined>(undefined)

export function MatchRequestProvider({ children }: { children: ReactNode }) {
  const adapter = new MatchRequestApiAdapter()
  const createMatchRequestUseCase = new CreateMatchRequestUseCase(adapter)
  const findSentMatchRequestsUseCase = new FindSentMatchRequestsUseCase(adapter)
  const findSentMatchRequestsListUseCase = new FindSentMatchRequestsListUseCase(adapter)
  const findReceivedMatchRequestsUseCase = new FindReceivedMatchRequestsUseCase(adapter)
  const acceptMatchRequestUseCase = new AcceptMatchRequestUseCase(adapter)

  return (
    <MatchRequestContext.Provider value={{ createMatchRequestUseCase, findSentMatchRequestsUseCase, findSentMatchRequestsListUseCase, findReceivedMatchRequestsUseCase, acceptMatchRequestUseCase }}>
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
