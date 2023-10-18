package apiclient

// ResponseMetadata represent pagination and collection totals for multi-page responses.
type Meta struct {
	Pagination Pagination `json:"pagination"`
}

// Pagination and collection totals in the response.
type Pagination struct {
	Total       int   `json:"total"`
	Count       int   `json:"count"`
	PerPage     int   `json:"per_page"`
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	Links       Links `json:"links"`
}

// Pagination links for the current, previous and next parts of the whole collection.
type Links struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
	Next     string `json:"next"`
}
