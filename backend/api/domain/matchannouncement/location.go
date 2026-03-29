package matchannouncement

import "time"

// Location represents the geographic location of a match
type Location struct {
	Country   string
	Province  string
	Locality  string
	Latitude  float64
	Longitude float64
}

func NewLocation(country, province, locality string) Location {
	return Location{
		Country:  country,
		Province: province,
		Locality: locality,
	}
}

func NewLocationWithCoords(country, province, locality string, latitude, longitude float64) Location {
	return Location{
		Country:   country,
		Province:  province,
		Locality:  locality,
		Latitude:  latitude,
		Longitude: longitude,
	}
}

// HasCoords returns true if the location has valid GPS coordinates
func (l Location) HasCoords() bool {
	return l.Latitude != 0 || l.Longitude != 0
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
