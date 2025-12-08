package team

import (
	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
	"sportlink/api/domain/team"
)

// TODO faltan las stats
type Dto struct {
	EntityId string `dynamodbav:"EntityId"`
	Id       string `dynamodbav:"Id"`
	Category int    `dynamodbav:"Category"`
	Sport    string `dynamodbav:"Sport"`
}

func (d *Dto) ToDomain() team.Entity {
	// Extract name from ID format: SPORT#<sport>#NAME#<name>
	// For backward compatibility, if ID doesn't match the format, use ID as name
	name := extractNameFromID(d.Id)

	return team.Entity{
		ID:       d.Id,
		Name:     name,
		Category: common.Category(d.Category),
		Sport:    common.Sport(d.Sport),
		Stats:    *common.NewStats(0, 0, 0), // Default stats (not persisted yet)
		Members:  []player.Entity{},         // Default empty members (not persisted yet)
	}
}

// extractNameFromID extracts the team name from the ID format SPORT#<sport>#NAME#<name>
// Returns the original ID if it doesn't match the format (for backward compatibility)
func extractNameFromID(id string) string {
	// Expected format: SPORT#<sport>#NAME#<name>
	prefix := "SPORT#"
	namePrefix := "#NAME#"

	if !contains(id, prefix) || !contains(id, namePrefix) {
		// Backward compatibility: if ID doesn't match format, assume it's just the name
		return id
	}

	// Find the position after #NAME#
	nameIndex := indexOf(id, namePrefix) + len(namePrefix)
	if nameIndex >= len(id) {
		return id
	}

	return id[nameIndex:]
}

// Helper functions for string operations
func contains(s, substr string) bool {
	return indexOf(s, substr) >= 0
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
