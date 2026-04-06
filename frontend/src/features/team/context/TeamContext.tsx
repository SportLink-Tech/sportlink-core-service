import { createContext, useContext, ReactNode } from 'react'
import { CreateTeamUseCase } from '../domain/usecases/CreateTeamUseCase'
import { SearchTeamUseCase } from '../domain/usecases/SearchTeamUseCase'
import { ListAccountTeamsUseCase } from '../domain/usecases/ListAccountTeamsUseCase'
import { UpdateTeamUseCase } from '../domain/usecases/UpdateTeamUseCase'
import { TeamApiAdapter } from '../infrastructure/adapters/TeamApiAdapter'

/**
 * Team Context Interface
 * Exposes use cases to UI components
 */
interface TeamContextType {
  createTeamUseCase: CreateTeamUseCase
  searchTeamUseCase: SearchTeamUseCase
  listAccountTeamsUseCase: ListAccountTeamsUseCase
  updateTeamUseCase: UpdateTeamUseCase
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
  const teamApiAdapter = new TeamApiAdapter()

  const createTeamUseCase = new CreateTeamUseCase(teamApiAdapter)
  const searchTeamUseCase = new SearchTeamUseCase(teamApiAdapter)
  const listAccountTeamsUseCase = new ListAccountTeamsUseCase(teamApiAdapter)
  const updateTeamUseCase = new UpdateTeamUseCase(teamApiAdapter)

  const value: TeamContextType = {
    createTeamUseCase,
    searchTeamUseCase,
    listAccountTeamsUseCase,
    updateTeamUseCase,
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
