package usecases

import (
	"sportlink/api/domain/matchannouncement"
)

type FindMatchAnnouncementUC struct {
	matchAnnouncementRepository matchannouncement.Repository
}

func NewFindMatchAnnouncementUC(matchAnnouncementRepository matchannouncement.Repository) *FindMatchAnnouncementUC {
	return &FindMatchAnnouncementUC{
		matchAnnouncementRepository: matchAnnouncementRepository,
	}
}

func (uc *FindMatchAnnouncementUC) Invoke(query matchannouncement.DomainQuery) (*[]matchannouncement.Entity, error) {
	result, err := uc.matchAnnouncementRepository.Find(query)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
