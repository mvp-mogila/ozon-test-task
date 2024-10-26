package dao

import "github.com/mvp-mogila/ozon-test-task/internal/models"

type Comment struct {
	ID       int
	PostID   int
	ParentID *int
	Content  string
}

func (c *Comment) GetID() int {
	return c.ID
}

func ConvertCreateCommentModelToDAO(m models.CreateCommentInput) Comment {
	return Comment{
		PostID:   m.PostID,
		ParentID: m.ParentID,
		Content:  m.Content,
	}
}

func ConvertCommentDAOToModel(d Comment) models.Comment {
	return models.Comment{
		ID:       d.ID,
		PostID:   d.PostID,
		ParentID: d.ParentID,
		Content:  d.Content,
	}
}
