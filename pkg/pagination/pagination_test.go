package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateMetadata(t *testing.T) {
	tests := []struct {
		name     string
		total    int64
		page     int
		limit    int
		expected PaginationMetadata
	}{
		{
			name:  "first page success",
			total: 25,
			page:  1,
			limit: 10,
			expected: PaginationMetadata{
				Total:       25,
				Page:        1,
				Limit:       10,
				TotalPages:  3,
				HasNext:     true,
				HasPrevious: false,
			},
		},
		{
			name:  "last page success",
			total: 25,
			page:  3,
			limit: 10,
			expected: PaginationMetadata{
				Total:       25,
				Page:        3,
				Limit:       10,
				TotalPages:  3,
				HasNext:     false,
				HasPrevious: true,
			},
		},
		{
			name:  "middle page success",
			total: 25,
			page:  2,
			limit: 10,
			expected: PaginationMetadata{
				Total:       25,
				Page:        2,
				Limit:       10,
				TotalPages:  3,
				HasNext:     true,
				HasPrevious: true,
			},
		},
		{
			name:  "exact pages success",
			total: 20,
			page:  1,
			limit: 10,
			expected: PaginationMetadata{
				Total:       20,
				Page:        1,
				Limit:       10,
				TotalPages:  2,
				HasNext:     true,
				HasPrevious: false,
			},
		},
		{
			name:  "zero total",
			total: 0,
			page:  1,
			limit: 10,
			expected: PaginationMetadata{
				Total:       0,
				Page:        1,
				Limit:       10,
				TotalPages:  0,
				HasNext:     false,
				HasPrevious: false,
			},
		},
		{
			name:  "zero limit defaults to 10",
			total: 25,
			page:  1,
			limit: 0,
			expected: PaginationMetadata{
				Total:       25,
				Page:        1,
				Limit:       10,
				TotalPages:  3,
				HasNext:     true,
				HasPrevious: false,
			},
		},
		{
			name:  "zero page defaults to 1",
			total: 25,
			page:  0,
			limit: 10,
			expected: PaginationMetadata{
				Total:       25,
				Page:        1,
				Limit:       10,
				TotalPages:  3,
				HasNext:     true,
				HasPrevious: false,
			},
		},
		{
			name:  "negative total handled",
			total: -1,
			page:  1,
			limit: 10,
			expected: PaginationMetadata{
				Total:       -1,
				Page:        1,
				Limit:       10,
				TotalPages:  0,
				HasNext:     false,
				HasPrevious: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateMetadata(tt.total, tt.page, tt.limit)
			assert.Equal(t, tt.expected, result)
		})
	}
}
