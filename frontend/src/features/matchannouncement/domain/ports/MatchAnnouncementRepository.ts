import { MatchAnnouncement, CreateMatchAnnouncementRequest, FindMatchAnnouncementsQuery } from '../../../../shared/types/matchAnnouncement'

export interface MatchAnnouncementRepository {
  create(request: CreateMatchAnnouncementRequest): Promise<{ data: MatchAnnouncement; status: number }>
  find(query: FindMatchAnnouncementsQuery): Promise<{ data: MatchAnnouncement[]; status: number }>
}

