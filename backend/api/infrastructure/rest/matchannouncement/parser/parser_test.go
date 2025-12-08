package parser

import (
	"sportlink/api/domain/common"
	"sportlink/api/domain/matchannouncement"
	"testing"
	"time"
)

func TestDefaultQueryParser_Sports(t *testing.T) {
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
			got, err := parser.Sports(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Sports() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !equalSports(got, tt.want) {
				t.Errorf("Sports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultQueryParser_Categories(t *testing.T) {
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
			got, err := parser.Categories(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Categories() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !equalCategories(got, tt.want) {
				t.Errorf("Categories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultQueryParser_Statuses(t *testing.T) {
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
			got, err := parser.Statuses(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Statuses() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !equalStatuses(got, tt.want) {
				t.Errorf("Statuses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultQueryParser_Date(t *testing.T) {
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
			got, err := parser.Date(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Date() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !got.Equal(tt.want) && !tt.want.IsZero() {
				t.Errorf("Date() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultQueryParser_Location(t *testing.T) {
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
			got := parser.Location(tt.country, tt.province, tt.locality)
			if tt.want == nil {
				if got != nil {
					t.Errorf("Location() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Errorf("Location() = nil, want %v", tt.want)
				return
			}
			if got.Country != tt.want.Country || got.Province != tt.want.Province || got.Locality != tt.want.Locality {
				t.Errorf("Location() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultQueryParser_Limit(t *testing.T) {
	parser := NewQueryParser()

	tests := []struct {
		name    string
		limit   string
		want    int
		wantErr bool
	}{
		{
			name:    "given valid limit when parsing then returns limit value",
			limit:   "9",
			want:    9,
			wantErr: false,
		},
		{
			name:    "given empty string when parsing limit then returns zero",
			limit:   "",
			want:    0,
			wantErr: false,
		},
		{
			name:    "given zero when parsing limit then returns zero",
			limit:   "0",
			want:    0,
			wantErr: false,
		},
		{
			name:    "given large number when parsing limit then returns value",
			limit:   "100",
			want:    100,
			wantErr: false,
		},
		{
			name:    "given invalid string when parsing limit then returns error",
			limit:   "invalid",
			want:    0,
			wantErr: true,
		},
		{
			name:    "given negative number when parsing limit then returns error",
			limit:   "-5",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.Limit(tt.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Limit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Limit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultQueryParser_Offset(t *testing.T) {
	parser := NewQueryParser()

	tests := []struct {
		name    string
		offset  string
		want    int
		wantErr bool
	}{
		{
			name:    "given valid offset when parsing then returns offset value",
			offset:  "9",
			want:    9,
			wantErr: false,
		},
		{
			name:    "given empty string when parsing offset then returns zero",
			offset:  "",
			want:    0,
			wantErr: false,
		},
		{
			name:    "given zero when parsing offset then returns zero",
			offset:  "0",
			want:    0,
			wantErr: false,
		},
		{
			name:    "given large number when parsing offset then returns value",
			offset:  "100",
			want:    100,
			wantErr: false,
		},
		{
			name:    "given invalid string when parsing offset then returns error",
			offset:  "invalid",
			want:    0,
			wantErr: true,
		},
		{
			name:    "given negative number when parsing offset then returns error",
			offset:  "-5",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.Offset(tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("Offset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Offset() = %v, want %v", got, tt.want)
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
