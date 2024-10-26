package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mvp-mogila/ozon-test-task/internal/models"
	"github.com/mvp-mogila/ozon-test-task/internal/repository/dao"
	"github.com/mvp-mogila/ozon-test-task/pkg/logger"
)

type PostPostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostPostgresRepository(pool *pgxpool.Pool) *PostPostgresRepository {
	return &PostPostgresRepository{
		db: pool,
	}
}

func (r *PostPostgresRepository) CreatePost(ctx context.Context, postInput models.CreatePostInput) (int, error) {
	q := `INSERT INTO post(title, content, allow_comments) VALUES($1, $2, $3) RETURNING id;`

	var newPostID int
	if err := r.db.QueryRow(ctx, q, postInput.Title, postInput.Content, postInput.AllowComments).Scan(&newPostID); err != nil {
		logger.Errorf(ctx, "CreatePost quering: %s", err.Error())
		return 0, err
	}
	return newPostID, nil
}

func (r *PostPostgresRepository) GetPost(ctx context.Context, postID int) (models.Post, error) {
	q := `SELECT id, title, content, allow_comments FROM post WHERE id = $1;`

	var postDAO dao.Post
	if err := r.db.QueryRow(ctx, q, postID).Scan(&postDAO.ID, &postDAO.Title, &postDAO.Content, &postDAO.AllowComments); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Infof(ctx, "GetPost %d not found", postID)
			return models.Post{}, models.ErrPostNotFound
		}
		logger.Errorf(ctx, "GetPost %d quering: %s", postID, err.Error())
		return models.Post{}, err
	}

	post := dao.ConvertPostDAOToModel(postDAO)
	return post, nil
}

func (r *PostPostgresRepository) GetPosts(ctx context.Context, offset int, limit int) ([]models.Post, error) {
	q := `SELECT id, title, content, allow_comments FROM post LIMIT $1 OFFSET $2;`

	rows, err := r.db.Query(ctx, q, limit, offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Infof(ctx, "GetPosts not found")
			return nil, models.ErrPostNotFound
		}
		logger.Errorf(ctx, "GetPosts quering: %s", err.Error())
		return nil, err
	}

	posts := make([]models.Post, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		var postDAO dao.Post
		if err := rows.Scan(&postDAO.ID, &postDAO.Title, &postDAO.Content, &postDAO.AllowComments); err != nil {
			logger.Errorf(ctx, "GetPosts scanning: %s", err.Error())
			return nil, err
		}

		post := dao.ConvertPostDAOToModel(postDAO)
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostPostgresRepository) GetPostCommentsPermission(ctx context.Context, postID int) (bool, error) {
	q := `SELECT allow_comments FROM post WHERE id = $1;`

	var allowComments bool
	err := r.db.QueryRow(ctx, q, postID).Scan(&allowComments)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Infof(ctx, "GetPostCommentsPermission %d not found", postID)
			return false, models.ErrPostNotFound
		}
		logger.Errorf(ctx, "GetPostCommentsPermission scanning: %s", err.Error())
		return false, err
	}
	return allowComments, nil
}
