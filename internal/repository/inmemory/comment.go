package inmemory

import (
	"context"

	"github.com/mvp-mogila/ozon-test-task/internal/models"
	"github.com/mvp-mogila/ozon-test-task/internal/repository/dao"
	"github.com/mvp-mogila/ozon-test-task/pkg/inmemory"
)

const defaultCommentStorageSize = 100

type CommentInMemoryRepository struct {
	totalSize int
	lastID    int
	storage   *inmemory.InMemoryStorage
}

func NewCommentInMemoryRepository() *CommentInMemoryRepository {
	return &CommentInMemoryRepository{
		totalSize: 0,
		lastID:    0,
		storage:   inmemory.NewInMemoryStorage(defaultCommentStorageSize),
	}
}

func (r *CommentInMemoryRepository) CreateComment(ctx context.Context, commentInput models.CreateCommentInput) (int, error) {
	newComment := dao.ConvertCreateCommentModelToDAO(commentInput)
	newComment.ID = r.lastID + 1

	err := r.storage.AddObject(ctx, &newComment)
	if err != nil {
		return 0, err
	}

	r.lastID++
	r.totalSize++
	return newComment.ID, nil
}

func (r *CommentInMemoryRepository) GetCommentsForPost(ctx context.Context, postID int, offset int, limit int) ([]models.Comment, error) {
	objects, err := r.storage.GetManyObjects(ctx, 0, r.totalSize)
	if err != nil {
		return nil, err
	}

	comments := make([]models.Comment, 0, limit)
	var count, skipped, idx int
	for idx = 0; idx < len(objects) && count < limit; idx++ {
		commentDAO, ok := objects[idx].(*dao.Comment)
		if !ok {
			return nil, models.ErrInternal
		}
		if commentDAO.PostID == postID && commentDAO.ParentID == nil {
			if skipped < offset {
				skipped++
			} else {
				comment := dao.ConvertCommentDAOToModel(*commentDAO)
				comments = append(comments, comment)
				count++
			}
		}
	}

	return comments, nil
}

func (r *CommentInMemoryRepository) GetCommentsForComment(ctx context.Context, parentCommentID int, offset int, limit int) ([]models.Comment, error) {
	objects, err := r.storage.GetManyObjects(ctx, offset, r.totalSize)
	if err != nil {
		return nil, err
	}

	comments := make([]models.Comment, 0, limit)
	var count, idx int
	for idx = 0; idx < len(objects) && count < limit; idx++ {
		commentDAO, ok := objects[idx].(*dao.Comment)
		if !ok {
			return nil, models.ErrInternal
		}
		if commentDAO.ParentID != nil && *commentDAO.ParentID == parentCommentID {
			comment := dao.ConvertCommentDAOToModel(*commentDAO)
			comments = append(comments, comment)
			count++
		}
	}

	return comments, nil
}

func (r *CommentInMemoryRepository) CheckCommentExistance(ctx context.Context, commentID int) bool {
	if obj := r.storage.GetObjectByID(ctx, commentID); obj != nil {
		return true
	}
	return false
}
