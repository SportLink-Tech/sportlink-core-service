export interface Match {
  id: string
  local_account_id: string
  visitor_account_id: string
  sport: string
  day: string
  status: 'ACCEPTED' | 'PLAYED' | 'CANCELLED'
  created_at: string
}
