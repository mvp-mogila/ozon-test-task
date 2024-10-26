package usecase

import (
	"context"
	"unicode/utf8"

	"github.com/mvp-mogila/ozon-test-task/internal/models"
	"github.com/mvp-mogila/ozon-test-task/pkg/logger"
	"github.com/mvp-mogila/ozon-test-task/pkg/pagination"
)

type PostsRepo interface {
	CreatePost(ctx context.Context, postInput models.CreatePostInput) (int, error)
	GetPost(ctx context.Context, postID int) (models.Post, error)
	GetPosts(ctx context.Context, offset int, limit int) ([]models.Post, error)
	GetPostCommentsPermission(ctx context.Context, postID int) (bool, error)
}

type CommentProvider interface {
	GetCommentsForPost(ctx context.Context, postID int, page int, commentsCount int) ([]models.Comment, error)
}

type PostUsecase struct {
	postRepo PostsRepo
}

func NewPostUsecase(r PostsRepo) *PostUsecase {
	return &PostUsecase{
		postRepo: r,
	}
}

func (u *PostUsecase) CreatePost(ctx context.Context, postInput models.CreatePostInput) (models.Post, error) {
	if utf8.RuneCountInString(postInput.Title) > 100 || utf8.RuneCountInString(postInput.Content) > 2000 {
		logger.Infof(ctx, "CreatePost: %s", models.ErrContentSizeExceeded.Error())
		return models.Post{}, models.ErrContentSizeExceeded
	}

	newPostID, err := u.postRepo.CreatePost(ctx, postInput)
	if err != nil {
		return models.Post{}, err
	}

	return models.Post{
		ID:            newPostID,
		Title:         postInput.Title,
		Content:       postInput.Content,
		AllowComments: postInput.AllowComments,
	}, nil
}

func (u *PostUsecase) GetPost(ctx context.Context, postID int) (models.Post, error) {
	if postID <= 0 {
		logger.Infof(ctx, "GetPost: %s", models.ErrInvalidIDValue.Error())
		return models.Post{}, models.ErrInvalidIDValue
	}

	post, err := u.postRepo.GetPost(ctx, postID)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func (u *PostUsecase) GetPosts(ctx context.Context, page int, postsCount int) ([]models.Post, error) {
	offset, limit, err := pagination.CountOffsetAndLimit(page, postsCount)
	if err != nil {
		logger.Infof(ctx, "GetPosts %s:", err.Error())
		return nil, err
	}

	posts, err := u.postRepo.GetPosts(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (u *PostUsecase) CheckCommentsPermission(ctx context.Context, postID int) (bool, error) {
	allow, err := u.postRepo.GetPostCommentsPermission(ctx, postID)
	if err != nil {
		return false, err
	}
	return allow, nil
}

func (u *PostUsecase) CheckPostExistance(ctx context.Context, postID int) (bool, error) {
	if postID <= 0 {
		logger.Infof(ctx, "CheckPostExistance: %s", models.ErrInvalidIDValue.Error())
		return false, models.ErrInvalidIDValue
	}

	_, err := u.postRepo.GetPost(ctx, postID)
	if err != nil {
		return false, err
	}
	return true, nil
}
