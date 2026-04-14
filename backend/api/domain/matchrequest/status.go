package matchrequest

import "fmt"

type Status string

const (
	StatusPending  Status = "PENDING"
	StatusAccepted Status = "ACCEPTED"
	StatusCancel   Status = "CANCEL"
	StatusRejected Status = "REJECTED"
)

func (s Status) String() string {
	return string(s)
}

func (s Status) IsValid() bool {
	switch s {
	case StatusPending, StatusAccepted, StatusRejected, StatusCancel:
		return true
	}
	return false
}

func ParseStatus(s string) (Status, error) {
	status := Status(s)
	if !status.IsValid() {
		return "", fmt.Errorf("invalid match request status: %s", s)
	}
	return status, nil
}
