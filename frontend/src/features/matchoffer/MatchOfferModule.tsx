import { ReactNode } from 'react'
import { MatchOfferProvider } from './context/MatchOfferContext'

/**
 * MatchOffer Module
 * Wraps MatchOfferProvider to modularize the feature
 * This module can be imported anywhere in the app
 * Following Hexagonal Architecture - Module enforces modularity
 */
interface MatchOfferModuleProps {
  children: ReactNode
}

export function MatchOfferModule({ children }: MatchOfferModuleProps) {
  return <MatchOfferProvider>{children}</MatchOfferProvider>
}

