package models

type Post struct {
	ID            int
	Title         string
	Content       string
	AllowComments bool
	Comments      []Comment
}

type CreatePostInput struct {
	Title         string
	Content       string
	AllowComments bool
}
