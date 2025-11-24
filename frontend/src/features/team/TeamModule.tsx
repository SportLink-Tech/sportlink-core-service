import { ReactNode } from 'react'
import { TeamProvider } from './context/TeamContext'

/**
 * Team Module
 * Wraps TeamProvider to modularize the feature
 * This module can be imported anywhere in the app
 * Following Hexagonal Architecture - Module enforces modularity
 */
interface TeamModuleProps {
  children: ReactNode
}

export function TeamModule({ children }: TeamModuleProps) {
  return <TeamProvider>{children}</TeamProvider>
}

