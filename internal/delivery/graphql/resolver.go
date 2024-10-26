package graphql

import (
	"context"

	"github.com/mvp-mogila/ozon-test-task/internal/delivery/graphql/dto"
	"github.com/mvp-mogila/ozon-test-task/internal/models"
)

type PostsUsecase interface {
	CreatePost(ctx context.Context, postInput models.CreatePostInput) (models.Post, error)
	GetPost(ctx context.Context, id int) (models.Post, error)
	GetPosts(ctx context.Context, page int, postsCount int) ([]models.Post, error)
}

type CommentsUsecase interface {
	CreateComment(ctx context.Context, commentInput models.CreateCommentInput) (models.Comment, error)
	GetCommentsForPost(ctx context.Context, postID int, page int, commentsCount int) ([]models.Comment, error)
	GetCommentsForComment(ctx context.Context, parentCommentID int, commentsCount int) ([]models.Comment, error)
	CreateSubscriber(ctx context.Context, postID int) (<-chan *dto.Comment, error)
}

type Resolver struct {
	PostUsecase    PostsUsecase
	CommentUsecase CommentsUsecase
}
