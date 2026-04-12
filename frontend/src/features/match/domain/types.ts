export interface Match {
  id: string
  participants: string[]
  sport: string
  day: string
  status: 'ACCEPTED' | 'PLAYED' | 'CANCELLED'
  created_at: string
}
