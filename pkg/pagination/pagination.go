package pagination

import "math"

// PaginationParameter holds common pagination inputs
type PaginationParameter struct {
	Page  int    `form:"page" binding:"omitempty,min=1"`
	Limit int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Sort  string `form:"sort"`
}

// PaginationMetadata holds pagination information for the client
type PaginationMetadata struct {
	Total       int64 `json:"total"`
	Page        int   `json:"page"`
	Limit       int   `json:"limit"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}

// PaginationResult holds the data and its pagination metadata
type PaginationResult struct {
	Data     interface{}        `json:"data"`
	Metadata PaginationMetadata `json:"metadata"`
}

// CalculateMetadata generates PaginationMetadata from total count and parameters
func CalculateMetadata(total int64, page, limit int) PaginationMetadata {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return PaginationMetadata{
		Total:       total,
		Page:        page,
		Limit:       limit,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}
}
