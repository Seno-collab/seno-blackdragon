package dto

type PaginationRequest struct {
	Page     int    `json:"page" form:"page" binding:"gte=1"`
	PageSize int    `json:"page_size" form:"page_size" binding:"gte=1,lte=100"`
	SortBy   string `json:"sort_by,omitempty" form:"sort_by"`
	OrderBy  string `json:"order_by,omitempty" form:"order_by"`
	Query    string `json:"query,omitempty" form:"query"`
}

type PaginationResponse[T any] struct {
	Items       []T   `json:"items"`
	Total       int64 `json:"total"`
	Page        int   `json:"page"`
	PageSize    int   `json:"page_size"`
	TotalPages  int   `json:"total_pages"`
	HasNextPage bool  `json:"has_next_page"`
}
