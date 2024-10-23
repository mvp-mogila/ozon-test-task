package models

type PostModel struct {
	ID             string
	Title          string
	Content        string
	CommetsAllowed bool
	Comments       []CommentModel
}
