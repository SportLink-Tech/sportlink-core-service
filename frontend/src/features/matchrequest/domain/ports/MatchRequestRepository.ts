export interface MatchRequest {
  id: string
  match_offer_id: string
  owner_account_id: string
  requester_account_id: string
  status: 'PENDING' | 'ACCEPTED' | 'REJECTED' | 'CANCEL'
  created_at: string
}

export interface MatchRequestRepository {
  create(requesterAccountId: string, matchOfferId: string): Promise<void>
  findSent(requesterAccountId: string, statuses?: string[]): Promise<MatchRequest[]>
  findReceived(ownerAccountId: string, statuses?: string[]): Promise<MatchRequest[]>
  accept(ownerAccountId: string, requestId: string): Promise<void>
  cancel(requesterAccountId: string, requestId: string): Promise<void>
}
