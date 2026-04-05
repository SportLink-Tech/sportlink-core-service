export interface MatchRequest {
  id: string
  match_announcement_id: string
  owner_account_id: string
  requester_account_id: string
  status: 'PENDING' | 'ACCEPTED' | 'REJECTED'
  created_at: string
}

export interface MatchRequestRepository {
  create(requesterAccountId: string, matchAnnouncementId: string): Promise<void>
  findSent(requesterAccountId: string, statuses?: string[]): Promise<MatchRequest[]>
}
