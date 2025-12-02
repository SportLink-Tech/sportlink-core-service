import { createContext, useContext, ReactNode } from 'react'
import { CreateMatchAnnouncementUseCase } from '../domain/usecases/CreateMatchAnnouncementUseCase'
import { FindMatchAnnouncementsUseCase } from '../domain/usecases/FindMatchAnnouncementsUseCase'
import { MatchAnnouncementApiAdapter } from '../infrastructure/adapters/MatchAnnouncementApiAdapter'

/**
 * MatchAnnouncement Context Interface
 * Exposes use cases to UI components
 */
interface MatchAnnouncementContextType {
  createMatchAnnouncementUseCase: CreateMatchAnnouncementUseCase
  findMatchAnnouncementsUseCase: FindMatchAnnouncementsUseCase
}

const MatchAnnouncementContext = createContext<MatchAnnouncementContextType | undefined>(undefined)

/**
 * MatchAnnouncement Provider
 * Implements Dependency Injection
 * Following Hexagonal Architecture:
 * 1. Creates adapter (MatchAnnouncementApiAdapter)
 * 2. Injects adapter into use cases
 * 3. Provides use cases to UI components
 */
export function MatchAnnouncementProvider({ children }: { children: ReactNode }) {
  // Create adapter instance (Infrastructure layer)
  const matchAnnouncementApiAdapter = new MatchAnnouncementApiAdapter()

  // Wire dependencies: Inject adapter into use cases (Domain layer)
  // This applies the Dependency Inversion Principle
  const createMatchAnnouncementUseCase = new CreateMatchAnnouncementUseCase(matchAnnouncementApiAdapter)
  const findMatchAnnouncementsUseCase = new FindMatchAnnouncementsUseCase(matchAnnouncementApiAdapter)

  const value: MatchAnnouncementContextType = {
    createMatchAnnouncementUseCase,
    findMatchAnnouncementsUseCase,
  }

  return <MatchAnnouncementContext.Provider value={value}>{children}</MatchAnnouncementContext.Provider>
}

/**
 * Custom Hook to use MatchAnnouncement Context
 * UI components use this hook to access use cases
 */
export function useMatchAnnouncementContext(): MatchAnnouncementContextType {
  const context = useContext(MatchAnnouncementContext)
  if (!context) {
    throw new Error('useMatchAnnouncementContext must be used within MatchAnnouncementProvider')
  }
  return context
}

