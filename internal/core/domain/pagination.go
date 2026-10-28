package domain

func ClampPagination(page, limit int) (int, int) {
	switch {
	case page <= 0:
		page = PageDefaultValue
	case page > PageMaxValue:
		page = PageMaxValue
	}

	switch {
	case limit <= 0:
		limit = LimitDefaultValue
	case limit > LimitMaxValue:
		limit = LimitMaxValue
	}

	return page, limit
}
