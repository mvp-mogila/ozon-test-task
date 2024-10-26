package inmemory

import (
	"context"

	"github.com/mvp-mogila/ozon-test-task/internal/models"
	"github.com/mvp-mogila/ozon-test-task/internal/repository/dao"
	"github.com/mvp-mogila/ozon-test-task/pkg/inmemory"
)

const defaultPostStorageSize = 100

type PostInMemoryRepository struct {
	totalSize int
	lastID    int
	storage   *inmemory.InMemoryStorage
}

func NewPostInMemoryRepository() *PostInMemoryRepository {
	return &PostInMemoryRepository{
		lastID:  0,
		storage: inmemory.NewInMemoryStorage(defaultPostStorageSize),
	}
}

func (r *PostInMemoryRepository) CreatePost(ctx context.Context, postInput models.CreatePostInput) (int, error) {
	newPost := dao.ConvertCreatePostModelToDAO(postInput)
	newPost.ID = r.lastID + 1

	err := r.storage.AddObject(ctx, &newPost)
	if err != nil {
		return 0, err
	}

	r.lastID++
	r.totalSize++
	return newPost.ID, nil
}

func (r *PostInMemoryRepository) GetPost(ctx context.Context, postID int) (models.Post, error) {
	object := r.storage.GetObjectByID(ctx, postID)

	if object == nil {
		return models.Post{}, models.ErrPostNotFound
	}

	postDAO, ok := object.(*dao.Post)
	if !ok {
		return models.Post{}, models.ErrInternal
	}

	post := dao.ConvertPostDAOToModel(*postDAO)
	return post, nil
}

func (r *PostInMemoryRepository) GetPosts(ctx context.Context, offset int, limit int) ([]models.Post, error) {
	objects, err := r.storage.GetManyObjects(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	posts := make([]models.Post, 0, limit)
	for _, obj := range objects {
		postDAO, ok := obj.(*dao.Post)
		if !ok {
			return nil, models.ErrInternal
		}
		posts = append(posts, dao.ConvertPostDAOToModel(*postDAO))
	}
	return posts, nil
}

func (r *PostInMemoryRepository) GetPostCommentsPermission(ctx context.Context, postID int) (bool, error) {
	object := r.storage.GetObjectByID(ctx, postID)

	if object == nil {
		return false, models.ErrPostNotFound
	}

	postDAO, ok := object.(*dao.Post)
	if !ok {
		return false, models.ErrInternal
	}
	return postDAO.AllowComments, nil
}
