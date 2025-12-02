package parser

import (
	"sportlink/api/domain/common"
	"sportlink/api/domain/matchannouncement"
	"testing"
	"time"
)

func TestDefaultQueryParser_ParseSports(t *testing.T) {
	parser := NewQueryParser()

	tests := []struct {
		name      string
		input     string
		want      []common.Sport
		wantError bool
	}{
		{
			name:      "empty string returns nil",
			input:     "",
			want:      nil,
			wantError: false,
		},
		{
			name:      "single sport",
			input:     "Football",
			want:      []common.Sport{common.Sport("Football")},
			wantError: false,
		},
		{
			name:      "multiple sports",
			input:     "Football,Paddle,Tennis",
			want:      []common.Sport{common.Sport("Football"), common.Sport("Paddle"), common.Sport("Tennis")},
			wantError: false,
		},
		{
			name:      "sports with spaces",
			input:     "Football, Paddle , Tennis",
			want:      []common.Sport{common.Sport("Football"), common.Sport("Paddle"), common.Sport("Tennis")},
			wantError: false,
		},
		{
			name:      "empty values are skipped",
			input:     "Football,,Paddle",
			want:      []common.Sport{common.Sport("Football"), common.Sport("Paddle")},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ParseSports(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("ParseSports() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !equalSports(got, tt.want) {
				t.Errorf("ParseSports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultQueryParser_ParseCategories(t *testing.T) {
	parser := NewQueryParser()

	tests := []struct {
		name      string
		input     string
		want      []common.Category
		wantError bool
	}{
		{
			name:      "empty string returns nil",
			input:     "",
			want:      nil,
			wantError: false,
		},
		{
			name:      "single category",
			input:     "5",
			want:      []common.Category{common.L5},
			wantError: false,
		},
		{
			name:      "multiple categories",
			input:     "1,3,5",
			want:      []common.Category{common.L1, common.L3, common.L5},
			wantError: false,
		},
		{
			name:      "categories with spaces",
			input:     "1, 3 , 5",
			want:      []common.Category{common.L1, common.L3, common.L5},
			wantError: false,
		},
		{
			name:      "invalid category format",
			input:     "invalid",
			want:      nil,
			wantError: true,
		},
		{
			name:      "invalid category value",
			input:     "99",
			want:      nil,
			wantError: true,
		},
		{
			name:      "empty values are skipped",
			input:     "1,,3",
			want:      []common.Category{common.L1, common.L3},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ParseCategories(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("ParseCategories() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !equalCategories(got, tt.want) {
				t.Errorf("ParseCategories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultQueryParser_ParseStatuses(t *testing.T) {
	parser := NewQueryParser()

	tests := []struct {
		name      string
		input     string
		want      []matchannouncement.Status
		wantError bool
	}{
		{
			name:      "empty string returns nil",
			input:     "",
			want:      nil,
			wantError: false,
		},
		{
			name:      "single status",
			input:     "PENDING",
			want:      []matchannouncement.Status{matchannouncement.StatusPending},
			wantError: false,
		},
		{
			name:      "multiple statuses",
			input:     "PENDING,CONFIRMED",
			want:      []matchannouncement.Status{matchannouncement.StatusPending, matchannouncement.StatusConfirmed},
			wantError: false,
		},
		{
			name:      "statuses with spaces",
			input:     "PENDING, CONFIRMED , CANCELLED",
			want:      []matchannouncement.Status{matchannouncement.StatusPending, matchannouncement.StatusConfirmed, matchannouncement.StatusCancelled},
			wantError: false,
		},
		{
			name:      "invalid status",
			input:     "INVALID",
			want:      nil,
			wantError: true,
		},
		{
			name:      "empty values are skipped",
			input:     "PENDING,,CONFIRMED",
			want:      []matchannouncement.Status{matchannouncement.StatusPending, matchannouncement.StatusConfirmed},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ParseStatuses(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("ParseStatuses() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !equalStatuses(got, tt.want) {
				t.Errorf("ParseStatuses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultQueryParser_ParseDate(t *testing.T) {
	parser := NewQueryParser()

	tests := []struct {
		name      string
		input     string
		want      time.Time
		wantError bool
	}{
		{
			name:      "empty string returns zero time",
			input:     "",
			want:      time.Time{},
			wantError: false,
		},
		{
			name:      "valid date",
			input:     "2025-12-01",
			want:      time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
			wantError: false,
		},
		{
			name:      "invalid date format",
			input:     "01-12-2025",
			want:      time.Time{},
			wantError: true,
		},
		{
			name:      "invalid date value",
			input:     "2025-13-01",
			want:      time.Time{},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ParseDate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("ParseDate() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !got.Equal(tt.want) && !tt.want.IsZero() {
				t.Errorf("ParseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultQueryParser_ParseLocation(t *testing.T) {
	parser := NewQueryParser()

	tests := []struct {
		name     string
		country  string
		province string
		locality string
		want     *matchannouncement.Location
	}{
		{
			name:     "all empty returns nil",
			country:  "",
			province: "",
			locality: "",
			want:     nil,
		},
		{
			name:     "all fields populated",
			country:  "Argentina",
			province: "Buenos Aires",
			locality: "Palermo",
			want: &matchannouncement.Location{
				Country:  "Argentina",
				Province: "Buenos Aires",
				Locality: "Palermo",
			},
		},
		{
			name:     "partial fields",
			country:  "Argentina",
			province: "",
			locality: "Palermo",
			want: &matchannouncement.Location{
				Country:  "Argentina",
				Province: "",
				Locality: "Palermo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.ParseLocation(tt.country, tt.province, tt.locality)
			if tt.want == nil {
				if got != nil {
					t.Errorf("ParseLocation() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Errorf("ParseLocation() = nil, want %v", tt.want)
				return
			}
			if got.Country != tt.want.Country || got.Province != tt.want.Province || got.Locality != tt.want.Locality {
				t.Errorf("ParseLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper functions for comparing slices

func equalSports(a, b []common.Sport) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalCategories(a, b []common.Category) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalStatuses(a, b []matchannouncement.Status) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
