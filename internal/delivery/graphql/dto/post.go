package dto

import "github.com/mvp-mogila/ozon-test-task/internal/models"

func ConvertCreatePostInputDTOToModel(in CreatePostInput) models.CreatePostInput {
	return models.CreatePostInput{
		Title:         in.Title,
		Content:       in.Content,
		AllowComments: in.AllowComments,
	}
}

func ConvertPostModelDataToDTO(m models.Post) PostData {
	return PostData{
		ID:            m.ID,
		Title:         m.Title,
		Content:       m.Content,
		AllowComments: m.AllowComments,
	}
}

// func ConvertPostModelToDTO(m models.Post) Post {
// 	return Post{
// 		PostData: &PostData{
// 			ID:            m.ID,
// 			Title:         m.Title,
// 			Content:       m.Content,
// 			AllowComments: m.AllowComments,
// 		},
// 		Comments: make([]*Comment, 0),
// 	}
// }
