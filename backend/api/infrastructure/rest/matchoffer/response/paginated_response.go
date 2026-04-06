package response

// PaginatedMatchOffersResponse represents a paginated response for match offers
type PaginatedMatchOffersResponse struct {
	Data       []MatchOfferResponse `json:"data"`
	Pagination PaginationInfo              `json:"pagination"`
}

// PaginationInfo contains pagination metadata
type PaginationInfo struct {
	Number int `json:"number"` // Current page number (1-based)
	OutOf  int `json:"out_of"` // Total number of pages
	Total  int `json:"total"`  // Total number of items matching the query
}
