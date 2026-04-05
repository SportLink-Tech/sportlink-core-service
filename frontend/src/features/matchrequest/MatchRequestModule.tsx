import { ReactNode } from 'react'
import { MatchRequestProvider } from './context/MatchRequestContext'

interface MatchRequestModuleProps {
  children: ReactNode
}

export function MatchRequestModule({ children }: MatchRequestModuleProps) {
  return <MatchRequestProvider>{children}</MatchRequestProvider>
}
