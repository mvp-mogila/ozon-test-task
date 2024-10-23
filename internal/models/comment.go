package models

type CommentModel struct {
	ID       string
	PostID   string
	ParentID *string
	Text     string
	Comments []*CommentModel
}
