package dto

import "github.com/mvp-mogila/ozon-test-task/internal/models"

func ConvertCommentModelToDTO(m models.Comment) Comment {
	return Comment{
		ID:       m.ID,
		PostID:   m.PostID,
		ParentID: m.ParentID,
		Content:  m.Content,
	}
}

func ConvertCreateCommentInputDTOToModel(in CreateCommentInput) models.CreateCommentInput {
	return models.CreateCommentInput{
		PostID:   in.PostID,
		ParentID: in.ParentID,
		Content:  in.Content,
	}
}
