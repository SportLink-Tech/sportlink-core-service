export interface Match {
  id: string
  participants: string[]
  sport: string
  day: string
  time_slot?: {
    start_time: string
    end_time: string
  }
  title: string
  status: 'ACCEPTED' | 'PLAYED' | 'CANCELLED'
  created_at: string
}
