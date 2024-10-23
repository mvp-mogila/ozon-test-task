package graphql

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type PostUsecase interface {
	CreatePost()
	GetPosts()
	GetPost()
}

type CommentsUsecase interface {
	CreateComment()
}

type Resolver struct {
	postUsecase    PostUsecase
	commentUsecase CommentsUsecase
}
