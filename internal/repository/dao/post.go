package dao

import "github.com/mvp-mogila/ozon-test-task/internal/models"

type Post struct {
	ID            int
	Title         string
	Content       string
	AllowComments bool
}

func (p *Post) GetID() int {
	return p.ID
}

func ConvertCreatePostModelToDAO(m models.CreatePostInput) Post {
	return Post{
		Title:         m.Title,
		Content:       m.Content,
		AllowComments: m.AllowComments,
	}
}

func ConvertPostDAOToModel(d Post) models.Post {
	return models.Post{
		ID:            d.ID,
		Title:         d.Title,
		Content:       d.Content,
		AllowComments: d.AllowComments,
	}
}
