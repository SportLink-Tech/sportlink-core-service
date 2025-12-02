import { ReactNode } from 'react'
import { MatchAnnouncementProvider } from './context/MatchAnnouncementContext'

/**
 * MatchAnnouncement Module
 * Wraps MatchAnnouncementProvider to modularize the feature
 * This module can be imported anywhere in the app
 * Following Hexagonal Architecture - Module enforces modularity
 */
interface MatchAnnouncementModuleProps {
  children: ReactNode
}

export function MatchAnnouncementModule({ children }: MatchAnnouncementModuleProps) {
  return <MatchAnnouncementProvider>{children}</MatchAnnouncementProvider>
}

