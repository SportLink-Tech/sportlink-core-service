import { createContext, useContext, ReactNode } from 'react'
import { CreateMatchOfferUseCase } from '../domain/usecases/CreateMatchOfferUseCase'
import { FindMatchOffersUseCase } from '../domain/usecases/FindMatchOffersUseCase'
import { FindAccountMatchOffersUseCase } from '../domain/usecases/FindAccountMatchOffersUseCase'
import { DeleteMatchOfferUseCase } from '../domain/usecases/DeleteMatchOfferUseCase'
import { MatchOfferApiAdapter } from '../infrastructure/adapters/MatchOfferApiAdapter'

/**
 * MatchOffer Context Interface
 * Exposes use cases to UI components
 */
interface MatchOfferContextType {
  createMatchOfferUseCase: CreateMatchOfferUseCase
  findMatchOffersUseCase: FindMatchOffersUseCase
  findAccountMatchOffersUseCase: FindAccountMatchOffersUseCase
  deleteMatchOfferUseCase: DeleteMatchOfferUseCase
}

const MatchOfferContext = createContext<MatchOfferContextType | undefined>(undefined)

/**
 * MatchOffer Provider
 * Implements Dependency Injection
 * Following Hexagonal Architecture:
 * 1. Creates adapter (MatchOfferApiAdapter)
 * 2. Injects adapter into use cases
 * 3. Provides use cases to UI components
 */
export function MatchOfferProvider({ children }: { children: ReactNode }) {
  // Create adapter instance (Infrastructure layer)
  const matchOfferApiAdapter = new MatchOfferApiAdapter()

  // Wire dependencies: Inject adapter into use cases (Domain layer)
  // This applies the Dependency Inversion Principle
  const createMatchOfferUseCase = new CreateMatchOfferUseCase(matchOfferApiAdapter)
  const findMatchOffersUseCase = new FindMatchOffersUseCase(matchOfferApiAdapter)
  const findAccountMatchOffersUseCase = new FindAccountMatchOffersUseCase(matchOfferApiAdapter)
  const deleteMatchOfferUseCase = new DeleteMatchOfferUseCase(matchOfferApiAdapter)

  const value: MatchOfferContextType = {
    createMatchOfferUseCase,
    findMatchOffersUseCase,
    findAccountMatchOffersUseCase,
    deleteMatchOfferUseCase,
  }

  return <MatchOfferContext.Provider value={value}>{children}</MatchOfferContext.Provider>
}

/**
 * Custom Hook to use MatchOffer Context
 * UI components use this hook to access use cases
 */
export function useMatchOfferContext(): MatchOfferContextType {
  const context = useContext(MatchOfferContext)
  if (!context) {
    throw new Error('useMatchOfferContext must be used within MatchOfferProvider')
  }
  return context
}

