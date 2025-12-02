package usecases

import (
	"fmt"
	"sportlink/api/domain/common"
	"sportlink/api/domain/matchannouncement"
	"sportlink/api/domain/team"
	"time"
)

type CreateMatchAnnouncementUC struct {
	matchAnnouncementRepository matchannouncement.Repository
	teamRepository              team.Repository
}

func NewCreateMatchAnnouncementUC(matchAnnouncementRepository matchannouncement.Repository, teamRepository team.Repository) *CreateMatchAnnouncementUC {
	return &CreateMatchAnnouncementUC{
		matchAnnouncementRepository: matchAnnouncementRepository,
		teamRepository:              teamRepository,
	}
}

func (uc *CreateMatchAnnouncementUC) Invoke(input matchannouncement.Entity) (*matchannouncement.Entity, error) {
	// Validate the announcement
	if err := uc.validateAnnouncement(input); err != nil {
		return nil, err
	}

	// Validate that the team exists
	if err := uc.validateTeamExists(input); err != nil {
		return nil, err
	}

	// Save the announcement
	err := uc.matchAnnouncementRepository.Save(input)
	if err != nil {
		return nil, fmt.Errorf("error while inserting match announcement in database: %w", err)
	}

	return &input, nil
}

func (uc *CreateMatchAnnouncementUC) validateAnnouncement(input matchannouncement.Entity) error {
	// Validate team name
	if input.TeamName == "" {
		return fmt.Errorf("team name cannot be empty")
	}

	// Validate sport
	if input.Sport == "" {
		return fmt.Errorf("sport cannot be empty")
	}

	// Validate day is not in the past
	now := time.Now().In(input.Location.GetTimezone())
	if input.Day.Before(now.Truncate(24 * time.Hour)) {
		return fmt.Errorf("day cannot be in the past")
	}

	// Validate time slot
	if input.TimeSlot.StartTime.IsZero() || input.TimeSlot.EndTime.IsZero() {
		return fmt.Errorf("time slot cannot be empty")
	}

	if input.TimeSlot.EndTime.Before(input.TimeSlot.StartTime) {
		return fmt.Errorf("end time cannot be before start time")
	}

	// Validate location
	if input.Location.Country == "" || input.Location.Province == "" || input.Location.Locality == "" {
		return fmt.Errorf("location must have country, province and locality")
	}

	// Validate status
	if !input.Status.IsValid() {
		return fmt.Errorf("invalid status")
	}

	// Validate created at
	if input.CreatedAt.IsZero() {
		return fmt.Errorf("created at cannot be empty")
	}

	return nil
}

func (uc *CreateMatchAnnouncementUC) validateTeamExists(input matchannouncement.Entity) error {
	// Search for the team by name and sport
	teams, err := uc.teamRepository.Find(team.DomainQuery{
		Name:   input.TeamName,
		Sports: []common.Sport{input.Sport},
	})
	if err != nil {
		return fmt.Errorf("error while finding team: %w", err)
	}

	// Check if the team exists
	if len(teams) == 0 {
		return fmt.Errorf("team '%s' for sport '%s' does not exist", input.TeamName, input.Sport)
	}

	return nil
}
