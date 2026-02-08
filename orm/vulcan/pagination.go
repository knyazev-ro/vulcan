package vulcan

func (q *Query[T]) Paginate(page int, perPage int) *Query[T] {
	// page 1, perPage 10 = offset 0 limit 10
	// page 2, perPage 10 = offset 10 limit 10
	// page 3, perPage 10 = offset 20 limit 10
	currPerPage := 10
	if perPage > 0 {
		currPerPage = perPage
	}

	currPage := 1
	if page > 1 {
		currPage = page
	}

	offset := (currPage - 1) * currPerPage
	q.Offset(offset)
	q.Limit(currPerPage)
	return q
}

func (q *Query[T]) CursorPaginate(col string, cursor any, perPage int) *Query[T] {
	if len(q.orderExp) <= 0 {
		q.OrderBy("asc", col)
	}

	currPerPage := 10
	if perPage > 0 {
		currPerPage = perPage
	}
	if cursor != nil {
		q.Where(col, ">", cursor)
	}
	q.Limit(currPerPage)
	return q
}
