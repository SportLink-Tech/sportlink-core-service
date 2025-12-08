package usecases_test

import (
	"sportlink/api/application/matchannouncement/usecases"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePageInfo(t *testing.T) {
	tests := []struct {
		name       string
		limit      int
		offset     int
		total      int
		wantNumber int
		wantOutOf  int
		wantTotal  int
	}{
		{
			name:       "given first page when calculating with limit and offset then returns page 1",
			limit:      9,
			offset:     0,
			total:      25,
			wantNumber: 1,
			wantOutOf:  3,
			wantTotal:  25,
		},
		{
			name:       "given second page when calculating with limit and offset then returns page 2",
			limit:      9,
			offset:     9,
			total:      25,
			wantNumber: 2,
			wantOutOf:  3,
			wantTotal:  25,
		},
		{
			name:       "given third page when calculating with limit and offset then returns page 3",
			limit:      9,
			offset:     18,
			total:      25,
			wantNumber: 3,
			wantOutOf:  3,
			wantTotal:  25,
		},
		{
			name:       "given exact division when calculating pages then returns correct total pages",
			limit:      9,
			offset:     0,
			total:      18,
			wantNumber: 1,
			wantOutOf:  2,
			wantTotal:  18,
		},
		{
			name:       "given remainder when calculating pages then rounds up total pages",
			limit:      9,
			offset:     0,
			total:      19,
			wantNumber: 1,
			wantOutOf:  3,
			wantTotal:  19,
		},
		{
			name:       "given no limit when calculating pages then returns single page",
			limit:      0,
			offset:     0,
			total:      25,
			wantNumber: 1,
			wantOutOf:  1,
			wantTotal:  25,
		},
		{
			name:       "given no results when calculating pages then returns zero pages",
			limit:      9,
			offset:     0,
			total:      0,
			wantNumber: 1,
			wantOutOf:  0,
			wantTotal:  0,
		},
		{
			name:       "given no limit and no results when calculating pages then returns zero pages",
			limit:      0,
			offset:     0,
			total:      0,
			wantNumber: 1,
			wantOutOf:  0,
			wantTotal:  0,
		},
		{
			name:       "given single result when calculating pages then returns one page",
			limit:      9,
			offset:     0,
			total:      1,
			wantNumber: 1,
			wantOutOf:  1,
			wantTotal:  1,
		},
		{
			name:       "given results less than limit when calculating pages then returns one page",
			limit:      9,
			offset:     0,
			total:      5,
			wantNumber: 1,
			wantOutOf:  1,
			wantTotal:  5,
		},
		{
			name:       "given last page when calculating with offset then returns correct page number",
			limit:      9,
			offset:     27,
			total:      35,
			wantNumber: 4,
			wantOutOf:  4,
			wantTotal:  35,
		},
		{
			name:       "given page with remainder when calculating last page then rounds up correctly",
			limit:      9,
			offset:     18,
			total:      26,
			wantNumber: 3,
			wantOutOf:  3,
			wantTotal:  26,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := usecases.CalculatePageInfo(tt.limit, tt.offset, tt.total)

			assert.Equal(t, tt.wantNumber, got.Number, "page number should match")
			assert.Equal(t, tt.wantOutOf, got.OutOf, "total pages should match")
			assert.Equal(t, tt.wantTotal, got.Total, "total items should match")
		})
	}
}
