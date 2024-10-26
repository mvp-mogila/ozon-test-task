package usecase

import (
	"context"
	"unicode/utf8"

	"github.com/mvp-mogila/ozon-test-task/internal/delivery/graphql/dto"
	"github.com/mvp-mogila/ozon-test-task/internal/models"
	"github.com/mvp-mogila/ozon-test-task/pkg/logger"
	"github.com/mvp-mogila/ozon-test-task/pkg/pagination"
)

const (
	defaultCommentsCount = 5
	limitCommentsCount   = 100
)

type CommentsRepo interface {
	CreateComment(ctx context.Context, commentInput models.CreateCommentInput) (int, error)
	GetCommentsForPost(ctx context.Context, postID int, offset int, limit int) ([]models.Comment, error)
	GetCommentsForComment(ctx context.Context, parentCommentID int, offset int, limit int) ([]models.Comment, error)
	CheckCommentExistance(ctx context.Context, commentID int) bool
}

type PostsProvider interface {
	CheckPostExistance(ctx context.Context, postID int) (bool, error)
	CheckCommentsPermission(ctx context.Context, postID int) (bool, error)
}

type CommentNotifier interface {
	CreateSubscriber(ctx context.Context, postID int) (ch chan *dto.Comment, subID int, err error)
	NotifySubscribers(ctx context.Context, comment models.Comment) error
	DeleteSubscriber(postID int, subID int) error
}

type CommentUsecase struct {
	commentRepo     CommentsRepo
	postProvider    PostsProvider
	commentNotifier CommentNotifier
}

func NewCommentUsecase(r CommentsRepo, pp PostsProvider, cn CommentNotifier) *CommentUsecase {
	return &CommentUsecase{
		commentRepo:     r,
		postProvider:    pp,
		commentNotifier: cn,
	}
}

func (u *CommentUsecase) CreateComment(ctx context.Context, commentInput models.CreateCommentInput) (models.Comment, error) {
	if commentInput.PostID <= 0 || (commentInput.ParentID != nil && *commentInput.ParentID <= 0) {
		logger.Infof(ctx, "CreateComment: %s", models.ErrInvalidIDValue.Error())
		return models.Comment{}, models.ErrInvalidIDValue
	}

	allow, err := u.postProvider.CheckCommentsPermission(ctx, commentInput.PostID)
	if err != nil {
		return models.Comment{}, err
	}
	if !allow {
		logger.Infof(ctx, "CreateComment: %s", models.ErrCommentsNotAllowed.Error())
		return models.Comment{}, models.ErrCommentsNotAllowed
	}

	if utf8.RuneCountInString(commentInput.Content) > 2000 {
		logger.Infof(ctx, "CreateComment: %s", models.ErrContentSizeExceeded.Error())
		return models.Comment{}, models.ErrContentSizeExceeded
	}

	if commentInput.ParentID != nil && *commentInput.ParentID > 0 {
		if exist := u.commentRepo.CheckCommentExistance(ctx, *commentInput.ParentID); !exist {
			logger.Infof(ctx, "CreateComment: %s", models.ErrNoParentComment.Error())
			return models.Comment{}, models.ErrNoParentComment
		}
	}

	newCommentID, err := u.commentRepo.CreateComment(ctx, commentInput)
	if err != nil {
		return models.Comment{}, err
	}

	newComment := models.Comment{
		ID:       newCommentID,
		PostID:   commentInput.PostID,
		ParentID: commentInput.ParentID,
		Content:  commentInput.Content,
	}

	if err := u.commentNotifier.NotifySubscribers(ctx, newComment); err != nil {
		logger.Errorf(ctx, "CreateComment notify: %s", err.Error())
	}
	return newComment, nil
}

func (u *CommentUsecase) GetCommentsForPost(ctx context.Context, postID int, page int, commentsCount int) ([]models.Comment, error) {
	if postID <= 0 {
		logger.Infof(ctx, "GetCommentsForPost: %s", models.ErrInvalidIDValue.Error())
		return nil, models.ErrInvalidIDValue
	}

	if page == 0 {
		page = 1
	}
	if commentsCount == 0 {
		commentsCount = defaultCommentsCount
	}

	offset, limit, err := pagination.CountOffsetAndLimit(page, commentsCount)
	if err != nil {
		logger.Infof(ctx, "GetCommentsForPost: %s", err.Error())
		return nil, err
	}

	comments, err := u.commentRepo.GetCommentsForPost(ctx, postID, offset, limit)
	if err != nil {
		return nil, err
	}
	return comments, err
}

func (u *CommentUsecase) GetCommentsForComment(ctx context.Context, parentCommentID int, commentsCount int) ([]models.Comment, error) {
	if parentCommentID <= 0 {
		logger.Infof(ctx, "GetCommentsForComment: %s", models.ErrInvalidIDValue.Error())
		return nil, models.ErrInvalidIDValue
	}

	if commentsCount == 0 {
		commentsCount = defaultCommentsCount
	}

	offset, limit, err := pagination.CountOffsetAndLimit(1, commentsCount)
	if err != nil {
		logger.Infof(ctx, "GetCommentsForComment: %s", err.Error())
		return nil, err
	}

	comments, err := u.commentRepo.GetCommentsForComment(ctx, parentCommentID, offset, limit)
	if err != nil {
		return nil, err
	}
	return comments, err
}

func (u *CommentUsecase) CreateSubscriber(ctx context.Context, postID int) (<-chan *dto.Comment, error) {
	if postID <= 0 {
		logger.Infof(ctx, "CreateSubscriber: %s", models.ErrInvalidIDValue.Error())
		return nil, models.ErrInvalidIDValue
	}

	exist, err := u.postProvider.CheckPostExistance(ctx, postID)
	if err != nil {
		return nil, err
	}

	if !exist {
		logger.Infof(ctx, "CreateSubscriber: %s", models.ErrPostNotFound.Error())
		return nil, models.ErrPostNotFound
	}

	commentCh, subID, err := u.commentNotifier.CreateSubscriber(ctx, postID)
	if err != nil {
		logger.Errorf(ctx, "CreateSubscriber notifier: %s", err.Error())
		return nil, err
	}

	go func() {
		<-ctx.Done()

		if err := u.commentNotifier.DeleteSubscriber(postID, subID); err != nil {
			ctx := context.Background()
			logger.Errorf(ctx, "DeleteSubscriber %d: %s", subID, err.Error())
		}
	}()
	return commentCh, nil
}
