package utils

type Pagination struct {
	Page     int
	PageSize int
}

func NewPagination(page, pageSize int) Pagination {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}
