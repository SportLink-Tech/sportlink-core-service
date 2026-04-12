import { createContext, useContext, ReactNode } from 'react'
import { CreateMatchOfferUseCase } from '../domain/usecases/CreateMatchOfferUseCase'
import { FindMatchOffersUseCase } from '../domain/usecases/FindMatchOffersUseCase'
import { FindAccountMatchOffersUseCase } from '../domain/usecases/FindAccountMatchOffersUseCase'
import { RetrieveMatchOfferUseCase } from '../domain/usecases/RetrieveMatchOfferUseCase'
import { DeleteMatchOfferUseCase } from '../domain/usecases/DeleteMatchOfferUseCase'
import { ConfirmMatchOfferUseCase } from '../domain/usecases/ConfirmMatchOfferUseCase'
import { MatchOfferApiAdapter } from '../infrastructure/adapters/MatchOfferApiAdapter'

interface MatchOfferContextType {
  createMatchOfferUseCase: CreateMatchOfferUseCase
  findMatchOffersUseCase: FindMatchOffersUseCase
  findAccountMatchOffersUseCase: FindAccountMatchOffersUseCase
  retrieveMatchOfferUseCase: RetrieveMatchOfferUseCase
  deleteMatchOfferUseCase: DeleteMatchOfferUseCase
  confirmMatchOfferUseCase: ConfirmMatchOfferUseCase
}

const MatchOfferContext = createContext<MatchOfferContextType | undefined>(undefined)

export function MatchOfferProvider({ children }: { children: ReactNode }) {
  const matchOfferApiAdapter = new MatchOfferApiAdapter()

  const createMatchOfferUseCase = new CreateMatchOfferUseCase(matchOfferApiAdapter)
  const findMatchOffersUseCase = new FindMatchOffersUseCase(matchOfferApiAdapter)
  const findAccountMatchOffersUseCase = new FindAccountMatchOffersUseCase(matchOfferApiAdapter)
  const retrieveMatchOfferUseCase = new RetrieveMatchOfferUseCase(matchOfferApiAdapter)
  const deleteMatchOfferUseCase = new DeleteMatchOfferUseCase(matchOfferApiAdapter)
  const confirmMatchOfferUseCase = new ConfirmMatchOfferUseCase(matchOfferApiAdapter)

  const value: MatchOfferContextType = {
    createMatchOfferUseCase,
    findMatchOffersUseCase,
    findAccountMatchOffersUseCase,
    retrieveMatchOfferUseCase,
    deleteMatchOfferUseCase,
    confirmMatchOfferUseCase,
  }

  return <MatchOfferContext.Provider value={value}>{children}</MatchOfferContext.Provider>
}

export function useMatchOfferContext(): MatchOfferContextType {
  const context = useContext(MatchOfferContext)
  if (!context) {
    throw new Error('useMatchOfferContext must be used within MatchOfferProvider')
  }
  return context
}
