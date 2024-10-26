package pagination

import "errors"

var (
	errInvalidPage  = errors.New("invalid page value")
	errInvalidCount = errors.New("invalid count value")
)

func CountOffsetAndLimit(page, count int) (offset, limit int, err error) {
	if page <= 0 {
		return 0, 0, errInvalidPage
	}

	if count <= 0 {
		return 0, 0, errInvalidCount
	}

	err = nil
	limit = count
	offset = (page - 1) * limit
	return
}
