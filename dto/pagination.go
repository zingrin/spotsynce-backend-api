package dto

// PaginationParams holds common pagination defaults and normalization.
type PaginationParams struct {
	Page  int
	Limit int
}

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

// NormalizePagination applies defaults and caps to pagination parameters.
func NormalizePagination(page, limit int) PaginationParams {
	if page < 1 {
		page = DefaultPage
	}
	if limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	return PaginationParams{Page: page, Limit: limit}
}

// Offset calculates the database offset from page and limit.
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.Limit
}

// TotalPages calculates total pages from total record count.
func TotalPages(total int64, limit int) int {
	if total == 0 {
		return 0
	}
	pages := int(total) / limit
	if int(total)%limit > 0 {
		pages++
	}
	return pages
}
