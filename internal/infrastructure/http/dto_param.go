package http

type PaginationQuery struct {
	Page     int32 `form:"page"`
	PageSize int32 `form:"page_size"`
}

func (p *PaginationQuery) ActualPage() int32 {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}

func (p *PaginationQuery) Offset() int32 {
	page := p.Page
	if page < 1 {
		page = 1
	}

	return (page - 1) * p.PageSize
}

func (p *PaginationQuery) Limit() int32 {
	pageSize := p.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	return pageSize
}
