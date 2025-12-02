package matchannouncement

import "time"

// Location represents the geographic location of a match
type Location struct {
	Country  string
	Province string
	Locality string
}

func NewLocation(country, province, locality string) Location {
	return Location{
		Country:  country,
		Province: province,
		Locality: locality,
	}
}

// GetTimezone returns the timezone associated with the location
// By default returns GMT-3 (Argentina/Buenos Aires)
// In the future it can be extended to map different locations to their corresponding timezones
func (l Location) GetTimezone() *time.Location {
	// Default timezone: GMT-3 (Argentina)
	location, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		// Fallback to fixed GMT-3 offset
		return time.FixedZone("GMT-3", -3*60*60)
	}
	return location
}
