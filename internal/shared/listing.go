package shared

import "time"

type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

// Common types

// SortOrder defines a sort rule for listing queries.
// IMPORTANT: Field must be validated against a whitelist before use in SQL to avoid injection.
type SortOrder struct {
	Field     string
	Direction SortDirection
}

// ListOptions defines common parameters for paginated and sorted list queries.
type ListOptions struct {
	CreatedAfter time.Time
	Limit        int // should be validated to enforce min / max limits
	SortOrders   []SortOrder
}
