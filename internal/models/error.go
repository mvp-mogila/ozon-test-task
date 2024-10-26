package models

import "errors"

var (
	ErrInternal            = errors.New("internal server error")
	ErrInvalidIDValue      = errors.New("invalid ID value")
	ErrCommentsNotAllowed  = errors.New("comments for this post not allowed")
	ErrContentSizeExceeded = errors.New("content size exceeded")
	ErrPostNotFound        = errors.New("post with provided ID not found")
	ErrNoParentComment     = errors.New("parent comment with provided ID not found")
)
