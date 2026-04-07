import {
  MatchOffer,
  CreateMatchOfferRequest,
  FindMatchOffersQuery,
  PaginatedMatchOffersResponse,
} from '../../../../shared/types/matchOffer'

export interface MatchOfferRepository {
  create(accountId: string, request: CreateMatchOfferRequest): Promise<{ data: MatchOffer; status: number }>
  find(query: FindMatchOffersQuery): Promise<{ data: PaginatedMatchOffersResponse; status: number }>
  findByAccount(accountId: string, statuses?: string[]): Promise<{ data: MatchOffer[]; status: number }>
  retrieve(offerId: string): Promise<MatchOffer>
  delete(accountId: string, offerId: string): Promise<void>
}

