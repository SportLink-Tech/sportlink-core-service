import {
  MatchAnnouncement,
  CreateMatchAnnouncementRequest,
  FindMatchAnnouncementsQuery,
  PaginatedMatchAnnouncementsResponse,
} from '../../../../shared/types/matchAnnouncement'

export interface MatchAnnouncementRepository {
  create(request: CreateMatchAnnouncementRequest): Promise<{ data: MatchAnnouncement; status: number }>
  find(query: FindMatchAnnouncementsQuery): Promise<{ data: PaginatedMatchAnnouncementsResponse; status: number }>
}

