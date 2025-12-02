package matchannouncement

import "fmt"

// Status represents the state of a match announcement
type Status string

const (
	StatusPending   Status = "PENDING"   // Announcement published, waiting for responses
	StatusConfirmed Status = "CONFIRMED" // Match confirmed with another team
	StatusCancelled Status = "CANCELLED" // Announcement cancelled by the team
	StatusExpired   Status = "EXPIRED"   // Announcement expired by TTL
)

// AllStatus returns all valid statuses
func AllStatus() []Status {
	return []Status{
		StatusPending,
		StatusConfirmed,
		StatusCancelled,
		StatusExpired,
	}
}

// IsValid checks if a status is valid
func (s Status) IsValid() bool {
	switch s {
	case StatusPending, StatusConfirmed, StatusCancelled, StatusExpired:
		return true
	default:
		return false
	}
}

// String returns the string representation of the status
func (s Status) String() string {
	return string(s)
}

// ParseStatus converts a string to Status
func ParseStatus(s string) (Status, error) {
	status := Status(s)
	if !status.IsValid() {
		return "", fmt.Errorf("invalid status: %s", s)
	}
	return status, nil
}
