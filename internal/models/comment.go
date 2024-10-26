package models

type Comment struct {
	ID       int
	PostID   int
	ParentID *int
	Content  string
	Comments []*Comment
}

type CreateCommentInput struct {
	PostID   int
	ParentID *int
	Content  string
}
