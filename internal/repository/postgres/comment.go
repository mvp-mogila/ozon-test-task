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

type CommentPostgresRepository struct {
	db *pgxpool.Pool
}

func NewCommentPostgresRepository(pool *pgxpool.Pool) *CommentPostgresRepository {
	return &CommentPostgresRepository{
		db: pool,
	}
}

func (r *CommentPostgresRepository) CreateComment(ctx context.Context, commentInput models.CreateCommentInput) (int, error) {
	q := `INSERT INTO comment(content, post_id, parent_id) VALUES($1, $2, $3) RETURNING id;`

	var newCommentID int
	if err := r.db.QueryRow(ctx, q, commentInput.Content, commentInput.PostID, commentInput.ParentID).Scan(&newCommentID); err != nil {
		logger.Errorf(ctx, "CreateComment quering: %s", err.Error())
		return 0, err
	}
	return newCommentID, nil
}

func (r *CommentPostgresRepository) GetCommentsForPost(ctx context.Context, postID int, offset int, limit int) ([]models.Comment, error) {
	q := `SELECT id, content, post_id, parent_id
		FROM comment WHERE post_id = $1 AND parent_id IS NULL
		LIMIT $2
		OFFSET $3;`

	rows, err := r.db.Query(ctx, q, postID, limit, offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Infof(ctx, "GetCommentsForPost %d not found", postID)
			return nil, nil
		}
		logger.Errorf(ctx, "GetCommentsForPost %d quering: %s", postID, err.Error())
		return nil, err
	}

	comments := make([]models.Comment, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		var commentDAO dao.Comment
		if err := rows.Scan(&commentDAO.ID, &commentDAO.Content, &commentDAO.PostID, &commentDAO.ParentID); err != nil {
			logger.Errorf(ctx, "GetCommentsForPost %d scanning: %s", postID, err.Error())
			return nil, err
		}

		comment := dao.ConvertCommentDAOToModel(commentDAO)
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *CommentPostgresRepository) GetCommentsForComment(ctx context.Context, parentCommentID int, offset int, limit int) ([]models.Comment, error) {
	q := `SELECT id, content, post_id, parent_id
		FROM comment WHERE parent_id = $1
		LIMIT $2
		OFFSET $3;`

	rows, err := r.db.Query(ctx, q, parentCommentID, limit, offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Infof(ctx, "GetCommentsForComment %d not found", parentCommentID)
			return nil, nil
		}
		logger.Errorf(ctx, "GetCommentsForComment %d quering: %s", parentCommentID, err.Error())
		return nil, err
	}

	comments := make([]models.Comment, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		var commentDAO dao.Comment
		if err := rows.Scan(&commentDAO.ID, &commentDAO.Content, &commentDAO.PostID, &commentDAO.ParentID); err != nil {
			logger.Errorf(ctx, "GetCommentsForComment %d scanning: %s", parentCommentID, err.Error())
			return nil, err
		}

		comment := dao.ConvertCommentDAOToModel(commentDAO)
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *CommentPostgresRepository) CheckCommentExistance(ctx context.Context, commentID int) bool {
	q := `SELECT id FROM comment WHERE id = $1;`

	var id int
	if err := r.db.QueryRow(ctx, q, commentID).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Infof(ctx, "CheckCommentExistance %d not found", commentID)
			return false
		}
		logger.Errorf(ctx, "CheckCommentExistance %d quering %s:", commentID, err.Error())
		return false
	}
	return true
}
