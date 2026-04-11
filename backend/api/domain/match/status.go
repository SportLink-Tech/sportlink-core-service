package match

import "fmt"

type Status string

const (
	StatusAccepted  Status = "ACCEPTED"  // Partido confirmado, aún no jugado
	StatusPlayed    Status = "PLAYED"    // Partido jugado, resultado registrado
	StatusCancelled Status = "CANCELLED" // Partido cancelado
)

func (s Status) IsValid() bool {
	switch s {
	case StatusAccepted, StatusPlayed, StatusCancelled:
		return true
	default:
		return false
	}
}

func (s Status) String() string {
	return string(s)
}

func ParseStatus(s string) (Status, error) {
	status := Status(s)
	if !status.IsValid() {
		return "", fmt.Errorf("invalid match status: %s", s)
	}
	return status, nil
}
