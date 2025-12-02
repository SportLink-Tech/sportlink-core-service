import { createContext, useContext, ReactNode } from 'react'
import { CreateTeamUseCase } from '../domain/usecases/CreateTeamUseCase'
import { SearchTeamUseCase } from '../domain/usecases/SearchTeamUseCase'
import { TeamApiAdapter } from '../infrastructure/adapters/TeamApiAdapter'

/**
 * Team Context Interface
 * Exposes use cases to UI components
 */
interface TeamContextType {
  createTeamUseCase: CreateTeamUseCase
  searchTeamUseCase: SearchTeamUseCase
}

const TeamContext = createContext<TeamContextType | undefined>(undefined)

/**
 * Team Provider
 * Implements Dependency Injection
 * Following Hexagonal Architecture:
 * 1. Creates adapter (TeamApiAdapter)
 * 2. Injects adapter into use cases
 * 3. Provides use cases to UI components
 */
export function TeamProvider({ children }: { children: ReactNode }) {
  // Create adapter instance (Infrastructure layer)
  const teamApiAdapter = new TeamApiAdapter()

  // Wire dependencies: Inject adapter into use cases (Domain layer)
  // This applies the Dependency Inversion Principle
  const createTeamUseCase = new CreateTeamUseCase(teamApiAdapter)
  const searchTeamUseCase = new SearchTeamUseCase(teamApiAdapter)

  const value: TeamContextType = {
    createTeamUseCase,
    searchTeamUseCase,
  }

  return <TeamContext.Provider value={value}>{children}</TeamContext.Provider>
}

/**
 * Custom Hook to use Team Context
 * UI components use this hook to access use cases
 */
export function useTeamContext(): TeamContextType {
  const context = useContext(TeamContext)
  if (!context) {
    throw new Error('useTeamContext must be used within TeamProvider')
  }
  return context
}

